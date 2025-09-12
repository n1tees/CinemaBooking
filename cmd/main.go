package main

import (
	"log"
	"time"

	"CinemaBooking/config"
	"CinemaBooking/pkg/routes"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "CinemaBooking/docs" // swagger docs
)

func main() {
	// Загружаем конфиг
	config.LoadEnv()

	// Немного подождать (если БД поднимается дольше)
	time.Sleep(5 * time.Second)

	// Подключаемся к БД
	// db.InitDB()
	// defer db.CloseDB()

	// Создаём роутер
	r := routes.SetupRouter()

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
