package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"movie-web/internal/db"
	"movie-web/internal/model"
)

type movieReviewRequest struct {
	MovieTitle string `json:"movie_title" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Comment    string `json:"comment" binding:"required,max=1000"`
}

func GetMovieReviews(c *gin.Context) {
	database := db.DB
	if database == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Cơ sở dữ liệu hiện không khả dụng"})
		return
	}

	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil || movieID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã phim không hợp lệ"})
		return
	}

	reviews := make([]model.MovieReview, 0)
	if err := database.Where("movie_id = ?", movieID).Order("created_at desc").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải bình luận phim"})
		return
	}

	var summary struct {
		AverageRating float64 `json:"average_rating"`
		TotalReviews  int64   `json:"total_reviews"`
	}
	if err := database.Model(&model.MovieReview{}).
		Where("movie_id = ?", movieID).
		Select("COALESCE(AVG(rating), 0) AS average_rating, COUNT(*) AS total_reviews").
		Scan(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tính điểm đánh giá"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews":        reviews,
		"average_rating": summary.AverageRating,
		"total_reviews":  summary.TotalReviews,
	})
}

func GetMyMovieReviews(c *gin.Context) {
	database := db.DB
	if database == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Cơ sở dữ liệu hiện không khả dụng"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Vui lòng đăng nhập để xem đánh giá của bạn"})
		return
	}

	reviews := make([]model.MovieReview, 0)
	if err := database.Where("user_id = ?", userID).Order("updated_at desc").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải danh sách đánh giá của bạn"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func UpsertMovieReview(c *gin.Context) {
	database := db.DB
	if database == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Cơ sở dữ liệu hiện không khả dụng"})
		return
	}

	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil || movieID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã phim không hợp lệ"})
		return
	}

	var request movieReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng chọn điểm từ 1 đến 5 sao và nhập bình luận tối đa 1000 ký tự"})
		return
	}
	request.Comment = strings.TrimSpace(request.Comment)
	request.MovieTitle = strings.TrimSpace(request.MovieTitle)
	if request.Comment == "" || request.MovieTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên phim và nội dung bình luận không được để trống"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Vui lòng đăng nhập để gửi đánh giá"})
		return
	}
	var user model.User
	if err := database.Select("id", "username").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy tài khoản"})
		return
	}

	review := model.MovieReview{
		MovieID: movieID, UserID: user.ID, MovieTitle: request.MovieTitle,
		Username: user.Username, Rating: request.Rating, Comment: request.Comment,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	result := database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "movie_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"movie_title", "username", "rating", "comment", "updated_at"}),
	}).Create(&review)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu đánh giá phim"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"review": review, "message": "Đã lưu đánh giá của bạn"})
}
