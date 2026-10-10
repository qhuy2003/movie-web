package main

import (
	"log"

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
	// Initialize PostgreSQL Database Connection
	log.Println("🔄 Đang khởi tạo kết nối PostgreSQL movie_db...")
	db.InitDB()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// Public Movie Routes
	movieRoutes := r.Group("/movies")
	{
		movieRoutes.GET("", handler.GetMovies)
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
		apiRoutes.POST("/watchlist", handler.AddWatchlist)
		apiRoutes.DELETE("/watchlist/:movieId", handler.RemoveWatchlist)
	}

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
