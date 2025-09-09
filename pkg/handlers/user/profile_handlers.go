package handlers

import (
	"net/http"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetUserInfoHandler godoc
// @Summary Получить данные профиля
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.Profile
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile [get]
func GetUserInfoHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	profile, err := services.GetUserInfo(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ChangePasswordHandler godoc
// @Summary Изменить пароль
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body services.ChangePasswordInput true "Смена пароля"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile/password [patch]
func ChangePasswordHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var input services.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ChangePassword(userID.(uint), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "пароль изменён"})
}

// UpdateProfileHandler godoc
// @Summary Обновить данные профиля (частично)
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Поля для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile [patch]
func UpdateProfileHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateProfile(userID.(uint), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "профиль обновлён"})
}
