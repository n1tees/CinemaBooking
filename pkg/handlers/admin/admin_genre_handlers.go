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

// CreateGenreHandler godoc
// @Summary Создать жанр
// @Tags admin-genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dt.CreateGenreDTI true "Жанр"
// @Success 201 {object} dt.CreateGenreDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/genres [post]
func CreateGenreHandler(c *gin.Context) {
	var input dt.CreateGenreDTI
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	id, err := services.CreateGenre(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dt.CreateGenreDTO{ID: id})
}

// UpdateGenreHandler godoc
// @Summary Обновить жанр
// @Tags admin-genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID жанра"
// @Param input body object true "Поля для обновления"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/genres/{id} [patch]
func UpdateGenreHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "genre ID not found",
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

	if err := services.UpdateGenre(uint(id), updates); err != nil {
		if err.Error() == "жанр не найден" {
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
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
		Answer: "жанр обновлён",
	})
}

// DeleteGenreHandler godoc
// @Summary Удалить жанр
// @Tags admin-genres
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID жанра"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/genres/{id} [delete]
func DeleteGenreHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid genre ID",
		})
		return
	}

	if err := services.DeleteGenre(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "genre not found",
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
		Answer: "жанр удалён",
	})
}

// AssignGenreToFilmHandler godoc
// @Summary Привязать жанр к фильму
// @Tags admin-genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID фильма"
// @Param input body dt.AssignGenreDTI true "ID жанра"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/films/{id}/genres [post]
func AssignGenreToFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid film ID",
		})
		return
	}

	var input dt.AssignGenreDTI
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	if err := services.AssignGenreToFilm(uint(filmID), input.GenreID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "film or genre not found",
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
		Answer: "жанр привязан к фильму",
	})
}

// RemoveGenreFromFilmHandler godoc
// @Summary Убрать жанр у фильма
// @Tags admin-genres
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID фильма"
// @Param genre_id path int true "ID жанра"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/films/{id}/genres/{genre_id} [delete]
func RemoveGenreFromFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid film ID",
		})
		return
	}

	genreIDStr := c.Param("genre_id")
	genreID, err := strconv.ParseUint(genreIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid genre ID",
		})
		return
	}

	if err := services.RemoveGenreFromFilm(uint(filmID), uint(genreID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "film or genre not found",
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
		Answer: "жанр убран у фильма",
	})
}
