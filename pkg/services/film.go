package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/models"
	"errors"
	"fmt"
	"time"
)

// Получить конкретный фильм
func GetFilm(id uint) (*models.Film, error) {
	var film models.Film
	if err := db.DB.Preload("Genres").Preload("Posters").First(&film, id).Error; err != nil {
		return nil, err
	}
	return &film, nil
}

// Получить все фильмы с фильтром по жанрам
func GetAllFilms(genres []uint) ([]models.Film, error) {
	var films []models.Film
	query := db.DB.Preload("Posters").Preload("Genres")

	if len(genres) > 0 {
		// Выбираем только фильмы, у которых есть ВСЕ жанры из списка
		subQuery := db.DB.Table("film_genres").
			Select("film_id").
			Where("genre_id IN ?", genres).
			Group("film_id").
			Having("COUNT(DISTINCT genre_id) = ?", len(genres))

		query = query.Where("id IN (?)", subQuery)
	}

	if err := query.Find(&films).Error; err != nil {
		return nil, err
	}

	return films, nil
}

// ____________________________________________________ADMIN_ONLY____________________________________________________
// Создать фильм
func CreateFilm(input dt.CreateFilmDTI) (*dt.CreateFilmDTO, error) {
	// парсим дату
	releaseDate, err := time.Parse("2006-01-02", input.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid release_date format: %w", err)
	}

	// собираем модель фильма
	film := models.Film{
		Title:       input.Title,
		Desc:        input.Description,
		Duration:    input.Duration,
		AgeRating:   input.AgeRating,
		ReleaseDate: releaseDate,
	}

	// подтягиваем жанры, если есть
	if len(input.Genres) > 0 {
		var genres []models.Genre
		if err := db.DB.Find(&genres, input.Genres).Error; err != nil {
			return nil, fmt.Errorf("ошибка загрузки жанров: %w", err)
		}
		film.Genres = genres
	}

	// сохраняем фильм
	if err := db.DB.Create(&film).Error; err != nil {
		return nil, err
	}

	// возвращаем только ID
	return &dt.CreateFilmDTO{ID: film.ID}, nil
}

// Обновить фильм
func UpdateFilm(id uint, updates map[string]interface{}) error {
	filtered := FilterUpdates(updates)
	if len(filtered) == 0 {
		return errors.New("пустой запрос")
	}

	if err := db.DB.Model(&models.Film{}).
		Where("id = ?", id).
		Updates(filtered).Error; err != nil {
		return errors.New("ошибка при обновлении фильма")
	}
	return nil
}

// Удалить фильм
func DeleteFilm(id uint) error {
	if err := db.DB.Delete(&models.Film{}, id).Error; err != nil {
		return err
	}
	return nil
}
