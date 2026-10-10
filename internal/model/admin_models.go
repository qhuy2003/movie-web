package model

import (
	"time"
)

type GenreDB struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type MovieDB struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	OriginName  string    `json:"origin_name"`
	Slug        string    `gorm:"index" json:"slug"`
	Year        int       `json:"year"`
	Rating      float64   `json:"rating"`
	Quality     string    `json:"quality"`
	Duration    string    `json:"duration"`
	Director    string    `json:"director"`
	Cast        string    `json:"cast"`        // Comma separated or text
	Genres      string    `json:"genres"`      // Comma separated or text
	PosterURL   string    `json:"poster_url"`
	ThumbURL    string    `json:"thumb_url"`
	Description string    `json:"description"`
	VideoURL    string    `json:"video_url"`
	CreatedAt   time.Time `json:"created_at"`
}
