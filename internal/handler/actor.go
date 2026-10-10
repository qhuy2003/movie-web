package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const actorsAPIURL = "https://vsmov.com/api/dien-vien"

type actorsAPIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Items []ActorItem `json:"items"`
	} `json:"data"`
}

type ActorItem struct {
	ID       int    `json:"_id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ThumbURL string `json:"thumb_url"`
}

func fetchActors() ([]ActorItem, error) {
	req, err := http.NewRequest(http.MethodGet, actorsAPIURL, nil)
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
		return nil, fmt.Errorf("actors API returned status %s", res.Status)
	}

	var payload actorsAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Data.Items, nil
}

func GetActors(c *gin.Context) {
	searchQuery := strings.ToLower(strings.TrimSpace(c.Query("q")))

	actors, err := fetchActors()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể lấy danh sách diễn viên", "detail": err.Error()})
		return
	}

	if searchQuery != "" {
		filtered := make([]ActorItem, 0)
		for _, a := range actors {
			if strings.Contains(strings.ToLower(a.Name), searchQuery) || strings.Contains(strings.ToLower(a.Slug), searchQuery) {
				filtered = append(filtered, a)
			}
		}
		c.JSON(http.StatusOK, filtered)
		return
	}

	// Limit response if list is very large
	if len(actors) > 100 {
		c.JSON(http.StatusOK, actors[:100])
		return
	}

	c.JSON(http.StatusOK, actors)
}
