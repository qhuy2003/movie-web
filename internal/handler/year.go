package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

const yearsAPIURL = "https://vsmov.com/api/nam"

func GetYears(c *gin.Context) {
	proxyYearRequest(c, yearsAPIURL)
}

func GetMoviesByYear(c *gin.Context) {
	year := c.Param("year")
	yearNumber, err := strconv.Atoi(year)
	if err != nil || yearNumber < 1888 || yearNumber > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Năm phát hành không hợp lệ"})
		return
	}

	params := url.Values{}
	for _, key := range []string{"limit", "page", "type", "status"} {
		if value := c.Query(key); value != "" {
			params.Set(key, value)
		}
	}
	targetURL := fmt.Sprintf("%s/%s", yearsAPIURL, year)
	if query := params.Encode(); query != "" {
		targetURL += "?" + query
	}
	proxyYearRequest(c, targetURL)
}

func proxyYearRequest(c *gin.Context, targetURL string) {
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Accept", "application/json")
	res, err := movieAPIClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể gọi API năm phát hành", "detail": err.Error()})
		return
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		c.JSON(http.StatusBadGateway, gin.H{"error": "API năm phát hành trả về lỗi", "status": res.StatusCode})
		return
	}
	var payload interface{}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Dữ liệu năm phát hành không hợp lệ", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payload)
}
