package handlers

import (
	"net/http"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetMyPaymentsHandler godoc
// @Summary Получить историю пополнений
// @Tags wallet
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dt.PaymentHistoryDTO
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /wallet/payments [get]
func GetMyPaymentsHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user not found in context",
		})
		return
	}

	payments, err := services.GetMyPayments(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	var result []dt.PaymentHistoryDTO
	for _, p := range payments {
		result = append(result, dt.PaymentHistoryDTO{
			Amount:    p.Amount,
			Operation: string(p.Operation),
			CreatedAt: p.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}

// RefillMyBalanceHandler godoc
// @Summary Пополнить баланс
// @Tags wallet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dt.RefillBalanceDTI true "Сумма пополнения"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /wallet/refill [post]
func RefillMyBalanceHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user not found in context",
		})
		return
	}

	var input dt.RefillBalanceDTI
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
	}

	if err := services.RefillMyBalance(userID.(uint), input.Amount); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.ServAnswerDTO{
		Answer: "баланс успешно пополнен",
	})
}

// GetBalanceHandler godoc
// @Summary Получить текущий баланс
// @Tags wallet
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dt.PaymentDTO
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /wallet/balance [get]
func GetBalanceHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user not found in context",
		})
		return
	}

	balance, err := services.GetBalance(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.PaymentDTO{Balance: balance})
}
