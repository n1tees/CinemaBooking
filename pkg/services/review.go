package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/models"
	"errors"
)

// Получить отзывы по фильму
func GetReviewsByFilm(filmID uint) ([]models.Review, error) {
	var reviews []models.Review

	if err := db.DB.Preload("User.Profile").
		Where("film_id = ? AND status = ?", filmID, models.ReviewApproved).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}

// Добавить отзыв (уходит на модерацию)
func AddReview(userID, filmID, rating uint, comment string) error {
	review := models.Review{
		FilmID: filmID,
		UserID: userID,
		Rating: rating,
		Coment: comment,
		Status: models.ReviewPending,
	}

	var count int64
	db.DB.Model(&models.Review{}).
		Where("film_id = ? AND user_id = ?", filmID, userID).
		Count(&count)
	if count > 0 {
		return errors.New("вы уже оставили отзыв на этот фильм")
	}

	if err := db.DB.Create(&review).Error; err != nil {
		return err
	}

	return nil
}

// Получить рейтинг фильма
func GetFilmRating(filmID uint) (*float64, error) {
	var avgRating float64

	if err := db.DB.Model(&models.Review{}).
		Where("film_id = ? AND status = ?", filmID, models.ReviewApproved).
		Select("AVG(rating)").
		Scan(&avgRating).Error; err != nil {
		return nil, err
	}

	if avgRating == 0 {
		return nil, nil // нет отзывов
	}

	return &avgRating, nil
}

// ____________________________________________________ADMIN_ONLY____________________________________________________
// Апрувнуть отзыв
func ApproveReview(reviewID uint) error {
	return db.DB.Model(&models.Review{}).
		Where("id = ?", reviewID).
		Update("status", models.ReviewApproved).Error
}

// Отклонить отзыв
func RejectReview(reviewID uint) error {
	return db.DB.Model(&models.Review{}).
		Where("id = ?", reviewID).
		Update("status", models.ReviewRejected).Error
}

// Удалить отзыв пользователем (свой)
func DeleteOwnReview(userID, reviewID uint) error {
	var review models.Review

	// Проверяем, что отзыв существует и принадлежит этому пользователю
	if err := db.DB.First(&review, "id = ? AND user_id = ?", reviewID, userID).Error; err != nil {
		return err
	}

	// Удаляем
	if err := db.DB.Delete(&review).Error; err != nil {
		return err
	}

	return nil
}
