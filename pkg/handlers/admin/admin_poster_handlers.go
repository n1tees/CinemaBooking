package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/models"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// CreatePosterHandler godoc
// @Summary Создать афишу
// @Tags admin-posters
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.Poster true "Данные афиши"
// @Success 201 {object} models.Poster
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/posters [post]
func CreatePosterHandler(c *gin.Context) {
	var poster models.Poster
	if err := c.ShouldBindJSON(&poster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreatePoster(&poster); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, poster)
}

// UpdatePosterHandler godoc
// @Summary Обновить афишу
// @Tags admin-posters
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID афиши"
// @Param input body map[string]interface{} true "Поля для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/posters/{id} [patch]
func UpdatePosterHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poster ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdatePoster(uint(id), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "афиша обновлена"})
}

// DeletePosterHandler godoc
// @Summary Удалить афишу
// @Tags admin-posters
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID афиши"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/posters/{id} [delete]
func DeletePosterHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid poster ID"})
		return
	}

	if err := services.DeletePoster(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "афиша удалена"})
}
