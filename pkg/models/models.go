package models

import (
	"time"

	"gorm.io/datatypes"
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

type UserType string

const (
	Customer UserType = "customer"
	Admin    UserType = "admin"
)

type ReviewStatus string

const (
	ReviewPending  ReviewStatus = "pending"  // отправлен на проверку
	ReviewApproved ReviewStatus = "approved" // опубликован
	ReviewRejected ReviewStatus = "rejected" // отменен
)

// Пользователи и авторизация
type AuthCredential struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Login        string `gorm:"type:varchar(50);unique;not null"`
	PasswordHash []byte `gorm:"not null"`
}

type Profile struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	FirstName  string    `gorm:"type:varchar(50);not null"`
	SecondName string    `gorm:"type:varchar(50)"`
	Phone      string    `gorm:"type:varchar(11);unique;not null"`
	Email      string    `gorm:"type:varchar(50)"`
	BirthDay   time.Time `gorm:"type:date"`
	Balance    float64   `gorm:"type:numeric(12,2);not null"`
	Bonus      float64   `gorm:"type:numeric(12,2);not null"`
}

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	AuthID    uint
	Auth      AuthCredential
	ProfileID uint
	Profile   Profile
	UserType  UserType `gorm:"type:varchar(25);not null"`
}

// Кинотеатры и залы
type Cinema struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Name     string `gorm:"type:varchar(100);not null"`
	Location string `gorm:"type:varchar(100)"`
	Phone    string `gorm:"type:varchar(11)"`
	Email    string `gorm:"type:varchar(50)"`
}

type HallType struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Name string `gorm:"type:varchar(50);not null"`
	Desc string `gorm:"type:varchar(50)"`
}

type CinemaHall struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	CinemaID   uint
	Cinema     Cinema
	HallTypeID uint
	HallType   HallType
	Name       string `gorm:"type:varchar(50);not null"`
	Capacity   uint
	Structure  datatypes.JSON `gorm:"type:jsonb"`
}

// Фильмы
type Genre struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Name string `gorm:"type:varchar(50);not null"`
	Desc string `gorm:"type:varchar(50)"`
}

type Film struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Title       string `gorm:"type:varchar(50);not null"`
	Desc        string `gorm:"type:varchar(100)"`
	Duration    uint   // минуты
	AgeRating   uint
	ReleaseDate time.Time
	Genres      []Genre `gorm:"many2many:film_genres;"`
	Posters     []Poster
}

type Poster struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	FilmID   uint `gorm:"not null;constraint:OnDelete:CASCADE;"`
	Film     Film
	ImageURL string `gorm:"type:varchar(100)"`
}

// Сеансы и бронирование
type Session struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	FilmID    uint `gorm:"not null"`
	Film      Film
	HallID    uint `gorm:"not null"`
	Hall      CinemaHall
	StartTime time.Time `gorm:"not null"`
	Price     float64   `gorm:"type:numeric(12,2);not null"`
}

type Booking struct {
	ID uint `gorm:"primaryKey"`

	SessionID  uint `gorm:"not null;index:idx_seat,unique"`
	Session    Session
	CustomerID uint `gorm:"not null"`
	Customer   User

	RowNum  uint `gorm:"not null;index:idx_seat,unique"`
	SeatNum uint `gorm:"not null;index:idx_seat,unique"`

	SpendBonus    float64 `gorm:"type:numeric(12,2);default:0"`
	ReceivedBonus float64 `gorm:"type:numeric(12,2);default:0"`
	TotalPrice    float64 `gorm:"type:numeric(12,2);not null"`

	Status BookingStatus `gorm:"type:varchar(20);not null"`
}

// Финансы
type PaymentHistory struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	UserID    uint
	User      User
	Amount    float64          `gorm:"type:numeric(12,2)"`
	Desc      string           `gorm:"type:varchar(50)"`
	Operation PaymentOperation `gorm:"type:varchar(20);not null"`
}

type BonusHistory struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	UserID    uint
	User      User
	Amount    float64        `gorm:"type:numeric(12,2);not null"`
	Desc      string         `gorm:"type:varchar(50)"`
	Operation BonusOperation `gorm:"type:varchar(20);not null"`
}

type FilmGenre struct {
	FilmID  uint `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	GenreID uint `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
}

type Review struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	FilmID uint
	Film   Film
	UserID uint
	User   User
	Rating uint
	Coment string       `gorm:"type:varchar(100)"`
	Status ReviewStatus `gorm:"type:varchar(20);not null;default:'pending'"`
}
