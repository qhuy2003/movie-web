package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"movie-web/internal/db"
	"movie-web/internal/model"
)

// ==================== USER MANAGEMENT ====================

func GetAdminUsers(c *gin.Context) {
	var users []model.User
	if db.DB != nil {
		db.DB.Order("created_at desc").Find(&users)
	}
	c.JSON(http.StatusOK, users)
}

func UpdateUserRole(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID người dùng không hợp lệ"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	role := strings.ToLower(strings.TrimSpace(req.Role))
	if role != "admin" && role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role phải là 'admin' hoặc 'user'"})
		return
	}

	if db.DB != nil {
		if err := db.DB.Model(&model.User{}).Where("id = ?", userID).Update("role", role).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật role"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã cập nhật quyền thành công", "role": role})
}

func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID người dùng không hợp lệ"})
		return
	}

	if db.DB != nil {
		if err := db.DB.Delete(&model.User{}, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa người dùng"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa người dùng thành công"})
}

// ==================== GENRE MANAGEMENT ====================

func GetAdminGenres(c *gin.Context) {
	var genres []model.GenreDB
	if db.DB != nil {
		db.DB.Order("id asc").Find(&genres)
	}
	c.JSON(http.StatusOK, genres)
}

func CreateGenre(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên thể loại là bắt buộc"})
		return
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Name), " ", "-"))
	}

	newGenre := model.GenreDB{
		Name:      strings.TrimSpace(req.Name),
		Slug:      slug,
		CreatedAt: time.Now(),
	}

	if db.DB != nil {
		if err := db.DB.Create(&newGenre).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Thể loại này đã tồn tại hoặc không thể tạo"})
			return
		}
	}

	c.JSON(http.StatusCreated, newGenre)
}

func UpdateGenre(c *gin.Context) {
	idParam := c.Param("id")
	genreID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID thể loại không hợp lệ"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Name), " ", "-"))
	}

	if db.DB != nil {
		if err := db.DB.Model(&model.GenreDB{}).Where("id = ?", genreID).Updates(map[string]interface{}{
			"name": strings.TrimSpace(req.Name),
			"slug": slug,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật thể loại"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã cập nhật thể loại thành công"})
}

func DeleteGenre(c *gin.Context) {
	idParam := c.Param("id")
	genreID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID thể loại không hợp lệ"})
		return
	}

	if db.DB != nil {
		if err := db.DB.Delete(&model.GenreDB{}, genreID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa thể loại"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa thể loại thành công"})
}

func ImportGenresFromAPI(c *gin.Context) {
	if db.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "PostgreSQL chưa kết nối; không thể đồng bộ thể loại"})
		return
	}

	apiGenres, err := fetchGenres()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể lấy danh sách thể loại từ API", "detail": err.Error()})
		return
	}

	genres := make([]model.GenreDB, 0, len(apiGenres))
	for _, genre := range apiGenres {
		genres = append(genres, model.GenreDB{Name: genre.Name, Slug: genre.Slug, CreatedAt: time.Now()})
	}
	result := db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&genres)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu danh sách thể loại vào PostgreSQL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Đã đồng bộ %d thể loại; thêm mới %d thể loại.", len(apiGenres), result.RowsAffected),
		"count":   result.RowsAffected,
	})
}

// ==================== MOVIE MANAGEMENT ====================

func GetAdminMovies(c *gin.Context) {
	var dbMovies []model.MovieDB
	if db.DB != nil {
		db.DB.Order("created_at desc").Find(&dbMovies)
	}
	c.JSON(http.StatusOK, dbMovies)
}

func CreateMovieAdmin(c *gin.Context) {
	var movie model.MovieDB
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu phim không hợp lệ: " + err.Error()})
		return
	}

	if movie.Slug == "" {
		movie.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(movie.Title), " ", "-"))
	}
	movie.CreatedAt = time.Now()

	if db.DB != nil {
		if err := db.DB.Create(&movie).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu phim vào PostgreSQL"})
			return
		}
	}

	c.JSON(http.StatusCreated, movie)
}

func UpdateMovieAdmin(c *gin.Context) {
	idParam := c.Param("id")
	movieID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID phim không hợp lệ"})
		return
	}

	var req model.MovieDB
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if db.DB != nil {
		if err := db.DB.Model(&model.MovieDB{}).Where("id = ?", movieID).Updates(map[string]interface{}{
			"title":       req.Title,
			"origin_name": req.OriginName,
			"slug":        req.Slug,
			"year":        req.Year,
			"rating":      req.Rating,
			"quality":     req.Quality,
			"duration":    req.Duration,
			"director":    req.Director,
			"cast":        req.Cast,
			"genres":      req.Genres,
			"poster_url":  req.PosterURL,
			"thumb_url":   req.ThumbURL,
			"description": req.Description,
			"video_url":   req.VideoURL,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật phim"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã cập nhật thông tin phim thành công"})
}

func DeleteMovieAdmin(c *gin.Context) {
	idParam := c.Param("id")
	movieID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID phim không hợp lệ"})
		return
	}

	if db.DB != nil {
		if err := db.DB.Delete(&model.MovieDB{}, movieID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa phim"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa phim thành công"})
}

func ImportMoviesFromAPI(c *gin.Context) {
	apiMovies, err := fetchMovies(1)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không thể lấy danh sách phim từ API", "detail": err.Error()})
		return
	}

	importedCount := 0
	if db.DB != nil {
		for _, m := range apiMovies {
			var existing model.MovieDB
			if err := db.DB.Where("title = ? OR slug = ?", m.Title, m.Slug).First(&existing).Error; err != nil {
				newMovie := model.MovieDB{
					Title:       m.Title,
					OriginName:  m.OriginName,
					Slug:        m.Slug,
					Year:        m.Year,
					Rating:      m.Rating,
					Quality:     "HD",
					PosterURL:   m.PosterURL,
					ThumbURL:    m.ThumbURL,
					Description: "Phim được tự động đồng bộ từ API vsmov.",
					CreatedAt:   time.Now(),
				}
				if err := db.DB.Create(&newMovie).Error; err == nil {
					importedCount++
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Đã nhập thành công %d bộ phim từ API vào PostgreSQL!", importedCount),
		"count":   importedCount,
	})
}
