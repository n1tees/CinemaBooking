package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/models"
	"errors"
	"time"

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

		// 2. Загружаем пользователя и проверяем бонусы
		var profile models.Profile
		if err := tx.First(&profile, userID).Error; err != nil {
			return err
		}
		if profile.Bonus < int(spendBonus) {
			return errors.New("недостаточно бонусов")
		}

		// 3. Загружаем сеанс
		var session models.Session
		if err := tx.First(&session, sessionID).Error; err != nil {
			return err
		}

		// 4. Считаем стоимость и бонусы
		totalPrice := session.Price - spendBonus
		if totalPrice < 0 {
			totalPrice = 0
		}
		receivedBonus := totalPrice * 0.1

		// 5. Обновляем баланс пользователя
		if err := tx.Model(&profile).
			Update("bonus", gorm.Expr("bonus - ? + ?", spendBonus, receivedBonus)).Error; err != nil {
			return err
		}
		if err := tx.Model(&profile).
			Update("balance", gorm.Expr("balance - ?", totalPrice)).Error; err != nil {
			return err
		}

		// 6. Создаём запись о бронировании
		booking = models.Booking{
			SessionID:     sessionID,
			CustomerID:    userID,
			RowNum:        row,
			SeatNum:       seat,
			SpendBonus:    spendBonus,
			ReceivedBonus: receivedBonus,
			TotalPrice:    totalPrice,
			Status:        models.BookingPaid,
			CreatedAt:     time.Now(),
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		// 7. Записываем историю
		payment := models.PaymentHistory{
			UserID:    userID,
			Amount:    -totalPrice,
			Operation: models.PaymentSpend,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		bonus := models.BonusHistory{
			UserID:    userID,
			Amount:    -int(spendBonus) + int(receivedBonus),
			Operation: models.BonusEarn,
		}
		if err := tx.Create(&bonus).Error; err != nil {
			return err
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
			return err
		}

		if booking.CustomerID != userID {
			return errors.New("нельзя отменить чужое бронирование")
		}
		if booking.Status != models.BookingPaid {
			return errors.New("бронирование нельзя отменить")
		}

		// Загружаем профиль
		var profile models.Profile
		if err := tx.First(&profile, booking.CustomerID).Error; err != nil {
			return err
		}

		// Возвращаем деньги
		if err := tx.Model(&profile).
			Update("balance", gorm.Expr("balance + ?", booking.TotalPrice)).Error; err != nil {
			return err
		}

		// Снимаем бонусы (если не хватает — просто не возвращаем)
		if profile.Bonus >= int(booking.ReceivedBonus) {
			if err := tx.Model(&profile).
				Update("bonus", gorm.Expr("bonus - ?", booking.ReceivedBonus)).Error; err != nil {
				return err
			}
		}

		// Обновляем статус брони
		if err := tx.Model(&booking).
			Update("status", models.BookingCanceled).Error; err != nil {
			return err
		}

		// История: возврат денег
		payment := models.PaymentHistory{
			UserID:    booking.CustomerID,
			Amount:    booking.TotalPrice,
			Operation: models.PaymentDeposit,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		// История: списание бонусов
		bonus := models.BonusHistory{
			UserID:    booking.CustomerID,
			Amount:    -int(booking.ReceivedBonus),
			Operation: models.BonusRedeem,
		}
		if err := tx.Create(&bonus).Error; err != nil {
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
