package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	// Replace legacy booking indexes. Earlier versions accidentally made either
	// booking_id or show_seat_id globally unique, preventing multi-seat bookings.
	if database.Migrator().HasTable(&model.BookingSeat{}) {
		for _, indexName := range []string{"idx_booking_seat", "idx_booking_seats_show_seat_id"} {
			// These names are fixed constants. DROP INDEX IF EXISTS also handles
			// partially migrated databases where GORM's index lookup is unreliable.
			if err := database.Exec("DROP INDEX IF EXISTS " + indexName).Error; err != nil {
				log.Printf("Unable to replace legacy booking index %s: %v", indexName, err)
				return nil
			}
		}
	}

	// Auto Migrate database tables
	err = database.AutoMigrate(
		&model.User{}, &model.Watchlist{}, &model.GenreDB{}, &model.MovieDB{},
		&model.Cinema{}, &model.Auditorium{}, &model.Showtime{}, &model.ShowSeat{},
		&model.Booking{}, &model.BookingSeat{},
	)
	if err != nil {
		log.Printf("⚠️ Lỗi AutoMigrate: %v", err)
	}

	DB = database
	if err := seedBookingCatalog(database); err != nil {
		log.Printf("Unable to seed demo cinemas and rooms: %v", err)
	}
	return database
}

func seedBookingCatalog(database *gorm.DB) error {
	catalog := []struct {
		cinema model.Cinema
		rooms  []model.Auditorium
	}{
		{model.Cinema{Name: "CineStream Landmark", Address: "720A Điện Biên Phủ, Bình Thạnh", City: "TP. Hồ Chí Minh"}, []model.Auditorium{{Name: "Phòng 1", RowCount: 6, SeatsPerRow: 10}, {Name: "Phòng 2", RowCount: 5, SeatsPerRow: 8}}},
		{model.Cinema{Name: "CineStream Crescent Mall", Address: "101 Tôn Dật Tiên, Quận 7", City: "TP. Hồ Chí Minh"}, []model.Auditorium{{Name: "Phòng 1", RowCount: 5, SeatsPerRow: 8}, {Name: "Phòng IMAX", RowCount: 7, SeatsPerRow: 12}}},
		{model.Cinema{Name: "CineStream Vincom Đồng Khởi", Address: "72 Lê Thánh Tôn, Quận 1", City: "TP. Hồ Chí Minh"}, []model.Auditorium{{Name: "Phòng 1", RowCount: 6, SeatsPerRow: 10}, {Name: "Phòng 2", RowCount: 4, SeatsPerRow: 8}}},
	}

	for _, item := range catalog {
		if err := database.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "name"}}, DoNothing: true}).Create(&item.cinema).Error; err != nil {
			return err
		}
		var cinema model.Cinema
		if err := database.Where("name = ?", item.cinema.Name).First(&cinema).Error; err != nil {
			return err
		}
		for _, room := range item.rooms {
			room.CinemaID = cinema.ID
			if err := database.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "cinema_id"}, {Name: "name"}}, DoNothing: true}).Create(&room).Error; err != nil {
				return err
			}
		}
		for roomNumber := 3; roomNumber <= 8; roomNumber++ {
			room := model.Auditorium{
				CinemaID: cinema.ID,
				Name:     fmt.Sprintf("Phòng %d", roomNumber), RowCount: 5 + roomNumber%3,
				SeatsPerRow: 8 + roomNumber%3,
			}
			if err := database.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "cinema_id"}, {Name: "name"}}, DoNothing: true}).Create(&room).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
