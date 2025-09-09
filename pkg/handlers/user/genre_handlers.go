package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetAllGenresHandler godoc
// @Summary Получить все жанры
// @Tags genres
// @Produce json
// @Success 200 {array} models.Genre
// @Failure 500 {object} map[string]string
// @Router /genres [get]
func GetAllGenresHandler(c *gin.Context) {
	genres, err := services.GetAllGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, genres)
}

// GetGenresByFilmHandler godoc
// @Summary Получить жанры фильма
// @Tags genres
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {array} models.Genre
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /films/{id}/genres [get]
func GetGenresByFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	genres, err := services.GetGenresByFilm(uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genres)
}
