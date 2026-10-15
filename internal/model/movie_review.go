package model

import "time"

type MovieReview struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MovieID    int       `gorm:"not null;uniqueIndex:idx_movie_review_user" json:"movie_id"`
	UserID     uint      `gorm:"not null;uniqueIndex:idx_movie_review_user" json:"user_id"`
	MovieTitle string    `gorm:"not null" json:"movie_title"`
	Username   string    `gorm:"not null" json:"username"`
	Rating     int       `gorm:"not null" json:"rating"`
	Comment    string    `gorm:"type:text;not null" json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
