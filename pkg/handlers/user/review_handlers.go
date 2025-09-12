package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetReviewsByFilmHandler godoc
// @Summary Получить отзывы по фильму
// @Tags reviews
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {array} dt.ReviewDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /films/{id}/reviews [get]
func GetReviewsByFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid film ID",
		})
		return
	}

	reviews, err := services.GetReviewsByFilm(uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	var result []dt.ReviewDTO
	for _, r := range reviews {
		result = append(result, dt.ReviewDTO{
			UserID:  r.UserID,
			Rating:  r.Rating,
			Comment: r.Coment,
		})
	}

	c.JSON(http.StatusOK, result)
}

// AddReviewHandler godoc
// @Summary Добавить отзыв (на модерацию)
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID фильма"
// @Param input body dt.CreateReviewDTI true "Отзыв"
// @Success 201 {object} dt.CreateReviewDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /films/{id}/reviews [post]
func AddReviewHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid film ID",
		})
		return
	}

	var input dt.CreateReviewDTI
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, dt.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "user not found in context",
		})
		return
	}

	if err := services.AddReview(userID.(uint), uint(filmID), input.Rating, input.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dt.CreateReviewDTO{
		Status: "отзыв отправлен на модерацию",
	})
}

// GetFilmRatingHandler godoc
// @Summary Получить рейтинг фильма
// @Tags reviews
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {object} dt.FilmRatingDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /films/{id}/rating [get]
func GetFilmRatingHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid film ID",
		})
		return
	}

	rating, err := services.GetFilmRating(uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	if rating == nil {
		c.JSON(http.StatusOK, dt.FilmRatingDTO{
			Rating:  nil,
			Message: "нет отзывов",
		})
		return
	}

	c.JSON(http.StatusOK, dt.FilmRatingDTO{
		Rating: rating,
	})
}
