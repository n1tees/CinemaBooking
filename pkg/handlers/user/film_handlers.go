package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetFilmHandler godoc
// @Summary Получить фильм по ID
// @Tags films
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {object} dt.FilmDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /films/{id} [get]
func GetFilmHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Invalid film_id",
		})
		return
	}

	film, err := services.GetFilm(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "film not found",
		})
		return
	}
	genresModel, err := services.GetGenresByFilm(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	var genresDTO []dt.GenreDTO
	for _, g := range genresModel {
		genresDTO = append(genresDTO, dt.GenreDTO{
			Name: g.Name,
		})
	}

	c.JSON(http.StatusOK, dt.FilmDTO{
		Title:       film.Title,
		Description: film.Desc,
		AgeRating:   film.AgeRating,
		Duration:    film.Duration,
		ReleaseDate: film.ReleaseDate.Format("2006-01-02"),

		Genres: genresDTO,
	})
}

// GetAllFilmsHandler godoc
// @Summary Получить список фильмов
// @Tags films
// @Produce json
// @Param genres query string false "Список ID жанров через запятую"
// @Success 200 {array} dt.FilmDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /films [get]
func GetAllFilmsHandler(c *gin.Context) {
	var genreIDs []uint

	genresQuery := c.Query("genres")
	if genresQuery != "" {
		parts := strings.Split(genresQuery, ",")
		for _, p := range parts {
			val, err := strconv.ParseUint(strings.TrimSpace(p), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, dt.ErrorResponse{
					Code:    "INVALID_INPUT",
					Message: "invalid genre ID in query",
				})
				return
			}
			genreIDs = append(genreIDs, uint(val))
		}
	}

	films, err := services.GetAllFilms(genreIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, films)
}
