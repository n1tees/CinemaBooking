package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/models"
	"errors"
)

// Посмотреть афиши
func GetAllPosters() ([]models.Poster, error) {
	var posters []models.Poster

	if err := db.DB.Find(&posters).Error; err != nil {
		return nil, err
	}

	return posters, nil
}

// ____________________________________________________ADMIN_ONLY____________________________________________________
// Создать афишу
func CreatePoster(poster *models.Poster) error {

	if err := db.DB.Create(poster).Error; err != nil {
		return err
	}

	return nil

}

// Изменить афишу
func UpdatePoster(id uint, updates map[string]interface{}) error {
	// Фильтруем значения
	filtered := FilterUpdates(updates)
	if len(filtered) == 0 {
		return errors.New("пустой запрос на обновление")
	}

	// Обновляем только нужные поля
	if err := db.DB.Model(&models.Poster{}).
		Where("id = ?", id).
		Updates(filtered).Error; err != nil {
		return errors.New("ошибка при обновлении афиши")
	}

	return nil
}

// Удалить афишу
func DeletePoster(id uint) error {
	if err := db.DB.Delete(&models.Poster{}, id).Error; err != nil {
		return err
	}
	return nil
}
