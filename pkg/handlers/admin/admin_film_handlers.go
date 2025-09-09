package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/models"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// CreateFilmHandler godoc
// @Summary Создать фильм
// @Tags admin-films
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.Film true "Фильм"
// @Success 201 {object} models.Film
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/films [post]
func CreateFilmHandler(c *gin.Context) {
	var film models.Film
	if err := c.ShouldBindJSON(&film); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateFilm(&film); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, film)
}

// UpdateFilmHandler godoc
// @Summary Обновить фильм
// @Tags admin-films
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID фильма"
// @Param input body map[string]interface{} true "Поля для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/films/{id} [patch]
func UpdateFilmHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateFilm(uint(id), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "фильм обновлён"})
}

// DeleteFilmHandler godoc
// @Summary Удалить фильм
// @Tags admin-films
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/films/{id} [delete]
func DeleteFilmHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	if err := services.DeleteFilm(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "фильм удалён"})
}
