package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateSessionHandler godoc
// @Summary Создать новый сеанс
// @Tags admin-sessions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dt.CreateSessionDTI true "Данные для создания сеанса"
// @Success 201 {object} dt.CreateSessionDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/sessions [post]
func CreateSessionHandler(c *gin.Context) {
	var input dt.CreateSessionDTI
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	dto, err := services.CreateSession(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dt.CreateSessionDTO{ID: dto.ID})
}

// UpdateSessionHandler godoc
// @Summary Обновить сеанс (частично)
// @Tags admin-sessions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID сеанса"
// @Param input body object true "Поля для обновления"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/sessions/{id} [patch]
func UpdateSessionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid session ID",
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

	if err := services.UpdateSession(uint(id), updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "session not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.ServAnswerDTO{
		Answer: "сеанс обновлён",
	})
}

// DeleteSessionHandler godoc
// @Summary Удалить сеанс
// @Tags admin-sessions
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID сеанса"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/sessions/{id} [delete]
func DeleteSessionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid session ID",
		})
		return
	}

	if err := services.DeleteSession(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "session not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dt.ServAnswerDTO{
		Answer: "сеанс удалён",
	})
}
