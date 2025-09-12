package handlers

import (
	"net/http"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// CreateBookingHandler godoc
// @Summary Забронировать билет
// @Tags bookings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dt.CreateBookingDTI true "Данные для бронирования"
// @Success 201 {object} dt.CreateBookingDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /bookings [post]
func CreateBookingHandler(c *gin.Context) {
	var input dt.CreateBookingDTI

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	input.UserID = userID.(uint)
	booking, err := services.CreateBooking(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, dt.CreateBookingDTO{
		ID:     booking.ID,
		Status: booking.Status,
	})
}
