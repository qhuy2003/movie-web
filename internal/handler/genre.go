package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const genreAPIURL = "https://vsmov.com/api/the-loai"

type genreAPIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Items []struct {
			ID   int    `json:"_id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"items"`
	} `json:"data"`
}

type GenreItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func fetchGenres() ([]GenreItem, error) {
	req, err := http.NewRequest(http.MethodGet, genreAPIURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; CinemaStream/1.0)")

	res, err := movieAPIClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("genre API returned status %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var payload genreAPIResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("genre API returned invalid JSON (%s): %w", res.Header.Get("Content-Type"), err)
	}
	if len(payload.Data.Items) == 0 {
		return nil, fmt.Errorf("genre API returned no items (%s)", res.Header.Get("Content-Type"))
	}

	result := make([]GenreItem, 0, len(payload.Data.Items))
	for _, item := range payload.Data.Items {
		result = append(result, GenreItem{
			ID:   item.ID,
			Name: item.Name,
			Slug: item.Slug,
		})
	}
	return result, nil
}

func GetGenres(c *gin.Context) {
	genres, err := fetchGenres()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể lấy danh sách thể loại từ vsmov API", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, genres)
}
