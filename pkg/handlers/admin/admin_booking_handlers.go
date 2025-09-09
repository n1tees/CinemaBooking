package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// CancelBookingHandler godoc
// @Summary Отменить бронирование
// @Tags bookings
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID бронирования"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings/{id} [delete]
func CancelBookingHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	idStr := c.Param("id")
	bookingID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	err = services.CancelBooking(uint(bookingID), userID.(uint))
	if err != nil {
		switch err.Error() {
		case "бронирование не найдено":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "нельзя отменить чужое бронирование":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "бронирование нельзя отменить":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "бронирование отменено"})
}
