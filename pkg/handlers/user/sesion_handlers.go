package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetAllSessionsHandler godoc
// @Summary Получить все предстоящие сеансы
// @Tags sessions
// @Produce json
// @Success 200 {array} models.Session
// @Failure 500 {object} map[string]string
// @Router /sessions [get]
func GetAllSessionsHandler(c *gin.Context) {
	sessions, err := services.GetAllSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// GetSessionsByFilmHandler godoc
// @Summary Получить предстоящие сеансы по фильму
// @Tags sessions
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {array} models.Session
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sessions/film/{id} [get]
func GetSessionsByFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	sessions, err := services.GetSessionsByFilm(uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// GetAvailableSeatsHandler godoc
// @Summary Получить доступные места на сеанс
// @Tags sessions
// @Produce json
// @Param id path int true "ID сеанса"
// @Success 200 {array} services.SeatDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sessions/{id}/seats [get]
func GetAvailableSeatsHandler(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	seats, err := services.GetAvailableSeats(uint(sessionID))
	if err != nil {
		if err.Error() == "сеанс не найден" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}
