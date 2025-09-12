package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// CancelBookingHandler godoc
// @Summary Отменить бронирование
// @Tags bookings
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID бронирования"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /bookings/{id} [delete]
func CancelBookingHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user_id not found",
		})
		return
	}

	idStr := c.Param("id")
	bookingID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid booking ID",
		})
		return
	}

	err = services.CancelBooking(uint(bookingID), userID.(uint))
	if err != nil {
		switch err.Error() {
		case "бронирование не найдено":
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
		case "нельзя отменить чужое бронирование":
			c.JSON(http.StatusForbidden, dt.ErrorResponse{
				Code:    "FORBIDDEN",
				Message: err.Error(),
			})
		case "бронирование нельзя отменить":
			c.JSON(http.StatusBadRequest, dt.ErrorResponse{
				Code:    "INVALID_STATE",
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, dt.ServAnswerDTO{
		Answer: "бронирование отменено",
	})
}
