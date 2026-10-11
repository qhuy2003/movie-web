package model

import "time"

type Cinema struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

type Auditorium struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	CinemaID    uint   `gorm:"not null;uniqueIndex:idx_cinema_room" json:"cinema_id"`
	Name        string `gorm:"not null;uniqueIndex:idx_cinema_room" json:"name"`
	RowCount    int    `gorm:"not null" json:"row_count"`
	SeatsPerRow int    `gorm:"not null" json:"seats_per_row"`
}

type Showtime struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	MovieID       int        `gorm:"not null;index;uniqueIndex:idx_movie_cinema_start" json:"movie_id"`
	MovieTitle    string     `gorm:"not null" json:"movie_title"`
	PosterURL     string     `json:"poster_url"`
	CinemaID      uint       `gorm:"not null;uniqueIndex:idx_movie_cinema_start;uniqueIndex:idx_room_start" json:"cinema_id"`
	AuditoriumID  uint       `gorm:"not null;uniqueIndex:idx_room_start" json:"auditorium_id"`
	StartAt       time.Time  `gorm:"not null;uniqueIndex:idx_movie_cinema_start;uniqueIndex:idx_room_start" json:"start_at"`
	StandardPrice int        `gorm:"not null" json:"standard_price"`
	VIPPrice      int        `gorm:"not null" json:"vip_price"`
	Cinema        Cinema     `gorm:"foreignKey:CinemaID" json:"cinema"`
	Auditorium    Auditorium `gorm:"foreignKey:AuditoriumID" json:"auditorium"`
	CreatedAt     time.Time  `json:"created_at"`
}

type ShowSeat struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ShowtimeID uint   `gorm:"not null;uniqueIndex:idx_showtime_seat" json:"showtime_id"`
	Code       string `gorm:"not null;uniqueIndex:idx_showtime_seat" json:"code"`
	Row        string `gorm:"not null" json:"row"`
	Number     int    `gorm:"not null" json:"number"`
	Type       string `gorm:"not null;default:standard" json:"type"`
	Status     string `gorm:"not null;default:available;index" json:"status"`
	BookingID  *uint  `gorm:"index" json:"booking_id,omitempty"`
}

type Booking struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Code            string     `gorm:"uniqueIndex;not null" json:"code"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	ShowtimeID      uint       `gorm:"not null;index" json:"showtime_id"`
	TotalAmount     int        `gorm:"not null" json:"total_amount"`
	Status          string     `gorm:"not null;index" json:"status"`
	PaymentStatus   string     `gorm:"not null;index" json:"payment_status"`
	ZaloPayTransID  *string    `gorm:"uniqueIndex" json:"-"`
	ZaloPayOrderURL string     `gorm:"type:text" json:"-"`
	ExpiresAt       time.Time  `gorm:"index" json:"expires_at"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type BookingSeat struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	BookingID  uint     `gorm:"not null;uniqueIndex:idx_booking_seat" json:"booking_id"`
	ShowSeatID uint     `gorm:"not null;uniqueIndex:idx_booking_seat;index" json:"show_seat_id"`
	ShowSeat   ShowSeat `gorm:"foreignKey:ShowSeatID" json:"seat"`
}
