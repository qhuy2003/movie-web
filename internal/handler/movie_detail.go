package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const movieDetailAPIURL = "https://vsmov.com/api/phim/"

type MovieDetailResponse struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Movie  struct {
		ID         int      `json:"_id"`
		Name       string   `json:"name"`
		OriginName string   `json:"origin_name"`
		Slug       string   `json:"slug"`
		Content    string   `json:"content"`
		Type       string   `json:"type"`
		Status     string   `json:"status"`
		PosterURL  string   `json:"poster_url"`
		ThumbURL   string   `json:"thumb_url"`
		TrailerURL string   `json:"trailer_url"`
		Time       string   `json:"time"`
		Quality    string   `json:"quality"`
		Lang       string   `json:"lang"`
		Year       int      `json:"year"`
		Actor      []string `json:"actor"`
		Director   []string `json:"director"`
		Category   []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"category"`
		Country []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"country"`
		TMDB struct {
			VoteAverage string `json:"vote_average"`
		} `json:"tmdb"`
	} `json:"movie"`
	Episodes []interface{} `json:"episodes"`
}

func fetchMovieDetailBySlug(slug string) (*MovieDetailResponse, error) {
	url := fmt.Sprintf("%s%s", movieDetailAPIURL, slug)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	res, err := movieAPIClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("movie detail API returned status %s", res.Status)
	}

	var payload MovieDetailResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func GetMovieDetailBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug không được để trống"})
		return
	}

	detail, err := fetchMovieDetailBySlug(slug)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể lấy thông tin chi tiết phim", "detail": err.Error()})
		return
	}

	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.JSON(http.StatusOK, detail)
}
