package handlers

import (
	"net/http"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetMyPaymentsHandler godoc
// @Summary Получить историю пополнений
// @Tags wallet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.PaymentHistory
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet/payments [get]
func GetMyPaymentsHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	payments, err := services.GetMyPayments(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// RefillMyBalanceHandler godoc
// @Summary Пополнить баланс
// @Tags wallet
// @Security BearerAuth
// @Accept json
// @Produce json
//
//	@Param input body struct {
//	  Amount float64 `json:"amount" binding:"required"`
//	} true "Сумма пополнения"
//
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet/refill [post]
func RefillMyBalanceHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.RefillMyBalance(userID.(uint), input.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "баланс успешно пополнен"})
}

// GetBalanceHandler godoc
// @Summary Получить текущий баланс
// @Tags wallet
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]float64
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet/balance [get]
func GetBalanceHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	balance, err := services.GetBalance(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
