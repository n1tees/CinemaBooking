package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/models"
	"errors"
	"time"
)

// Получить все сеансы фильмов
func GetAllSessions() ([]models.Session, error) {
	var sessions []models.Session
	if err := db.DB.Preload("Film").Preload("Hall").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// Получить сеансы по фильму
func GetSessionsByFilm(filmID uint) ([]models.Session, error) {
	var sessions []models.Session
	if err := db.DB.Preload("Hall").Where("film_id = ?", filmID).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// ____________________________________________________ADMIN_ONLY____________________________________________________
// Создать сеанс
func CreateSession(filmID, hallID uint, start time.Time, price float64) (*models.Session, error) {
	session := models.Session{
		FilmID:    filmID,
		HallID:    hallID,
		StartTime: start,
		Price:     price,
	}
	if err := db.DB.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// Обновить сеанс
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
