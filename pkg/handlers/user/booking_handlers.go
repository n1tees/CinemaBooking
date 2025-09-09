package handlers

import (
	"net/http"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// CreateBookingHandler godoc
// @Summary Забронировать билет
// @Tags bookings
// @Security BearerAuth
// @Accept json
// @Produce json
//
//	@Param input body struct {
//	  SessionID  uint    `json:"session_id" binding:"required"`
//	  Row        uint    `json:"row" binding:"required"`
//	  Seat       uint    `json:"seat" binding:"required"`
//	  SpendBonus float64 `json:"spend_bonus"`
//	} true "Данные для бронирования"
//
// @Success 201 {object} models.Booking
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings [post]
func CreateBookingHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var input struct {
		SessionID  uint    `json:"session_id" binding:"required"`
		Row        uint    `json:"row" binding:"required"`
		Seat       uint    `json:"seat" binding:"required"`
		SpendBonus float64 `json:"spend_bonus"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := services.CreateBooking(
		userID.(uint),
		input.SessionID,
		input.Row,
		input.Seat,
		input.SpendBonus,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}
