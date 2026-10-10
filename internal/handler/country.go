package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const countriesAPIURL = "https://vsmov.com/api/quoc-gia"

type countriesAPIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Items []CountryItem `json:"items"`
	} `json:"data"`
}

type CountryItem struct {
	ID   int    `json:"_id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func fetchCountries() ([]CountryItem, error) {
	req, err := http.NewRequest(http.MethodGet, countriesAPIURL, nil)
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
		return nil, fmt.Errorf("countries API returned status %s", res.Status)
	}

	var payload countriesAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Data.Items, nil
}

func GetCountries(c *gin.Context) {
	countries, err := fetchCountries()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể lấy danh sách quốc gia", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, countries)
}

func GetMoviesByCountry(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug quốc gia không được để trống"})
		return
	}

	// Forward query params (limit, page, year, type, status)
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

	targetURL := fmt.Sprintf("%s/%s?%s", countriesAPIURL, slug, params.Encode())

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Accept", "application/json")

	res, err := movieAPIClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể gọi API quốc gia", "detail": err.Error()})
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
