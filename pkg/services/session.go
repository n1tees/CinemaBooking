package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type HallStructure struct {
	Rows []struct {
		Row   uint `json:"row"`
		Seats uint `json:"seats"`
	} `json:"rows"`
}

type SeatDTO struct {
	Row   uint   `json:"row"`
	Seat  uint   `json:"seat"`
	State string `json:"state"` // "free" или "taken"
}

// Получить все предстоящие сеансы (от сегодня и на 2 месяца вперёд)
func GetAllSessions() ([]models.Session, error) {
	var sessions []models.Session

	start := time.Now().Truncate(24 * time.Hour)
	end := start.AddDate(0, 2, 0) // +2 месяца

	if err := db.DB.Preload("Film").Preload("Hall").
		Where("start_time >= ? AND start_time < ?", start, end).
		Order("start_time ASC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

// Получить предстоящие сеансы по фильму
func GetSessionsByFilm(filmID uint) ([]models.Session, error) {
	var sessions []models.Session

	start := time.Now().Truncate(24 * time.Hour)
	end := start.AddDate(0, 2, 0) // +2 месяца

	if err := db.DB.Preload("Hall").
		Where("film_id = ? AND start_time >= ? AND start_time < ?", filmID, start, end).
		Order("start_time ASC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

// Получить свободные места
func GetAvailableSeats(sessionID uint) ([]SeatDTO, error) {
	// 1. Находим сеанс с залом
	var session models.Session
	if err := db.DB.Preload("Hall").First(&session, sessionID).Error; err != nil {
		return nil, errors.New("сеанс не найден")
	}

	// 2. Парсим JSON структуру зала
	var structure HallStructure
	if err := json.Unmarshal(session.Hall.Structure, &structure); err != nil {
		return nil, errors.New("ошибка парсинга структуры зала")
	}

	// 3. Получаем занятые места
	var taken []struct {
		Row  uint
		Seat uint
	}
	if err := db.DB.Model(&models.Booking{}).
		Select("row_num as row, seat_num as seat").
		Where("session_id = ? AND status IN ?", sessionID, []models.BookingStatus{
			models.BookingReserved,
			models.BookingPaid,
		}).
		Find(&taken).Error; err != nil {
		return nil, err
	}

	takenMap := make(map[string]bool)
	for _, s := range taken {
		key := fmt.Sprintf("%d-%d", s.Row, s.Seat)
		takenMap[key] = true
	}

	// 4. Формируем все места
	var seats []SeatDTO
	for _, row := range structure.Rows {
		for seat := uint(1); seat <= row.Seats; seat++ {
			key := fmt.Sprintf("%d-%d", row.Row, seat)
			state := "free"
			if takenMap[key] {
				state = "taken"
			}
			seats = append(seats, SeatDTO{
				Row:   row.Row,
				Seat:  seat,
				State: state,
			})
		}
	}

	return seats, nil
}

// ____________________________________________________ADMIN_ONLY____________________________________________________
// Создать сеанс
// CreateSession создает новый сеанс
func CreateSession(input dt.CreateSessionDTI) (*dt.CreateSessionDTO, error) {
	// Проверка даты и цены
	if input.Start.Before(time.Now()) {
		return nil, errors.New("нельзя создать сеанс в прошлом")
	}
	if input.Price <= 0 {
		return nil, errors.New("цена должна быть больше нуля")
	}

	session := models.Session{
		FilmID:    input.FilmID,
		HallID:    input.HallID,
		StartTime: input.Start,
		Price:     input.Price,
	}

	if err := db.DB.Create(&session).Error; err != nil {
		return nil, err
	}

	return &dt.CreateSessionDTO{ID: session.ID}, nil
}

// Обновить сеанс (частично)
func UpdateSession(id uint, updates map[string]interface{}) error {
	filtered := FilterUpdates(updates)
	if len(filtered) == 0 {
		return errors.New("пустой запрос")
	}

	if err := db.DB.Model(&models.Session{}).
		Where("id = ?", id).
		Updates(filtered).Error; err != nil {
		return errors.New("ошибка при обновлении сеанса")
	}
	return nil
}

// Удалить сеанс
func DeleteSession(id uint) error {
	if err := db.DB.Delete(&models.Session{}, id).Error; err != nil {
		return err
	}
	return nil
}
