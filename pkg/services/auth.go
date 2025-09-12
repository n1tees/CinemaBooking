package services

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"

	"CinemaBooking/config"
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/dt"

	"CinemaBooking/pkg/models"
)

// Создать пользователя
func RegUser(input dt.RegisterDTI) (uint, error) {

	// Проверка, что логин свободен
	if _, err := searchAuthByLogin(input.Login); err == nil {
		return 0, errors.New("логин занят")
	}

	// Проверка, что телефон свободен
	if _, err := searchProfileByPhone(input.Phone); err == nil {
		return 0, errors.New("номер телефона привязан к другому аккаунту")
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("ошибка при хешировании пароля")
	}

	var user models.User
	// Начинаем транзакцию
	err = db.DB.Transaction(func(tx *gorm.DB) error {

		// 1. Создаём AuthCredential
		auth := models.AuthCredential{
			Login:        input.Login,
			PasswordHash: hash,
		}
		if err := tx.Create(&auth).Error; err != nil {
			return errors.New("ошибка создания логина")
		}

		// 2. Создаём Profile
		profile := models.Profile{
			FirstName:  input.FirstName,
			SecondName: input.SecondName,
			Phone:      input.Phone,
			BirthDay:   input.ParsedBirthDay,
		}
		if err := tx.Create(&profile).Error; err != nil {
			return errors.New("ошибка создания профиля")
		}

		// 3. Создаём User
		user = models.User{
			UserType:  models.Customer,
			AuthID:    auth.ID,
			ProfileID: profile.ID,
		}
		if err := tx.Create(&user).Error; err != nil {
			return errors.New("ошибка создания пользователя")
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	// Если всё ок — получить ID пользователя
	return user.ID, nil
}

// Авторизация пользователя
func LoginUser(input dt.LoginDTI) (string, error) {

	var auth *models.AuthCredential
	var err error

	if auth, err = verifyCredentials(input.Login, input.Password); err != nil {
		return "", err
	}

	var user *models.User
	if user, err = searchUserByAuth(auth); err != nil {
		return "", err
	}

	return generateJWT(user.ID, user.UserType)
}

// ____________________________________________________INTERNAL____________________________________________________
// Генерация токена JWT
func generateJWT(userID uint, userType models.UserType) (string, error) {

	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 2).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		return "", errors.New("ошибка при генерации токена")
	}

	return tokenString, nil
}

// Проверка Логин-Пароль
func verifyCredentials(login string, password string) (*models.AuthCredential, error) {

	// Проверка существования логина
	var auth *models.AuthCredential
	auth, err := searchAuthByLogin(login)
	if err != nil {
		return nil, errors.New("неверный логин или пароль")
	}

	// Cравниваем введённый пароль с захэшированным паролем из БД
	if err := bcrypt.CompareHashAndPassword(auth.PasswordHash, []byte(password)); err != nil {
		return nil, errors.New("неверный логин или пароль")
	}
	return auth, nil

}
