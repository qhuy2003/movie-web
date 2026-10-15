package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"

	"movie-web/internal/db"
	"movie-web/internal/handler"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	loadEnvFile(".env")
	// Initialize PostgreSQL Database Connection
	log.Println("🔄 Đang khởi tạo kết nối PostgreSQL movie_db...")
	db.InitDB()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// Public Movie Routes
	movieRoutes := r.Group("/movies")
	{
		movieRoutes.GET("", handler.GetMovies)
		movieRoutes.GET("/:id/reviews", handler.GetMovieReviews)
		movieRoutes.GET("/:id", handler.GetMovieByID)
		movieRoutes.GET("/detail/:slug", handler.GetMovieDetailBySlug)
		movieRoutes.POST("", handler.CreateMovie)
	}

	// Public Catalog & Search Routes
	r.GET("/genres", handler.GetGenres)
	r.GET("/actors", handler.GetActors)
	r.GET("/countries", handler.GetCountries)
	r.GET("/countries/:slug", handler.GetMoviesByCountry)
	r.GET("/years", handler.GetYears)
	r.GET("/years/:year", handler.GetMoviesByYear)
	r.GET("/catalog/:slug", handler.GetMoviesByCatalogSlug)
	r.GET("/cinemas", handler.GetCinemas)
	r.GET("/showtimes", handler.GetShowtimes)
	r.GET("/showtimes/:id/seats", handler.GetShowtimeSeats)

	// Authentication Routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handler.Register)
		authRoutes.POST("/login", handler.Login)
		authRoutes.GET("/me", handler.AuthMiddleware(), handler.GetMe)
	}

	// Protected User Watchlist Routes
	apiRoutes := r.Group("/api")
	apiRoutes.Use(handler.AuthMiddleware())
	{
		apiRoutes.GET("/watchlist", handler.GetUserWatchlist)
		apiRoutes.GET("/reviews", handler.GetMyMovieReviews)
		apiRoutes.POST("/movies/:id/reviews", handler.UpsertMovieReview)
		apiRoutes.POST("/watchlist", handler.AddWatchlist)
		apiRoutes.DELETE("/watchlist/:movieId", handler.RemoveWatchlist)
		apiRoutes.GET("/bookings", handler.GetMyBookings)
		apiRoutes.GET("/bookings/:code", handler.GetMyBooking)
		apiRoutes.POST("/bookings", handler.CreateBooking)
		apiRoutes.POST("/bookings/:code/pay", handler.StartZaloPayPayment)
		apiRoutes.POST("/bookings/:code/verify-payment", handler.VerifyZaloPayPayment)
		apiRoutes.POST("/bookings/:code/cancel", handler.CancelBooking)
	}
	r.POST("/payments/zalopay/callback", handler.ZaloPayCallback)

	// Protected Admin Management Routes
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(handler.AuthMiddleware(), handler.AdminMiddleware())
	{
		// Users
		adminRoutes.GET("/users", handler.GetAdminUsers)
		adminRoutes.PUT("/users/:id/role", handler.UpdateUserRole)
		adminRoutes.DELETE("/users/:id", handler.DeleteUser)

		// Genres
		adminRoutes.GET("/genres", handler.GetAdminGenres)
		adminRoutes.POST("/genres", handler.CreateGenre)
		adminRoutes.PUT("/genres/:id", handler.UpdateGenre)
		adminRoutes.DELETE("/genres/:id", handler.DeleteGenre)
		adminRoutes.POST("/genres/import", handler.ImportGenresFromAPI)

		// Movies
		adminRoutes.GET("/movies", handler.GetAdminMovies)
		adminRoutes.POST("/movies", handler.CreateMovieAdmin)
		adminRoutes.POST("/movies/import", handler.ImportMoviesFromAPI)
		adminRoutes.PUT("/movies/:id", handler.UpdateMovieAdmin)
		adminRoutes.DELETE("/movies/:id", handler.DeleteMovieAdmin)
	}

	log.Println("🚀 Go Server đang chạy tại port :8080")
	r.Run(":8080")
}

func loadEnvFile(path string) {
	candidates := []string{path}
	if _, sourceFile, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates, filepath.Join(filepath.Dir(sourceFile), "..", "..", path))
	}
	for _, candidate := range candidates {
		loaded, err := readEnvFile(candidate)
		if err != nil {
			log.Printf("Unable to read environment file %s: %v", candidate, err)
		}
		if loaded && os.Getenv("ZALOPAY_APP_ID") != "" && os.Getenv("ZALOPAY_KEY1") != "" && os.Getenv("ZALOPAY_KEY2") != "" {
			log.Printf("Loaded local environment settings from %s", candidate)
			return
		}
	}
}

func readEnvFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	loaded := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), "\"'")
		if key == "" {
			continue
		}
		if current, exists := os.LookupEnv(key); exists && current != "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			log.Printf("Unable to set environment variable %s", key)
		} else {
			loaded = true
		}
	}
	if err := scanner.Err(); err != nil {
		return loaded, err
	}
	return loaded, nil
}
