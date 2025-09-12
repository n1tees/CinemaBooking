package dt

import (
	"CinemaBooking/pkg/models"
	"time"
)

// ErrorResponseDTO godoc
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"error"`
}

// RegisterDTI godoc
type RegisterDTI struct {
	FirstName      string    `json:"first_name" binding:"required"`
	SecondName     string    `json:"second_name"`
	Phone          string    `json:"phone" binding:"required"`
	Email          string    `json:"email" binding:"required,email"`
	BirthDay       string    `json:"birthday"`
	Login          string    `json:"login" binding:"required"`
	Password       string    `json:"password" binding:"required,min=6"`
	ParsedBirthDay time.Time `json:"-"`
}

// RegisterDTO godoc
type RegisterDTO struct {
	ID uint `json:"id"`
}

// LoginDTI godoc
type LoginDTI struct {
	Login    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginDTO godoc
type LoginDTO struct {
	JWT string `json:"jwt"`
}

// BonusBalanceDTO godoc
type BonusBalanceDTO struct {
	Balance float64 `json:"balance"`
}

// BonusHistoryDTO godoc
type BonusHistoryDTO struct {
	HistoryBalance []models.BonusHistory `json:"history_balance"`
}

// CreateBookingDTI godoc
type CreateBookingDTI struct {
	UserID    uint `json:"user_id" binding:"required"`
	SessionID uint `json:"session_id" binding:"required"`
	RowNum    uint `json:"row_num" binding:"required"`
	SeatNum   uint `json:"seat_num" binding:"required"`
	UseBonus  bool `json:"use_bonus" binding:"required"`
}

// BookingDTO godoc
type CreateBookingDTO struct {
	ID     uint                 `json:"status_id"`
	Status models.BookingStatus `json:"status"`
}

// GetFilmDTO godoc
type FilmDTO struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	AgeRating   uint       `json:"age_rating"`
	Duration    uint       `json:"duration"`
	ReleaseDate string     `json:"release_date"`
	Genres      []GenreDTO `json:"genres"`
}

// GetGenreDTO godoc
type GenreDTO struct {
	Name string `json:"name"`
}

// PaymentHistoryDTO godoc
type PaymentHistoryDTO struct {
	Amount    float64   `json:"amount"`
	Operation string    `json:"operation"`
	CreatedAt time.Time `json:"created_at"`
}

// PaymentDTO godoc
type PaymentDTO struct {
	Balance float64 `json:"balance"`
}

// RefillBalanceDTI godoc
type RefillBalanceDTI struct {
	Amount float64 `json:"amount" binding:"required"`
}

// PosterDTO godoc
type PosterDTO struct {
	FilmID uint   `json:"film_id"`
	URL    string `json:"url"`
}

// ProfileDTO godoc
type ProfileDTO struct {
	FirstName  string  `json:"first_name"`
	SecondName string  `json:"second_name"`
	Email      string  `json:"email"`
	Balance    float64 `json:"balance"`
	Bonus      float64 `json:"bonus"`
}

// ChangePasswordDTI godoc
type ChangePasswordDTI struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	RepeatNew   string `json:"repeat_new" binding:"required"`
}

// ChangePasswordDTI godoc
type ChangePasswordDTO struct {
	UserID uint   `json:"user_id"`
	Status string `json:"status"`
}

// UpdateProfileDTO godoc
type UpdateProfileDTO struct {
	Status string `json:"status"`
}

// ReviewDTI godoc
type ReviewDTI struct {
	FilmID uint `json:"film_id" binding:"required"`
}

// ReviewDTO godoc
type ReviewDTO struct {
	UserID  uint   `json:"user_id"`
	Rating  uint   `json:"rating"`
	Comment string `json:"comment"`
}

// CreateReviewDTI godoc
type CreateReviewDTI struct {
	Rating  uint   `json:"rating" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}

// CreateReviewDTO godoc
type CreateReviewDTO struct {
	Status string `json:"status"`
}

// FilmRatingDTO godoc
type FilmRatingDTO struct {
	Rating  *float64 `json:"rating,omitempty"`
	Message string   `json:"message,omitempty"`
}

// SessionDTO godoc
type SessionDTO struct {
	ID        uint      `json:"id"`
	FilmID    uint      `json:"film_id"`
	HallID    uint      `json:"hall_id"`
	StartTime time.Time `json:"start_time"`
	Price     float64   `json:"price"`
}

// SeatDTO godoc
type SeatDTO struct {
	Row   uint   `json:"row"`
	Seat  uint   `json:"seat"`
	State string `json:"state"` // free / taken
}

// ServAnswerDTO godoc
// для RefillBalanceDTO и CancelBookingDTO
type ServAnswerDTO struct {
	Answer string `json:"answer"`
}

// CreateFilmDTI godoc
type CreateFilmDTI struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Duration    uint   `json:"duration" binding:"required"`
	AgeRating   uint   `json:"age_rating" binding:"required"`
	ReleaseDate string `json:"release_date" binding:"required"`
	Genres      []uint `json:"genres"`
}

// CreateFilmDTO godoc
type CreateFilmDTO struct {
	ID uint `json:"id"`
}

// CreateGenreDTI godoc
type CreateGenreDTI struct {
	Name string `json:"name" binding:"required"`
}

// CreateGenreDTO godoc
type CreateGenreDTO struct {
	ID uint `json:"id"`
}

// AssignGenreDTI godoc
type AssignGenreDTI struct {
	GenreID uint `json:"genre_id" binding:"required"`
}

// CreateSessionDTI godoc
type CreateSessionDTI struct {
	FilmID uint      `json:"film_id" binding:"required"`
	HallID uint      `json:"hall_id" binding:"required"`
	Start  time.Time `json:"start" binding:"required"`
	Price  float64   `json:"price" binding:"required"`
}

// CreateSessionDTO godoc
type CreateSessionDTO struct {
	ID uint `json:"id"`
}

// CreatePosterDTI godoc
type CreatePosterDTI struct {
	FilmID   uint   `json:"film_id" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
}

// CreatePosterDTO godoc
type CreatePosterDTO struct {
	ID uint `json:"id"`
}
