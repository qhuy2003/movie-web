package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const catalogListAPIURL = "https://vsmov.com/api/danh-sach"

func GetMoviesByCatalogSlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		slug = "phim-moi-cap-nhat"
	}

	params := url.Values{}
	if limit := c.Query("limit"); limit != "" {
		params.Set("limit", limit)
	}
	if page := c.Query("page"); page != "" {
		params.Set("page", page)
	}
	if year := c.Query("year"); year != "" {
		params.Set("year", year)
	}
	if movieType := c.Query("type"); movieType != "" {
		params.Set("type", movieType)
	}
	if status := c.Query("status"); status != "" {
		params.Set("status", status)
	}

	targetURL := fmt.Sprintf("%s/%s?%s", catalogListAPIURL, slug, params.Encode())

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Accept", "application/json")

	res, err := movieAPIClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể lấy danh sách phim từ catalog API", "detail": err.Error()})
		return
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
