package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"movie-web/internal/db"
	"movie-web/internal/model"
)

var jwtSecret = []byte("cinestream_jwt_secret_key_2026")

// In-Memory fallback store if DB is disconnected
var memoryUsers = make(map[string]model.User)

func GenerateJWT(user model.User) (string, error) {
	role := user.Role
	if role == "" {
		role = "user"
	}
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu đăng ký không hợp lệ: " + err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mã hóa mật khẩu"})
		return
	}

	newUser := model.User{
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Username:  strings.TrimSpace(req.Username),
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: time.Now(),
	}

	if db.DB != nil {
		var existingUser model.User
		if err := db.DB.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email này đã được sử dụng"})
			return
		}

		// Self promotion if first user or admin email
		var count int64
		db.DB.Model(&model.User{}).Count(&count)
		if count == 0 {
			newUser.Role = "admin"
		}

		if err := db.DB.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu tài khoản vào cơ sở dữ liệu"})
			return
		}
	} else {
		// Fallback memory
		if _, exists := memoryUsers[newUser.Email]; exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Email này đã được sử dụng"})
			return
		}
		if len(memoryUsers) == 0 {
			newUser.Role = "admin"
		}
		newUser.ID = uint(len(memoryUsers) + 1)
		memoryUsers[newUser.Email] = newUser
	}

	token, err := GenerateJWT(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
		return
	}

	c.JSON(http.StatusCreated, model.AuthResponse{
		Token: token,
		User:  newUser,
	})
}

func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email hoặc mật khẩu không đúng định dạng"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	var user model.User

	if db.DB != nil {
		if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không chính xác"})
			return
		}
	} else {
		// Memory fallback
		var ok bool
		user, ok = memoryUsers[email]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không chính xác"})
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không chính xác"})
		return
	}

	token, err := GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
		return
	}

	c.JSON(http.StatusOK, model.AuthResponse{
		Token: token,
		User:  user,
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Yêu cầu đăng nhập"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format Token không hợp lệ"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token đã hết hạn hoặc không hợp lệ"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Lỗi đọc claims"})
			c.Abort()
			return
		}

		userIDClaim, ok := claims["user_id"].(float64)
		if !ok || userIDClaim <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token khÃ´ng cÃ³ user_id há»£p lá»‡"})
			c.Abort()
			return
		}
		userID := uint(userIDClaim)
		c.Set("userID", userID)

		// Read the current role instead of trusting a potentially stale JWT role.
		role := "user"
		if db.DB != nil {
			var user model.User
			if err := db.DB.Select("role").First(&user, userID).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "KhÃ´ng tÃ¬m tháº¥y tÃ i khoáº£n"})
				c.Abort()
				return
			}
			role = user.Role
		} else {
			for _, user := range memoryUsers {
				if user.ID == userID {
					role = user.Role
					break
				}
			}
		}
		c.Set("userRole", role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("userRole")
		if !exists || roleVal.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Yêu cầu quyền Quản Trị Viên (Admin) để thực hiện thao tác này"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Chưa đăng nhập"})
		return
	}

	userID := userIDVal.(uint)
	var user model.User

	if db.DB != nil {
		if err := db.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy người dùng"})
			return
		}
	} else {
		for _, u := range memoryUsers {
			if u.ID == userID {
				user = u
				break
			}
		}
	}

	c.JSON(http.StatusOK, user)
}
