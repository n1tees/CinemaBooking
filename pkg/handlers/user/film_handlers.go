package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetFilmHandler godoc
// @Summary Получить фильм по ID
// @Tags films
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {object} models.Film
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /films/{id} [get]
func GetFilmHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	film, err := services.GetFilm(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "film not found"})
		return
	}

	c.JSON(http.StatusOK, film)
}

// GetAllFilmsHandler godoc
// @Summary Получить список фильмов
// @Tags films
// @Produce json
// @Param genres query string false "Список ID жанров через запятую"
// @Success 200 {array} models.Film
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /films [get]
func GetAllFilmsHandler(c *gin.Context) {
	var genreIDs []uint

	genresQuery := c.Query("genres")
	if genresQuery != "" {
		parts := strings.Split(genresQuery, ",")
		for _, p := range parts {
			val, err := strconv.ParseUint(strings.TrimSpace(p), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre ID in query"})
				return
			}
			genreIDs = append(genreIDs, uint(val))
		}
	}

	films, err := services.GetAllFilms(genreIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, films)
}
