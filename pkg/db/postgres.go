package db

import (
	"CinemaBooking/config"
	"CinemaBooking/pkg/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.GetDBConnString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	log.Println("Успешное подключение к базе данных, подготовка к выполнению миграции")

	if err := Migrate(db); err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	DB = db
	log.Println("Пропуск миграций, не забыть вернуть")
	log.Println("Успешное выполнение миграций")
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Ошибка при получении SQL-соединения: %v", err)
		return
	}
	sqlDB.Close()
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.AuthCredential{},
		&models.User{},
		&models.Profile{},

		&models.Cinema{},
		&models.HallType{},
		&models.CinemaHall{},
		&models.HallSeat{},

		&models.Genre{},
		&models.Film{},
		&models.Poster{},
		&models.FilmGenre{},
		&models.Review{},

		&models.Session{},
		&models.Booking{},

		&models.PaymentHistory{},
		&models.BonusHistory{},
	)
}
