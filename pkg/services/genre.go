package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/models"
	"errors"
)

// Получить все жанры
func GetAllGenres() ([]models.Genre, error) {
	var genres []models.Genre
	if err := db.DB.Order("name ASC").Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}

// ____________________________________________________ADMIN_ONLY____________________________________________________
// Создать жанр
func CreateGenre(input dt.CreateGenreDTI) (uint, error) {
	var count int64
	db.DB.Model(&models.Genre{}).Where("name = ?", input.Name).Count(&count)
	if count > 0 {
		return 0, errors.New("жанр с таким названием уже существует")
	}

	genre := models.Genre{
		Name: input.Name,
	}

	if err := db.DB.Create(&genre).Error; err != nil {
		return 0, err
	}
	return genre.ID, nil
}

// Обновить жанр
func UpdateGenre(id uint, updates map[string]interface{}) error {
	filtered := FilterUpdates(updates)
	if len(filtered) == 0 {
		return errors.New("пустой запрос")
	}

	if err := db.DB.Model(&models.Genre{}).
		Where("id = ?", id).
		Updates(filtered).Error; err != nil {
		return errors.New("ошибка при обновлении жанра")
	}

	return nil
}

// Удалить жанр
func DeleteGenre(id uint) error {
	if err := db.DB.Delete(&models.Genre{}, id).Error; err != nil {
		return err
	}
	return nil
}
