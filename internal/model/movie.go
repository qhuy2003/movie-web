package model

type Movie struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Year       int      `json:"year"`
	Rating     float64  `json:"rating"`
	Genres     []string `json:"genres"`
	OriginName string   `json:"origin_name,omitempty"`
	PosterURL  string   `json:"poster_url,omitempty"`
	ThumbURL   string   `json:"thumb_url,omitempty"`
	Slug       string   `json:"slug,omitempty"`
}
