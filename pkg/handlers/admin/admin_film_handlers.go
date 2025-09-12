package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// CreateFilmHandler godoc
// @Summary Создать фильм
// @Tags admin-films
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dt.CreateFilmDTI true "Фильм"
// @Success 201 {object} dt.CreateFilmDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/films [post]
func CreateFilmHandler(c *gin.Context) {
	var input dt.CreateFilmDTI
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	dto, err := services.CreateFilm(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	// возвращаем только ID
	c.JSON(http.StatusCreated, dto)
}

// UpdateFilmHandler godoc
// @Summary Обновить фильм
// @Tags admin-films
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID фильма"
// @Param input body object true "Поля для обновления"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/films/{id} [patch]
func UpdateFilmHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "film ID not found",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	if err := services.UpdateFilm(uint(id), updates); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.ServAnswerDTO{
		Answer: "фильм обновлён",
	})
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
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid film ID",
		})
		return
	}

	if err := services.DeleteFilm(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "фильм удалён"})
}
