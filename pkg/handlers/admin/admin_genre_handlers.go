package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/models"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// CreateGenreHandler godoc
// @Summary Создать жанр
// @Tags admin-genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.Genre true "Жанр"
// @Success 201 {object} models.Genre
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/genres [post]
func CreateGenreHandler(c *gin.Context) {
	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateGenre(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, genre)
}

// UpdateGenreHandler godoc
// @Summary Обновить жанр
// @Tags admin-genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID жанра"
// @Param input body map[string]interface{} true "Поля для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/genres/{id} [patch]
func UpdateGenreHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateGenre(uint(id), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "жанр обновлён"})
}

// DeleteGenreHandler godoc
// @Summary Удалить жанр
// @Tags admin-genres
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID жанра"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/genres/{id} [delete]
func DeleteGenreHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre ID"})
		return
	}

	if err := services.DeleteGenre(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "жанр удалён"})
}

// AssignGenreToFilmHandler godoc
// @Summary Привязать жанр к фильму
// @Tags admin-genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID фильма"
// @Param input body struct{ GenreID uint `json:"genre_id" binding:"required"` } true "ID жанра"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/films/{id}/genres [post]
func AssignGenreToFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	var input struct {
		GenreID uint `json:"genre_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.AssignGenreToFilm(uint(filmID), input.GenreID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "жанр привязан к фильму"})
}

// RemoveGenreFromFilmHandler godoc
// @Summary Убрать жанр у фильма
// @Tags admin-genres
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID фильма"
// @Param genre_id path int true "ID жанра"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/films/{id}/genres/{genre_id} [delete]
func RemoveGenreFromFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	genreIDStr := c.Param("genre_id")
	genreID, err := strconv.ParseUint(genreIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre ID"})
		return
	}

	if err := services.RemoveGenreFromFilm(uint(filmID), uint(genreID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "жанр убран у фильма"})
}
