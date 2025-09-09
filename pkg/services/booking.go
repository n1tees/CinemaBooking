package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/models"
	"errors"

	"gorm.io/gorm"
)

// Забронировать билет
func CreateBooking(userID, sessionID, row, seat uint, spendBonus float64) (*models.Booking, error) {
	var booking models.Booking

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Проверка занятости места
		taken, err := isSeatTaken(tx, sessionID, row, seat)
		if err != nil {
			return err
		}
		if taken {
			return errors.New("место уже занято")
		}

		// 2. Загружаем пользователя и профиль
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return errors.New("пользователь не найден")
		}
		var profile models.Profile
		if err := tx.First(&profile, user.ProfileID).Error; err != nil {
			return errors.New("профиль не найден")
		}

		// 3. Проверяем бонусы
		if profile.Bonus < int(spendBonus) {
			return errors.New("недостаточно бонусов")
		}

		// 4. Загружаем сеанс
		var session models.Session
		if err := tx.First(&session, sessionID).Error; err != nil {
			return errors.New("сеанс не найден")
		}

		// 5. Считаем стоимость и бонусы
		totalPrice := session.Price - spendBonus
		if totalPrice < 0 {
			totalPrice = 0
		}
		receivedBonus := totalPrice * 0.1

		// 6. Обновляем баланс пользователя
		if err := tx.Model(&profile).
			Update("bonus", gorm.Expr("bonus - ? + ?", spendBonus, receivedBonus)).Error; err != nil {
			return err
		}
		if err := tx.Model(&profile).
			Update("balance", gorm.Expr("balance - ?", totalPrice)).Error; err != nil {
			return err
		}

		// 7. Создаём запись о бронировании
		booking = models.Booking{
			SessionID:     sessionID,
			CustomerID:    userID,
			RowNum:        row,
			SeatNum:       seat,
			SpendBonus:    spendBonus,
			ReceivedBonus: receivedBonus,
			TotalPrice:    totalPrice,
			Status:        models.BookingPaid,
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		// 8. Записываем историю оплат
		if totalPrice > 0 {
			payment := models.PaymentHistory{
				UserID:    userID,
				Amount:    -totalPrice,
				Operation: models.PaymentSpend,
			}
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
		}

		// 9. Записываем историю бонусов (раздельно списание и начисление)
		if spendBonus > 0 {
			bonusSpend := models.BonusHistory{
				UserID:    userID,
				Amount:    -int(spendBonus),
				Operation: models.BonusRedeem,
			}
			if err := tx.Create(&bonusSpend).Error; err != nil {
				return err
			}
		}
		if receivedBonus > 0 {
			bonusEarn := models.BonusHistory{
				UserID:    userID,
				Amount:    int(receivedBonus),
				Operation: models.BonusEarn,
			}
			if err := tx.Create(&bonusEarn).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// ____________________________________________________ADMIN_ONLY____________________________________________________
// Отменить бронирование
func CancelBooking(bookingID, userID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.First(&booking, bookingID).Error; err != nil {
			return errors.New("бронирование не найдено")
		}

		if booking.CustomerID != userID {
			return errors.New("нельзя отменить чужое бронирование")
		}
		if booking.Status != models.BookingPaid {
			return errors.New("бронирование нельзя отменить")
		}

		// Загружаем пользователя и профиль
		var user models.User
		if err := tx.First(&user, booking.CustomerID).Error; err != nil {
			return errors.New("пользователь не найден")
		}
		var profile models.Profile
		if err := tx.First(&profile, user.ProfileID).Error; err != nil {
			return errors.New("профиль не найден")
		}

		// Возвращаем деньги
		if booking.TotalPrice > 0 {
			if err := tx.Model(&profile).
				Update("balance", gorm.Expr("balance + ?", booking.TotalPrice)).Error; err != nil {
				return err
			}

			// Записываем возврат в историю оплат
			payment := models.PaymentHistory{
				UserID:    booking.CustomerID,
				Amount:    booking.TotalPrice,
				Operation: models.PaymentDeposit,
			}
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
		}

		// Снимаем бонусы, которые начислялись за покупку
		if booking.ReceivedBonus > 0 && profile.Bonus >= int(booking.ReceivedBonus) {
			if err := tx.Model(&profile).
				Update("bonus", gorm.Expr("bonus - ?", booking.ReceivedBonus)).Error; err != nil {
				return err
			}

			// Записываем в историю бонусов
			bonus := models.BonusHistory{
				UserID:    booking.CustomerID,
				Amount:    -int(booking.ReceivedBonus),
				Operation: models.BonusRedeem,
			}
			if err := tx.Create(&bonus).Error; err != nil {
				return err
			}
		}

		// Обновляем статус брони
		if err := tx.Model(&booking).
			Update("status", models.BookingCanceled).Error; err != nil {
			return err
		}

		return nil
	})
}

// ____________________________________________________INTERNAL____________________________________________________
func isSeatTaken(tx *gorm.DB, sessionID, row, seat uint) (bool, error) {
	var count int64

	err := tx.Model(&models.Booking{}).
		Where("session_id = ? AND row_num = ? AND seat_num = ? AND status IN ?",
			sessionID, row, seat, []models.BookingStatus{
				models.BookingReserved,
				models.BookingPaid,
			}).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
