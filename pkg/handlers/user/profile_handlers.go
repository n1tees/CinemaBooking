package handlers

import (
	"net/http"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetUserInfoHandler godoc
// @Summary Получить данные профиля
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dt.ProfileDTO
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /profile [get]
func GetUserInfoHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user_id not found",
		})
		return
	}

	profile, err := services.GetUserInfo(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.ProfileDTO{
		FirstName:  profile.FirstName,
		SecondName: profile.SecondName,
		Email:      profile.Email,
		Balance:    profile.Balance,
		Bonus:      profile.Bonus,
	})
}

// ChangePasswordHandler godoc
// @Summary Изменить пароль
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dt.ChangePasswordDTI true "Смена пароля"
// @Success 200 {object} dt.ChangePasswordDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /profile/password [patch]
func ChangePasswordHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user_id not found",
		})
		return
	}

	var input dt.ChangePasswordDTI
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	if err := services.ChangePassword(userID.(uint), input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.ChangePasswordDTO{
		UserID: userID.(uint),
		Status: "Пароль изменен",
	})
}

// UpdateProfileHandler godoc
// @Summary Обновить данные профиля (частично)
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body object true "Поля для обновления"
// @Success 200 {object} dt.UpdateProfileDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /profile [patch]
func UpdateProfileHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user not found in context",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	if err := services.UpdateProfile(userID.(uint), updates); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.UpdateProfileDTO{
		Status: "профиль обновлён",
	})
}
