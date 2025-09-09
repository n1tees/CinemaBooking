package handlers

import (
	"net/http"
	"strconv"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// ApproveReviewHandler godoc
// @Summary Одобрить отзыв
// @Tags admin-reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/reviews/{id}/approve [patch]
func ApproveReviewHandler(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	if err := services.ApproveReview(uint(reviewID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "отзыв одобрен"})
}

// RejectReviewHandler godoc
// @Summary Отклонить отзыв
// @Tags admin-reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/reviews/{id}/reject [patch]
func RejectReviewHandler(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	if err := services.RejectReview(uint(reviewID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "отзыв отклонен"})
}

// DeleteOwnReviewHandler godoc
// @Summary Удалить свой отзыв
// @Tags reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reviews/{id} [delete]
func DeleteOwnReviewHandler(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	if err := services.DeleteOwnReview(userID.(uint), uint(reviewID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "отзыв удалён"})
}
