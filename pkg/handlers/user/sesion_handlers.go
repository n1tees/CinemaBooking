package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetAllSessionsHandler godoc
// @Summary Получить все предстоящие сеансы
// @Tags sessions
// @Produce json
// @Success 200 {array} dt.SessionDTO
// @Failure 500 {object} dt.ErrorResponse
// @Router /sessions [get]
func GetAllSessionsHandler(c *gin.Context) {
	sessions, err := services.GetAllSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	var result []dt.SessionDTO
	for _, s := range sessions {
		result = append(result, dt.SessionDTO{
			ID:        s.ID,
			FilmID:    s.FilmID,
			HallID:    s.HallID,
			StartTime: s.StartTime,
			Price:     s.Price,
		})
	}

	c.JSON(http.StatusOK, result)
}

// GetSessionsByFilmHandler godoc
// @Summary Получить предстоящие сеансы по фильму
// @Tags sessions
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {array} dt.SessionDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /sessions/film/{id} [get]
func GetSessionsByFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid film ID",
		})
		return
	}

	sessions, err := services.GetSessionsByFilm(uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	var result []dt.SessionDTO
	for _, s := range sessions {
		result = append(result, dt.SessionDTO{
			ID:        s.ID,
			FilmID:    s.FilmID,
			HallID:    s.HallID,
			StartTime: s.StartTime,
			Price:     s.Price,
		})
	}

	c.JSON(http.StatusOK, result)
}

// GetSeatsBySessionHandler godoc
// @Summary Получить все места на сеанс с их статусом (free/taken)
// @Tags sessions
// @Produce json
// @Param id path int true "ID сеанса"
// @Success 200 {array} dt.SeatDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /sessions/{id}/seats [get]
func GetAvailableSeatsHandler(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid session ID",
		})
		return
	}

	seats, err := services.GetAvailableSeats(uint(sessionID))
	if err != nil {
		if err.Error() == "сеанс не найден" {
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, seats)
}
