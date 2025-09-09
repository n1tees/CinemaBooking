package handlers

import (
	"net/http"
	"strconv"
	"time"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// CreateSessionHandler godoc
// @Summary Создать новый сеанс
// @Tags admin-sessions
// @Security BearerAuth
// @Accept json
// @Produce json
//
//	@Param input body struct {
//	  FilmID uint      `json:"film_id" binding:"required"`
//	  HallID uint      `json:"hall_id" binding:"required"`
//	  Start  time.Time `json:"start" binding:"required"`
//	  Price  float64   `json:"price" binding:"required"`
//	} true "Данные для создания сеанса"
//
// @Success 201 {object} models.Session
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/sessions [post]
func CreateSessionHandler(c *gin.Context) {
	var input struct {
		FilmID uint      `json:"film_id" binding:"required"`
		HallID uint      `json:"hall_id" binding:"required"`
		Start  time.Time `json:"start" binding:"required"`
		Price  float64   `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := services.CreateSession(input.FilmID, input.HallID, input.Start, input.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// UpdateSessionHandler godoc
// @Summary Обновить сеанс (частично)
// @Tags admin-sessions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID сеанса"
// @Param input body map[string]interface{} true "Поля для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/sessions/{id} [patch]
func UpdateSessionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateSession(uint(id), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "сеанс обновлен"})
}

// DeleteSessionHandler godoc
// @Summary Удалить сеанс
// @Tags admin-sessions
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID сеанса"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/sessions/{id} [delete]
func DeleteSessionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	if err := services.DeleteSession(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "сеанс удален"})
}
