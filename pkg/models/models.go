package models

import (
	"time"

	"gorm.io/gorm"
)

// ENUM’ы
type BookingStatus string

const (
	BookingReserved BookingStatus = "reserved"
	BookingPaid     BookingStatus = "paid"
	BookingCanceled BookingStatus = "canceled"
)

type PaymentOperation string

const (
	PaymentDeposit PaymentOperation = "deposit"
	PaymentSpend   PaymentOperation = "spend"
)

type BonusOperation string

const (
	BonusEarn   BonusOperation = "earn"
	BonusRedeem BonusOperation = "redeem"
)

// Пользователи и авторизация
type AuthCredential struct {
	gorm.Model
	Login        string `gorm:"type:varchar(50);unique;not null"`
	PasswordHash []byte `gorm:"not null"`
}

type Profile struct {
	gorm.Model
	FirstName  string    `gorm:"type:varchar(50);not null"`
	SecondName string    `gorm:"type:varchar(50)"`
	Phone      string    `gorm:"type:varchar(11);unique;not null"`
	Email      string    `gorm:"type:varchar(50)"`
	BirthDay   time.Time `gorm:"type:date"`
	Balance    float64   `gorm:"type:numeric(12,2);default:0"`
	Bonus      int       `gorm:"default:0"`
}

type User struct {
	gorm.Model
	AuthID    uint
	Auth      AuthCredential
	ProfileID uint
	Profile   Profile
}

// Кинотеатры и залы
type Cinema struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null"`
	Location string `gorm:"type:varchar(100)"`
	Phone    string `gorm:"type:varchar(11)"`
	Email    string `gorm:"type:varchar(50)"`
}

type HallType struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null"`
	Desc string `gorm:"type:varchar(50)"`
}

type CinemaHall struct {
	gorm.Model
	CinemaID   uint
	Cinema     Cinema
	HallTypeID uint
	HallType   HallType
	Name       string `gorm:"type:varchar(50);not null"`
	Capacity   uint
}

type HallSeat struct {
	gorm.Model
	HallID uint
	Hall   CinemaHall
	Row    uint
	Seat   uint
}

// Фильмы
type Genre struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null"`
	Desc string `gorm:"type:varchar(50)"`
}

type Film struct {
	gorm.Model
	Title       string `gorm:"type:varchar(50);not null"`
	Desc        string `gorm:"type:varchar(100)"`
	Duration    uint   // минуты
	AgeRating   uint
	ReleaseDate time.Time
	Genres      []Genre `gorm:"many2many:film_genres;"`
	Posters     []Poster
}

type Poster struct {
	gorm.Model
	FilmID   uint
	Film     Film
	ImageURL string `gorm:"type:varchar(100)"`
}

// Сеансы и бронирование
type Session struct {
	gorm.Model
	FilmID    uint
	Film      Film
	HallID    uint
	Hall      CinemaHall
	StartTime time.Time
	Price     float64 `gorm:"type:numeric(12,2)"`
	Bookings  []Booking
}

type Booking struct {
	gorm.Model
	SeatID     uint
	Seat       HallSeat
	CustomerID uint
	Customer   User
	SessionID  uint
	Session    Session
	Status     BookingStatus `gorm:"type:varchar(20);not null"`
}

// Финансы
type PaymentHistory struct {
	gorm.Model
	UserID    uint
	User      User
	Amount    float64          `gorm:"type:numeric(12,2)"`
	Operation PaymentOperation `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
}

type BonusHistory struct {
	gorm.Model
	UserID    uint
	User      User
	Amount    int
	Operation BonusOperation `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
}

type FilmGenre struct {
	FilmID  uint `gorm:"primaryKey"`
	GenreID uint `gorm:"primaryKey"`
}
