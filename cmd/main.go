package main

import (
	"CinemaBooking/config"
	"CinemaBooking/pkg/db"
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	// "github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// swaggerFiles "github.com/swaggo/files"
	// "CinemaBooking/pkg/middleware"
	// "CinemaBooking/pkg/routes"
	// "CinemaBooking/pkg/services"
	// _ "CinemaBooking/docs"
)

func main() {
	config.LoadEnv()

	time.Sleep(5 * time.Second)

	db.InitDB()
	defer db.CloseDB()

	// запуск фонового обновления статусов
	go func() {
		for {
			services.CheckAndUpdateStatuses()
			time.Sleep(1 * time.Minute)
		}
	}()

	r := gin.Default()
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// Public routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.InitAuthRoutes(r)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())

	routes.InitRaceRoutes(api)
	routes.InitUserRoutes(api)
	routes.InitTrackRoutes(api)
	routes.InitKartodromRoutes(api)
	routes.InitPaymentRoutes(api)
	routes.InitBookingRoutes(api)
	routes.InitKartBookingRoutes(api)
	routes.InitKartRoutes(api)

	r.Run() // запускает на :8080
}
