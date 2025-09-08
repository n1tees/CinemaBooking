package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/models"
)

// Получить все жанры фильма
func GetGenresByFilm(filmID uint) ([]models.Genre, error) {
	var genres []models.Genre

	if err := db.DB.
		Joins("JOIN film_genres fg ON fg.genre_id = genres.id").
		Where("fg.film_id = ?", filmID).
		Find(&genres).Error; err != nil {
		return nil, err
	}

	return genres, nil
}

// ____________________________________________________ADMIN_ONLY____________________________________________________
// Привязать жанр к фильму
func AssignGenreToFilm(filmID uint, genreID uint) error {
	link := models.FilmGenre{
		FilmID:  filmID,
		GenreID: genreID,
	}

	if err := db.DB.Create(&link).Error; err != nil {
		return err
	}
	return nil
}

// Убрать жанр у фильма
func RemoveGenreFromFilm(filmID uint, genreID uint) error {
	if err := db.DB.Where("film_id = ? AND genre_id = ?", filmID, genreID).
		Delete(&models.FilmGenre{}).Error; err != nil {
		return err
	}
	return nil
}
