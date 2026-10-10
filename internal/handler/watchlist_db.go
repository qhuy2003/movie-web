package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"movie-web/internal/db"
	"movie-web/internal/model"
)

func GetUserWatchlist(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	var items []model.Watchlist
	if db.DB != nil {
		db.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&items)
	}

	c.JSON(http.StatusOK, items)
}

func AddWatchlist(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	var req struct {
		MovieID   int     `json:"movie_id" binding:"required"`
		Title     string  `json:"title"`
		PosterURL string  `json:"poster_url"`
		Rating    float64 `json:"rating"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	item := model.Watchlist{
		UserID:    userID,
		MovieID:   req.MovieID,
		Title:     req.Title,
		PosterURL: req.PosterURL,
		Rating:    req.Rating,
		CreatedAt: time.Now(),
	}

	if db.DB != nil {
		var count int64
		db.DB.Model(&model.Watchlist{}).Where("user_id = ? AND movie_id = ?", userID, req.MovieID).Count(&count)
		if count == 0 {
			db.DB.Create(&item)
		}
	}

	c.JSON(http.StatusCreated, item)
}

func RemoveWatchlist(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	movieIDParam := c.Param("movieId")
	movieID, err := strconv.Atoi(movieIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MovieID không hợp lệ"})
		return
	}

	if db.DB != nil {
		db.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).Delete(&model.Watchlist{})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa khỏi danh sách yêu thích"})
}
