package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"movie-web/internal/db"
	"movie-web/internal/model"
)

const (
	bookingHoldDuration = 10 * time.Minute
	standardSeatPrice   = 90000
	vipSeatPrice        = 130000
)

var cinemaLocation = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)

type bookingResponse struct {
	ID            uint           `json:"id"`
	Code          string         `json:"code"`
	Showtime      model.Showtime `json:"showtime"`
	SeatCodes     []string       `json:"seat_codes"`
	TotalAmount   int            `json:"total_amount"`
	Status        string         `json:"status"`
	PaymentStatus string         `json:"payment_status"`
	ExpiresAt     time.Time      `json:"expires_at"`
	PaidAt        *time.Time     `json:"paid_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
}

func bookingDB(c *gin.Context) bool {
	if db.DB != nil {
		return true
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Cơ sở dữ liệu đặt vé chưa sẵn sàng"})
	return false
}

func GetCinemas(c *gin.Context) {
	if !bookingDB(c) {
		return
	}
	var cinemas []model.Cinema
	if err := db.DB.Order("id asc").Find(&cinemas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải danh sách rạp"})
		return
	}
	c.JSON(http.StatusOK, cinemas)
}

func GetShowtimes(c *gin.Context) {
	if !bookingDB(c) {
		return
	}
	movieID, err := strconv.Atoi(c.Query("movie_id"))
	if err != nil || movieID <= 0 || strings.TrimSpace(c.Query("movie_title")) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần chọn phim hợp lệ trước khi tải suất chiếu"})
		return
	}
	showDate := c.DefaultQuery("date", time.Now().In(cinemaLocation).Format("2006-01-02"))
	date, err := time.ParseInLocation("2006-01-02", showDate, cinemaLocation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ngày chiếu không hợp lệ"})
		return
	}
	today := time.Now().In(cinemaLocation)
	if date.Before(time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, cinemaLocation)) || date.After(today.AddDate(0, 0, 14)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chỉ có thể xem lịch chiếu trong 14 ngày tới"})
		return
	}
	if err := ensureDemoShowtimes(movieID, c.Query("movie_title"), c.Query("poster_url"), date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo lịch chiếu mẫu", "detail": err.Error()})
		return
	}

	query := db.DB.Preload("Cinema").Preload("Auditorium").Where("movie_id = ? AND start_at >= ? AND start_at < ? AND start_at > ?", movieID, date, date.AddDate(0, 0, 1), time.Now())
	if cinemaID := c.Query("cinema_id"); cinemaID != "" {
		query = query.Where("cinema_id = ?", cinemaID)
	}
	var showtimes []model.Showtime
	if err := query.Order("start_at asc").Find(&showtimes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải suất chiếu"})
		return
	}
	type showtimeItem struct {
		model.Showtime
		AvailableSeats int64 `json:"available_seats"`
	}
	result := make([]showtimeItem, 0, len(showtimes))
	for _, showtime := range showtimes {
		if err := ensureShowtimeSeats(db.DB, showtime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo sơ đồ ghế"})
			return
		}
		var available int64
		db.DB.Model(&model.ShowSeat{}).Where("showtime_id = ? AND status = ?", showtime.ID, "available").Count(&available)
		result = append(result, showtimeItem{Showtime: showtime, AvailableSeats: available})
	}
	c.JSON(http.StatusOK, result)
}

func ensureDemoShowtimes(movieID int, title, poster string, date time.Time) error {
	var cinemas []model.Cinema
	if err := db.DB.Order("id asc").Find(&cinemas).Error; err != nil {
		return err
	}
	times := [][2]int{{10, 0}, {13, 30}, {19, 0}}
	for _, cinema := range cinemas {
		var rooms []model.Auditorium
		if err := db.DB.Where("cinema_id = ?", cinema.ID).Order("id asc").Find(&rooms).Error; err != nil {
			return err
		}
		for _, clock := range times {
			startAt := time.Date(date.Year(), date.Month(), date.Day(), clock[0], clock[1], 0, 0, cinemaLocation)
			var existing model.Showtime
			if err := db.DB.Where("movie_id = ? AND cinema_id = ? AND start_at = ?", movieID, cinema.ID, startAt).First(&existing).Error; err == nil {
				continue
			}
			for _, room := range rooms {
				var occupied int64
				db.DB.Model(&model.Showtime{}).Where("cinema_id = ? AND auditorium_id = ? AND start_at = ?", cinema.ID, room.ID, startAt).Count(&occupied)
				if occupied > 0 {
					continue
				}
				showtime := model.Showtime{
					MovieID: movieID, MovieTitle: strings.TrimSpace(title), PosterURL: poster,
					CinemaID: cinema.ID, AuditoriumID: room.ID, StartAt: startAt,
					StandardPrice: standardSeatPrice, VIPPrice: vipSeatPrice,
				}
				if err := db.DB.Create(&showtime).Error; err != nil {
					if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
						break
					}
					return err
				}
				break
			}
		}
	}
	return nil
}

func ensureShowtimeSeats(database *gorm.DB, showtime model.Showtime) error {
	room := showtime.Auditorium
	if room.ID == 0 {
		if err := database.First(&room, showtime.AuditoriumID).Error; err != nil {
			return err
		}
	}
	seats := make([]model.ShowSeat, 0, room.RowCount*room.SeatsPerRow)
	for rowIndex := 0; rowIndex < room.RowCount; rowIndex++ {
		row := string(rune('A' + rowIndex))
		for number := 1; number <= room.SeatsPerRow; number++ {
			seatType := "standard"
			if rowIndex >= room.RowCount/2-1 && rowIndex <= room.RowCount/2 {
				seatType = "vip"
			}
			seats = append(seats, model.ShowSeat{
				ShowtimeID: showtime.ID, Code: fmt.Sprintf("%s%d", row, number),
				Row: row, Number: number, Type: seatType, Status: "available",
			})
		}
	}
	return database.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "showtime_id"}, {Name: "code"}}, DoNothing: true}).Create(&seats).Error
}

func releaseExpiredBookings(tx *gorm.DB) error {
	now := time.Now()
	expiredBookings := tx.Model(&model.Booking{}).Select("id").Where("status = ? AND expires_at <= ?", "AWAITING_PAYMENT", now)
	if err := tx.Model(&model.ShowSeat{}).Where("booking_id IN (?)", expiredBookings).Updates(map[string]interface{}{"status": "available", "booking_id": nil}).Error; err != nil {
		return err
	}
	return tx.Model(&model.Booking{}).Where("status = ? AND expires_at <= ?", "AWAITING_PAYMENT", now).Updates(map[string]interface{}{"status": "EXPIRED", "payment_status": "EXPIRED"}).Error
}

func GetShowtimeSeats(c *gin.Context) {
	if !bookingDB(c) {
		return
	}
	showtimeID, err := strconv.Atoi(c.Param("id"))
	if err != nil || showtimeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Suất chiếu không hợp lệ"})
		return
	}
	var showtime model.Showtime
	if err := db.DB.Preload("Cinema").Preload("Auditorium").First(&showtime, showtimeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy suất chiếu"})
		return
	}
	if err := releaseExpiredBookings(db.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật ghế giữ chỗ hết hạn"})
		return
	}
	if err := ensureShowtimeSeats(db.DB, showtime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải sơ đồ ghế"})
		return
	}
	var seats []model.ShowSeat
	if err := db.DB.Where("showtime_id = ?", showtime.ID).Order("row asc, number asc").Find(&seats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải danh sách ghế"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"showtime": showtime, "seats": seats})
}

type createBookingRequest struct {
	ShowtimeID uint     `json:"showtime_id" binding:"required"`
	SeatCodes  []string `json:"seat_codes" binding:"required,min=1,max=8,dive,required"`
}

func CreateBooking(c *gin.Context) {
	if !bookingDB(c) {
		return
	}
	var request createBookingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chọn từ 1 đến 8 ghế hợp lệ"})
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Vui lòng đăng nhập để đặt vé"})
		return
	}
	codes := make([]string, 0, len(request.SeatCodes))
	seen := make(map[string]bool)
	for _, code := range request.SeatCodes {
		code = strings.ToUpper(strings.TrimSpace(code))
		if code == "" || seen[code] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Danh sách ghế bị trùng hoặc không hợp lệ"})
			return
		}
		seen[code] = true
		codes = append(codes, code)
	}
	sort.Strings(codes)

	var created model.Booking
	var apiErr error
	statusCode := http.StatusCreated
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := releaseExpiredBookings(tx); err != nil {
			return err
		}
		var showtime model.Showtime
		if err := tx.First(&showtime, request.ShowtimeID).Error; err != nil {
			apiErr = fmt.Errorf("Không tìm thấy suất chiếu")
			statusCode = http.StatusNotFound
			return nil
		}
		if !showtime.StartAt.After(time.Now()) {
			apiErr = fmt.Errorf("Suất chiếu này đã bắt đầu hoặc đã kết thúc")
			statusCode = http.StatusConflict
			return nil
		}
		var seats []model.ShowSeat
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("showtime_id = ? AND code IN ?", showtime.ID, codes).Find(&seats).Error; err != nil {
			return err
		}
		if len(seats) != len(codes) {
			apiErr = fmt.Errorf("Một hoặc nhiều ghế không thuộc sơ đồ suất chiếu này")
			statusCode = http.StatusBadRequest
			return nil
		}
		amount := 0
		for _, seat := range seats {
			if seat.Status != "available" {
				apiErr = fmt.Errorf("Ghế %s vừa được người khác giữ hoặc đặt", seat.Code)
				statusCode = http.StatusConflict
				return nil
			}
			if seat.Type == "vip" {
				amount += showtime.VIPPrice
			} else {
				amount += showtime.StandardPrice
			}
		}
		bookingCode, err := makeTicketCode()
		if err != nil {
			return err
		}
		now := time.Now()
		created = model.Booking{
			Code: bookingCode, UserID: userID.(uint), ShowtimeID: showtime.ID,
			TotalAmount: amount, Status: "AWAITING_PAYMENT", PaymentStatus: "PENDING",
			ExpiresAt: now.Add(bookingHoldDuration), CreatedAt: now,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		for _, seat := range seats {
			result := tx.Model(&model.ShowSeat{}).Where("id = ? AND status = ?", seat.ID, "available").Updates(map[string]interface{}{"status": "reserved", "booking_id": created.ID})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				apiErr = fmt.Errorf("Ghế %s vừa được người khác giữ hoặc đặt", seat.Code)
				statusCode = http.StatusConflict
				return apiErr
			}
			if err := tx.Create(&model.BookingSeat{BookingID: created.ID, ShowSeatID: seat.ID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if apiErr == nil {
			apiErr = err
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{"error": apiErr.Error()})
		return
	}
	if apiErr != nil {
		c.JSON(statusCode, gin.H{"error": apiErr.Error()})
		return
	}
	detail, err := loadBookingResponse(db.DB, created)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Đã giữ ghế nhưng không thể tải thông tin vé"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"booking": detail, "message": "Đã giữ ghế. Hoàn tất thanh toán trong 10 phút để xác nhận vé."})
}

func SimulateBookingPayment(c *gin.Context) {
	updateBookingStatus(c, "pay")
}

func CancelBooking(c *gin.Context) {
	updateBookingStatus(c, "cancel")
}

func updateBookingStatus(c *gin.Context, action string) {
	if !bookingDB(c) {
		return
	}
	userID, _ := c.Get("userID")
	var booking model.Booking
	var apiErr error
	statusCode := http.StatusOK
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := releaseExpiredBookings(tx); err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("code = ? AND user_id = ?", c.Param("code"), userID).First(&booking).Error; err != nil {
			apiErr = fmt.Errorf("Không tìm thấy mã đặt vé")
			statusCode = http.StatusNotFound
			return nil
		}
		if action == "pay" && booking.PaymentStatus == "PAID" {
			return nil
		}
		if booking.Status != "AWAITING_PAYMENT" || time.Now().After(booking.ExpiresAt) {
			apiErr = fmt.Errorf("Đơn đặt vé đã hết hạn hoặc không còn chờ thanh toán")
			statusCode = http.StatusConflict
			return nil
		}
		now := time.Now()
		if action == "cancel" {
			booking.Status = "CANCELLED"
			booking.PaymentStatus = "CANCELLED"
			if err := tx.Model(&model.ShowSeat{}).Where("booking_id = ?", booking.ID).Updates(map[string]interface{}{"status": "available", "booking_id": nil}).Error; err != nil {
				return err
			}
		} else {
			booking.Status = "CONFIRMED"
			booking.PaymentStatus = "PAID"
			booking.PaidAt = &now
			if err := tx.Model(&model.ShowSeat{}).Where("booking_id = ?", booking.ID).Update("status", "booked").Error; err != nil {
				return err
			}
		}
		return tx.Save(&booking).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật thanh toán"})
		return
	}
	if apiErr != nil {
		c.JSON(statusCode, gin.H{"error": apiErr.Error()})
		return
	}
	detail, err := loadBookingResponse(db.DB, booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải thông tin vé"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"booking": detail, "message": map[bool]string{true: "Thanh toán mô phỏng thành công. Vé đã được xác nhận.", false: "Đã hủy đơn và trả ghế về trạng thái còn trống."}[action == "pay"]})
}

func GetMyBookings(c *gin.Context) {
	if !bookingDB(c) {
		return
	}
	userID, _ := c.Get("userID")
	var bookings []model.Booking
	if err := db.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải lịch sử đặt vé"})
		return
	}
	result := make([]bookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		item, err := loadBookingResponse(db.DB, booking)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải lịch sử đặt vé"})
			return
		}
		result = append(result, item)
	}
	c.JSON(http.StatusOK, result)
}

func GetMyBooking(c *gin.Context) {
	if !bookingDB(c) {
		return
	}
	userID, _ := c.Get("userID")
	var booking model.Booking
	if err := db.DB.Where("code = ? AND user_id = ?", c.Param("code"), userID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy mã đặt vé"})
		return
	}
	result, err := loadBookingResponse(db.DB, booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải thông tin vé"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func loadBookingResponse(database *gorm.DB, booking model.Booking) (bookingResponse, error) {
	var showtime model.Showtime
	if err := database.Preload("Cinema").Preload("Auditorium").First(&showtime, booking.ShowtimeID).Error; err != nil {
		return bookingResponse{}, err
	}
	var seats []model.ShowSeat
	if err := database.Model(&model.ShowSeat{}).Joins("JOIN booking_seats ON booking_seats.show_seat_id = show_seats.id").Where("booking_seats.booking_id = ?", booking.ID).Order("show_seats.row asc, show_seats.number asc").Find(&seats).Error; err != nil {
		return bookingResponse{}, err
	}
	codes := make([]string, 0, len(seats))
	for _, seat := range seats {
		codes = append(codes, seat.Code)
	}
	return bookingResponse{
		ID: booking.ID, Code: booking.Code, Showtime: showtime, SeatCodes: codes,
		TotalAmount: booking.TotalAmount, Status: booking.Status,
		PaymentStatus: booking.PaymentStatus, ExpiresAt: booking.ExpiresAt,
		PaidAt: booking.PaidAt, CreatedAt: booking.CreatedAt,
	}, nil
}

func makeTicketCode() (string, error) {
	bytes := make([]byte, 5)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "CS-" + strings.ToUpper(hex.EncodeToString(bytes)), nil
}
