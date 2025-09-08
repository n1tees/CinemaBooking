package services

import (
	"fmt"
	"time"

	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/models"

	"gorm.io/gorm"
)

// Получить текущий баланс бонусов
func GetBonusBalance(userID uint) (int, error) {
	var profile models.Profile
	if err := db.DB.First(&profile, "id = ?", userID).Error; err != nil {
		return 0, err
	}
	return profile.Bonus, nil
}

// Получить историю бонусов
func GetBonusHistory(userID uint) ([]models.BonusHistory, error) {
	var history []models.BonusHistory

	if err := db.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

// ____________________________________________________INTERNAL____________________________________________________
// Добавить бонусы
func AddBonus(userID uint, amount int, description string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Обновляем баланс
		if err := tx.Model(&models.Profile{}).
			Where("id = ?", userID).
			Update("bonus", gorm.Expr("bonus + ?", amount)).Error; err != nil {
			return err
		}

		// Пишем в историю
		history := models.BonusHistory{
			UserID:    userID,
			Amount:    amount,
			Operation: models.BonusEarn,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})
}

// Списать бонусы
func SpendBonus(userID uint, amount int, description string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Проверяем баланс
		var profile models.Profile
		if err := tx.First(&profile, "id = ?", userID).Error; err != nil {
			return err
		}
		if profile.Bonus < amount {
			return fmt.Errorf("недостаточно бонусов")
		}

		// Списываем
		if err := tx.Model(&models.Profile{}).
			Where("id = ?", userID).
			Update("bonus", gorm.Expr("bonus - ?", amount)).Error; err != nil {
			return err
		}

		// Записываем в историю
		history := models.BonusHistory{
			UserID:    userID,
			Amount:    -amount,
			Operation: models.BonusRedeem,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})
}
