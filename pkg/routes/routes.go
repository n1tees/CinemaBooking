package routes

import (
	"CinemaBooking/pkg/handlers"
	"CinemaBooking/pkg/middleware"

	adminHandlers "CinemaBooking/pkg/handlers/admin"
	userHandlers "CinemaBooking/pkg/handlers/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//  AUTH
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterHandler)
		auth.POST("/login", handlers.LoginHandler)
	}

	//  PROFILE
	profile := r.Group("/profile", middleware.AuthRequired())
	{
		profile.GET("", userHandlers.GetUserInfoHandler)
		profile.PATCH("", userHandlers.UpdateProfileHandler)
		profile.PATCH("/password", userHandlers.ChangePasswordHandler)
	}

	//  WALLET
	wallet := r.Group("/wallet", middleware.AuthRequired())
	{
		wallet.GET("/balance", userHandlers.GetBalanceHandler)
		wallet.GET("/payments", userHandlers.GetMyPaymentsHandler)
		wallet.POST("/refill", userHandlers.RefillMyBalanceHandler)
	}

	//  BONUS
	bonus := r.Group("/bonus", middleware.AuthRequired())
	{
		bonus.GET("/balance", userHandlers.GetBonusBalanceHandler)
		bonus.GET("/history", userHandlers.GetBonusHistoryHandler)
	}

	//  FILMS
	films := r.Group("/films")
	{
		films.GET("", userHandlers.GetAllFilmsHandler)
		films.GET("/:id", userHandlers.GetFilmHandler)
		films.GET("/:id/genres", userHandlers.GetGenresByFilmHandler)
		films.GET("/:id/reviews", userHandlers.GetReviewsByFilmHandler)
		films.GET("/:id/rating", userHandlers.GetFilmRatingHandler)

		// отзывы (только авторизованный)
		filmsAuth := films.Group("/:id/reviews", middleware.AuthRequired())
		{
			filmsAuth.POST("", userHandlers.AddReviewHandler)
		}
	}

	//  GENRES (публичные)
	r.GET("/genres", userHandlers.GetAllGenresHandler)

	//  POSTERS
	r.GET("/posters", userHandlers.GetAllPostersHandler)

	//  SESSIONS
	sessions := r.Group("/sessions")
	{
		sessions.GET("", userHandlers.GetAllSessionsHandler)
		sessions.GET("/film/:id", userHandlers.GetSessionsByFilmHandler)
		sessions.GET("/:id/seats", userHandlers.GetAvailableSeatsHandler)
	}

	//  BOOKINGS
	bookings := r.Group("/bookings", middleware.AuthRequired())
	{
		bookings.POST("", userHandlers.CreateBookingHandler)
		bookings.DELETE("/:id", adminHandlers.CancelBookingHandler)
		// можно добавить GET /bookings для истории броней
	}

	//  ADMIN
	admin := r.Group("/admin", middleware.AuthRequired(), middleware.AdminOnly())
	{
		// фильмы
		admin.POST("/films", adminHandlers.CreateFilmHandler)
		admin.PATCH("/films/:id", adminHandlers.UpdateFilmHandler)
		admin.DELETE("/films/:id", adminHandlers.DeleteFilmHandler)

		// жанры
		admin.POST("/genres", adminHandlers.CreateGenreHandler)
		admin.PATCH("/genres/:id", adminHandlers.UpdateGenreHandler)
		admin.DELETE("/genres/:id", adminHandlers.DeleteGenreHandler)
		admin.POST("/films/:id/genres", adminHandlers.AssignGenreToFilmHandler)
		admin.DELETE("/films/:id/genres/:genre_id", adminHandlers.RemoveGenreFromFilmHandler)

		// сеансы
		admin.POST("/sessions", adminHandlers.CreateSessionHandler)
		admin.PATCH("/sessions/:id", adminHandlers.UpdateSessionHandler)
		admin.DELETE("/sessions/:id", adminHandlers.DeleteSessionHandler)

		// афиши
		admin.POST("/posters", adminHandlers.CreatePosterHandler)
		admin.PATCH("/posters/:id", adminHandlers.UpdatePosterHandler)
		admin.DELETE("/posters/:id", adminHandlers.DeletePosterHandler)

		// модерация отзывов
		admin.PATCH("/reviews/:id/approve", adminHandlers.ApproveReviewHandler)
		admin.PATCH("/reviews/:id/reject", adminHandlers.RejectReviewHandler)
	}

	return r
}
