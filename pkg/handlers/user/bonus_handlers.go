package handlers

import (
	"net/http"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetBonusBalanceHandler godoc
// @Summary Получить текущий баланс бонусов
// @Tags bonus
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dt.BonusBalanceDTO
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /bonus/balance [get]
func GetBonusBalanceHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user_id not found",
		})
		return
	}

	balance, err := services.GetBonusBalance(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.BonusBalanceDTO{
		Balance: balance,
	})
}

// GetBonusHistoryHandler godoc
// @Summary Получить историю бонусов
// @Tags bonus
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dt.BonusHistoryDTO
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /bonus/history [get]
func GetBonusHistoryHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user_id not found",
		})
		return
	}

	history, err := services.GetBonusHistory(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			dt.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		return
	}

	c.JSON(http.StatusOK, dt.BonusHistoryDTO{
		HistoryBalance: history,
	})
}
