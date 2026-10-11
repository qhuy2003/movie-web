package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"movie-web/internal/db"
	"movie-web/internal/model"
)

type zaloPayConfig struct {
	appID       string
	key1        string
	key2        string
	createURL   string
	queryURL    string
	callbackURL string
	frontendURL string
}

func getZaloPayConfig() (zaloPayConfig, error) {
	cfg := zaloPayConfig{
		appID:       strings.TrimSpace(os.Getenv("ZALOPAY_APP_ID")),
		key1:        strings.TrimSpace(os.Getenv("ZALOPAY_KEY1")),
		key2:        strings.TrimSpace(os.Getenv("ZALOPAY_KEY2")),
		createURL:   envOrDefault("ZALOPAY_CREATE_URL", "https://sb-openapi.zalopay.vn/v2/create"),
		queryURL:    envOrDefault("ZALOPAY_QUERY_URL", "https://sb-openapi.zalopay.vn/v2/query"),
		callbackURL: strings.TrimSpace(os.Getenv("ZALOPAY_CALLBACK_URL")),
		frontendURL: strings.TrimRight(envOrDefault("FRONTEND_URL", "http://localhost:5173"), "/"),
	}
	if cfg.appID == "" || cfg.key1 == "" || cfg.key2 == "" {
		return zaloPayConfig{}, fmt.Errorf("thiếu ZALOPAY_APP_ID, ZALOPAY_KEY1 hoặc ZALOPAY_KEY2")
	}
	if _, err := strconv.Atoi(cfg.appID); err != nil {
		return zaloPayConfig{}, fmt.Errorf("ZALOPAY_APP_ID phải là số")
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func hmacHex(key, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

type zaloPayCreateResponse struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
	OrderURL      string `json:"order_url"`
}

// StartZaloPayPayment creates a signed sandbox order for the current user's booking.
func StartZaloPayPayment(c *gin.Context) {
	if !bookingDB(c) {
		return
	}
	cfg, err := getZaloPayConfig()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ZaloPay Sandbox chưa được cấu hình trên backend."})
		return
	}
	userID, _ := c.Get("userID")
	var booking model.Booking
	if err := db.DB.Where("code = ? AND user_id = ?", c.Param("code"), userID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy mã đặt vé"})
		return
	}
	if booking.PaymentStatus == "PAID" {
		c.JSON(http.StatusConflict, gin.H{"error": "Đơn này đã thanh toán."})
		return
	}
	if booking.Status != "AWAITING_PAYMENT" || time.Now().After(booking.ExpiresAt) {
		c.JSON(http.StatusConflict, gin.H{"error": "Đơn đặt vé đã hết hạn hoặc không còn chờ thanh toán."})
		return
	}
	if booking.ZaloPayOrderURL != "" && booking.ZaloPayTransID != nil {
		if result, err := queryZaloPayOrder(cfg, *booking.ZaloPayTransID); err == nil && result.ReturnCode == 1 && result.Amount == int64(booking.TotalAmount) {
			if err := confirmPaidBooking(db.DB, booking.ID, result.ZPTransID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật trạng thái vé đã thanh toán."})
				return
			}
			if err := db.DB.First(&booking, booking.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải thông tin vé."})
				return
			}
			detail, err := loadBookingResponse(db.DB, booking)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải thông tin vé."})
				return
			}
			c.JSON(http.StatusOK, gin.H{"booking": detail})
			return
		}
		c.JSON(http.StatusOK, gin.H{"order_url": booking.ZaloPayOrderURL})
		return
	}

	appTransID := time.Now().In(cinemaLocation).Format("060102") + "_" + booking.Code
	appTime := time.Now().UnixMilli()
	item := "[]"
	redirectURL := cfg.frontendURL + "/booking?payment_return=1&booking=" + url.QueryEscape(booking.Code)
	embedData, _ := json.Marshal(map[string]interface{}{
		"redirecturl":              redirectURL,
		"preferred_payment_method": []string{},
	})
	embed := string(embedData)
	description := "Thanh toan ve CineStream " + booking.Code
	appUser := fmt.Sprintf("user-%v", userID)
	amount := strconv.Itoa(booking.TotalAmount)

	form := url.Values{}
	form.Set("app_id", cfg.appID)
	form.Set("app_user", appUser)
	form.Set("app_trans_id", appTransID)
	form.Set("app_time", strconv.FormatInt(appTime, 10))
	form.Set("amount", amount)
	form.Set("description", description)
	form.Set("item", item)
	form.Set("embed_data", embed)
	form.Set("bank_code", "")
	form.Set("expire_duration_seconds", "600")
	form.Set("redirect_url", redirectURL)
	if cfg.callbackURL != "" {
		form.Set("callback_url", cfg.callbackURL)
	}
	form.Set("mac", hmacHex(cfg.key1, strings.Join([]string{cfg.appID, appTransID, appUser, amount, strconv.FormatInt(appTime, 10), embed, item}, "|")))

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, cfg.createURL, strings.NewReader(form.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo yêu cầu thanh toán ZaloPay."})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không kết nối được ZaloPay Sandbox."})
		return
	}
	defer resp.Body.Close()
	var zp zaloPayCreateResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&zp); err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 || zp.ReturnCode != 1 || zp.OrderURL == "" {
		message := zp.ReturnMessage
		if message == "" {
			message = "ZaloPay không tạo được đơn thanh toán."
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": message})
		return
	}

	if err := db.DB.Model(&model.Booking{}).Where("id = ? AND status = ? AND payment_status = ?", booking.ID, "AWAITING_PAYMENT", "PENDING").Updates(map[string]interface{}{
		"zalo_pay_trans_id":  appTransID,
		"zalo_pay_order_url": zp.OrderURL,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Đã tạo lệnh ZaloPay nhưng không lưu được mã giao dịch."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_url": zp.OrderURL})
}

type zaloPayQueryResponse struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
	Amount        int64  `json:"amount"`
	ZPTransID     int64  `json:"zp_trans_id"`
}

func queryZaloPayOrder(cfg zaloPayConfig, appTransID string) (zaloPayQueryResponse, error) {
	form := url.Values{}
	form.Set("app_id", cfg.appID)
	form.Set("app_trans_id", appTransID)
	form.Set("mac", hmacHex(cfg.key1, cfg.appID+"|"+appTransID+"|"+cfg.key1))
	req, err := http.NewRequest(http.MethodPost, cfg.queryURL, strings.NewReader(form.Encode()))
	if err != nil {
		return zaloPayQueryResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return zaloPayQueryResponse{}, err
	}
	defer resp.Body.Close()
	var result zaloPayQueryResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
		return zaloPayQueryResponse{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return zaloPayQueryResponse{}, fmt.Errorf("ZaloPay query returned HTTP %d", resp.StatusCode)
	}
	return result, nil
}

// VerifyZaloPayPayment checks the order directly with ZaloPay; browser redirects are never trusted.
func VerifyZaloPayPayment(c *gin.Context) {
	if !bookingDB(c) {
		return
	}
	cfg, err := getZaloPayConfig()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ZaloPay Sandbox chưa được cấu hình trên backend."})
		return
	}
	userID, _ := c.Get("userID")
	var booking model.Booking
	if err := db.DB.Where("code = ? AND user_id = ?", c.Param("code"), userID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy mã đặt vé"})
		return
	}
	if booking.PaymentStatus != "PAID" && booking.ZaloPayTransID != nil {
		result, err := queryZaloPayOrder(cfg, *booking.ZaloPayTransID)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Không kết nối được ZaloPay để xác minh thanh toán."})
			return
		}
		if result.ReturnCode == 1 && result.Amount == int64(booking.TotalAmount) {
			if err := confirmPaidBooking(db.DB, booking.ID, result.ZPTransID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật trạng thái vé đã thanh toán."})
				return
			}
		}
	}
	if err := db.DB.First(&booking, booking.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải trạng thái đặt vé."})
		return
	}
	detail, err := loadBookingResponse(db.DB, booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải thông tin vé."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"booking": detail})
}

type zaloPayCallbackRequest struct {
	Data string `json:"data"`
	MAC  string `json:"mac"`
}

func ZaloPayCallback(c *gin.Context) {
	if db.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"return_code": 0, "return_message": "database unavailable"})
		return
	}
	cfg, err := getZaloPayConfig()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"return_code": 0, "return_message": "ZaloPay callback is not configured"})
		return
	}
	var callback zaloPayCallbackRequest
	if err := c.ShouldBindJSON(&callback); err != nil || callback.Data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"return_code": 0, "return_message": "invalid callback"})
		return
	}
	want, err := hex.DecodeString(hmacHex(cfg.key2, callback.Data))
	got, decodeErr := hex.DecodeString(callback.MAC)
	if err != nil || decodeErr != nil || !hmac.Equal(want, got) {
		c.JSON(http.StatusOK, gin.H{"return_code": -1, "return_message": "mac not equal"})
		return
	}
	var data struct {
		AppID      int64  `json:"app_id"`
		AppTransID string `json:"app_trans_id"`
		Amount     int64  `json:"amount"`
		ZPTransID  int64  `json:"zp_trans_id"`
	}
	if err := json.Unmarshal([]byte(callback.Data), &data); err != nil || strconv.FormatInt(data.AppID, 10) != cfg.appID {
		c.JSON(http.StatusOK, gin.H{"return_code": -1, "return_message": "invalid transaction data"})
		return
	}
	var booking model.Booking
	if err := db.DB.Where("zalo_pay_trans_id = ?", data.AppTransID).First(&booking).Error; err != nil || int64(booking.TotalAmount) != data.Amount {
		c.JSON(http.StatusOK, gin.H{"return_code": -1, "return_message": "transaction not found or amount mismatch"})
		return
	}
	if err := confirmPaidBooking(db.DB, booking.ID, data.ZPTransID); err != nil {
		c.JSON(http.StatusOK, gin.H{"return_code": 0, "return_message": "could not update booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"return_code": 1, "return_message": "success"})
}

func confirmPaidBooking(database *gorm.DB, bookingID uint, zpTransID int64) error {
	return database.Transaction(func(tx *gorm.DB) error {
		var booking model.Booking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&booking, bookingID).Error; err != nil {
			return err
		}
		if booking.PaymentStatus == "PAID" {
			return nil
		}
		if booking.Status != "AWAITING_PAYMENT" || time.Now().After(booking.ExpiresAt) {
			return fmt.Errorf("booking is no longer payable")
		}
		now := time.Now()
		booking.Status = "CONFIRMED"
		booking.PaymentStatus = "PAID"
		booking.PaidAt = &now
		if err := tx.Model(&model.ShowSeat{}).Where("booking_id = ?", booking.ID).Update("status", "booked").Error; err != nil {
			return err
		}
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}
		return nil
	})
}
