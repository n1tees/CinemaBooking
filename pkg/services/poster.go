package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/models"
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
func UpdatePoster(id uint, newData models.Poster) error {

	var poster models.Poster

	// Проверяем, есть ли такая афиша
	if err := db.DB.First(&poster, id).Error; err != nil {
		return err
	}

	if err := db.DB.Save(&newData).Error; err != nil {
		return err
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
