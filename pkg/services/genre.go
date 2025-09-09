package services

import (
	"CinemaBooking/pkg/db"
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
func CreateGenre(genre *models.Genre) error {

	var count int64
	db.DB.Model(&models.Genre{}).Where("name = ?", genre.Name).Count(&count)
	if count > 0 {
		return errors.New("жанр с таким названием уже существует")
	}

	if err := db.DB.Create(genre).Error; err != nil {
		return err
	}
	return nil
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
