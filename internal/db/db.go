package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"movie-web/internal/model"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "5433")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "movie_db")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("\n❌ LỖI KẾT NỐI POSTGRESQL [%s]: %v\n👉 Nguyên nhân: Sai mật khẩu user '%s' hoặc PostgreSQL chưa chạy trên port %s.\n👉 Hãy kiểm tra mật khẩu PostgreSQL trong file internal/db/db.go hoặc thiết lập biến môi trường DB_PASSWORD.\n", dbname, err, user, port)
		return nil
	}

	log.Printf("✅ Đã kết nối thành công cơ sở dữ liệu PostgreSQL [%s]", dbname)

	// Auto Migrate database tables
	err = database.AutoMigrate(&model.User{}, &model.Watchlist{}, &model.GenreDB{}, &model.MovieDB{})
	if err != nil {
		log.Printf("⚠️ Lỗi AutoMigrate: %v", err)
	}

	DB = database
	return database
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
