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

// ApproveReviewHandler godoc
// @Summary Одобрить отзыв
// @Tags admin-reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/reviews/{id}/approve [patch]
func ApproveReviewHandler(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid review ID",
		})
		return
	}

	if err := services.ApproveReview(uint(reviewID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "review not found",
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
		Answer: "отзыв одобрен",
	})
}

// RejectReviewHandler godoc
// @Summary Отклонить отзыв
// @Tags admin-reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /admin/reviews/{id}/reject [patch]
func RejectReviewHandler(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid review ID",
		})
		return
	}

	if err := services.RejectReview(uint(reviewID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "review not found",
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
		Answer: "отзыв отклонён",
	})
}

// DeleteOwnReviewHandler godoc
// @Summary Удалить свой отзыв
// @Tags reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} dt.ServAnswerDTO
// @Failure 400 {object} dt.ErrorResponse
// @Failure 401 {object} dt.ErrorResponse
// @Failure 403 {object} dt.ErrorResponse
// @Failure 404 {object} dt.ErrorResponse
// @Failure 500 {object} dt.ErrorResponse
// @Router /reviews/{id} [delete]
func DeleteOwnReviewHandler(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dt.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "invalid review ID",
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

	if err := services.DeleteOwnReview(userID.(uint), uint(reviewID)); err != nil {
		switch err.Error() {
		case "отзыв не найден":
			c.JSON(http.StatusNotFound, dt.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
		case "нельзя удалить чужой отзыв":
			c.JSON(http.StatusForbidden, dt.ErrorResponse{
				Code:    "FORBIDDEN",
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, dt.ServAnswerDTO{
		Answer: "отзыв удалён",
	})
}
