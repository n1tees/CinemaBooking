package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetReviewsByFilmHandler godoc
// @Summary Получить отзывы по фильму
// @Tags reviews
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {array} models.Review
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /films/{id}/reviews [get]
func GetReviewsByFilmHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	reviews, err := services.GetReviewsByFilm(uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// AddReviewHandler godoc
// @Summary Добавить отзыв (на модерацию)
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID фильма"
//
//	@Param input body struct {
//	  Rating  uint   `json:"rating" binding:"required"`
//	  Comment string `json:"comment" binding:"required"`
//	} true "Отзыв"
//
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /films/{id}/reviews [post]
func AddReviewHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	var input struct {
		Rating  uint   `json:"rating" binding:"required"`
		Comment string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	if err := services.AddReview(userID.(uint), uint(filmID), input.Rating, input.Comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "отзыв отправлен на модерацию"})
}

// GetFilmRatingHandler godoc
// @Summary Получить рейтинг фильма
// @Tags reviews
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /films/{id}/rating [get]
func GetFilmRatingHandler(c *gin.Context) {
	filmIDStr := c.Param("id")
	filmID, err := strconv.ParseUint(filmIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	rating, err := services.GetFilmRating(uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rating == nil {
		c.JSON(http.StatusOK, gin.H{"rating": nil, "message": "нет отзывов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rating": *rating})
}
