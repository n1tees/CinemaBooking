package handlers

import (
	"net/http"

	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetAllPostersHandler godoc
// @Summary Получить список афиш
// @Tags posters
// @Produce json
// @Success 200 {array} models.Poster
// @Failure 500 {object} map[string]string
// @Router /posters [get]
func GetAllPostersHandler(c *gin.Context) {
	posters, err := services.GetAllPosters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posters)
}
