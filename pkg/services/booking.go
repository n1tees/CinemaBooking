package services

import (
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/models"
	"errors"

	"gorm.io/gorm"
)

// Забронировать билет
func CreateBooking(input dt.CreateBookingDTI) (*models.Booking, error) {
	var booking models.Booking

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		// 1. Загружаем пользователя и профиль
		var user models.User
		if err := tx.First(&user, input.UserID).Error; err != nil {
			return errors.New("пользователь не найден")
		}
		var profile models.Profile
		if err := tx.First(&profile, user.ProfileID).Error; err != nil {
			return errors.New("профиль не найден")
		}

		// 2. Проверка занятости места
		taken, err := isSeatTaken(tx, input.SessionID, input.RowNum, input.SeatNum)
		if err != nil {
			return err
		}
		if taken {
			return errors.New("место занято")
		}

		// 3. Загружаем сеанс
		var session models.Session
		if err := tx.First(&session, input.SessionID).Error; err != nil {
			return errors.New("сеанс не найден")
		}
		var ReceivedBonus, SpendBonus, TotalPrice float64

		if input.UseBonus {

			ReceivedBonus = 0.0

			if profile.Bonus > session.Price {
				SpendBonus = session.Price
				TotalPrice = 0.0
			} else {
				TotalPrice = session.Price - profile.Bonus
				SpendBonus = profile.Bonus
			}
		} else {

			ReceivedBonus = session.Price * 0.1
			TotalPrice = session.Price
			SpendBonus = 0.0
		}

		// 6. Обновляем баланс пользователя
		if err := tx.Model(&profile).
			Update("bonus", gorm.Expr("bonus - ? + ?", SpendBonus, ReceivedBonus)).Error; err != nil {
			return err
		}
		if err := tx.Model(&profile).
			Update("balance", gorm.Expr("balance - ?", TotalPrice)).Error; err != nil {
			return err
		}

		// 7. Создаём запись о бронировании
		booking = models.Booking{
			SessionID:     input.SessionID,
			CustomerID:    input.UserID,
			RowNum:        input.RowNum,
			SeatNum:       input.SeatNum,
			SpendBonus:    SpendBonus,
			ReceivedBonus: ReceivedBonus,
			TotalPrice:    TotalPrice,
			Status:        models.BookingPaid,
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		// 8. Записываем историю оплат
		if TotalPrice > 0 {
			payment := models.PaymentHistory{
				UserID:    input.UserID,
				Amount:    TotalPrice,
				Operation: models.PaymentSpend,
			}
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
		}

		// 9. Записываем историю бонусов (раздельно списание и начисление)
		if SpendBonus > 0 {
			bonusSpend := models.BonusHistory{
				UserID:    input.UserID,
				Amount:    SpendBonus,
				Operation: models.BonusRedeem,
			}
			if err := tx.Create(&bonusSpend).Error; err != nil {
				return err
			}
		}
		if ReceivedBonus > 0 {
			bonusEarn := models.BonusHistory{
				UserID:    input.UserID,
				Amount:    ReceivedBonus,
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
		if booking.ReceivedBonus > 0 {
			if profile.Bonus >= booking.ReceivedBonus {
				// хватает бонусов → списываем бонусами
				if err := tx.Model(&profile).
					Update("bonus", gorm.Expr("bonus - ?", booking.ReceivedBonus)).Error; err != nil {
					return err
				}

				bonus := models.BonusHistory{
					UserID:    booking.CustomerID,
					Amount:    booking.ReceivedBonus,
					Operation: models.BonusRedeem,
				}
				if err := tx.Create(&bonus).Error; err != nil {
					return err
				}
			} else {
				// не хватает бонусов → остаток списываем с баланса
				missing := booking.ReceivedBonus - profile.Bonus

				if err := tx.Model(&profile).Updates(map[string]interface{}{
					"bonus":   0, // все бонусы обнуляем
					"balance": gorm.Expr("balance - ?", missing),
				}).Error; err != nil {
					return err
				}

				// История по бонусам
				if profile.Bonus > 0 {
					bonus := models.BonusHistory{
						UserID:    booking.CustomerID,
						Amount:    profile.Bonus,
						Operation: models.BonusRedeem,
					}
					if err := tx.Create(&bonus).Error; err != nil {
						return err
					}
				}

				// История по балансу (добор недостающих бонусов)
				payment := models.PaymentHistory{
					UserID:    booking.CustomerID,
					Amount:    missing,
					Operation: models.PaymentSpend, // или отдельный тип, например PaymentAdjustment
				}
				if err := tx.Create(&payment).Error; err != nil {
					return err
				}
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
