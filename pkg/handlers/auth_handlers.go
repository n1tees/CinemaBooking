package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// RegisterHandler godoc
// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dt.RegisterDTI true "Данные пользователя"
// @Success 201 {object} dt.RegisterDTO
// @Failure 400 {object} dt.ErrorResponse
// @Router /register [post]
func RegisterHandler(c *gin.Context) {
	var input dt.RegisterDTI

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Неверный формат запроса"})
		return
	}

	// Парсим дату рождения
	birth, err := makeDateByString(input.BirthDay)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: fmt.Sprintf("Неверный формат даты. Ожидается YYYY-MM-DD, получено %s", input.BirthDay),
		})

		return
	}

	// Передаём дату в структуру
	input.ParsedBirthDay = birth

	userID, err := services.RegUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dt.RegisterDTO{ID: userID})
}

// LoginHandler godoc
// @Summary Логин пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dt.LoginDTI true "Данные пользователя"
// @Success 200 {object} dt.LoginDTO
// @Failure 401 {object} dt.ErrorResponse
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var input dt.LoginDTI

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALVID_INPUT",
			Message: "Неверный формат запроса",
		})
		return
	}

	token, err := services.LoginUser(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.LoginDTO{JWT: token})
}

// функции для работы с датами и временем
func makeDateByString(date string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, errors.New("ошибка при парсинге даты")
	}
	return parsedDate, nil
}

// func makeTimeByString(timeStr string) (time.Time, error) {
// 	parsedTime, err := time.Parse("15:04", timeStr)
// 	if err != nil {
// 		return time.Time{}, errors.New("ошибка при парсинге времени")
// 	}
// 	return parsedTime, nil
// }
