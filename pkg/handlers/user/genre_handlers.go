package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetAllGenresHandler godoc
// @Summary Получить все жанры
// @Tags genres
// @Produce json
// @Success 200 {array} dt.GenreDTO
// @Failure 500 {object} dt.ErrorResponse
// @Router /genres [get]
func GetAllGenresHandler(c *gin.Context) {
	genres, err := services.GetAllGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	var result []dt.GenreDTO
	for _, g := range genres {
		result = append(result, dt.GenreDTO{
			Name: g.Name,
		})
	}

	c.JSON(http.StatusOK, result)
}

// GetGenresByFilmHandler godoc
// @Summary Получить жанры фильма
// @Tags genres
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {array} dt.GenreDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /films/{id}/genres [get]
func GetGenresByFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid film ID",
		})
		return
	}

	genres, err := services.GetGenresByFilm(uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	var result []dt.GenreDTO
	for _, g := range genres {
		result = append(result, dt.GenreDTO{
			Name: g.Name,
		})
	}

	c.JSON(http.StatusOK, result)
}
