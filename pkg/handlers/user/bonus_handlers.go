package handlers

import (
	"net/http"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetBonusBalanceHandler godoc
// @Summary Получить текущий баланс бонусов
// @Tags bonus
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bonus/balance [get]
func GetBonusBalanceHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	balance, err := services.GetBonusBalance(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// GetBonusHistoryHandler godoc
// @Summary Получить историю бонусов
// @Tags bonus
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.BonusHistory
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bonus/history [get]
func GetBonusHistoryHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	history, err := services.GetBonusHistory(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
