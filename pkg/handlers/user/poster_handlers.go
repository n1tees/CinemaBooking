package handlers

import (
	"net/http"

	"CinemaBooking/pkg/dt"
	"CinemaBooking/pkg/services"

	"github.com/gin-gonic/gin"
)

// GetAllPostersHandler godoc
// @Summary Получить список афиш
// @Tags posters
// @Produce json
// @Success 200 {array} dt.PosterDTO
// @Failure 500 {object} dt.ErrorResponse
// @Router /posters [get]
func GetAllPostersHandler(c *gin.Context) {
	posters, err := services.GetAllPosters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dt.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	var result []dt.PosterDTO
	for _, p := range posters {
		result = append(result, dt.PosterDTO{
			FilmID: p.ID,
			URL:    p.ImageURL,
		})
	}

	c.JSON(http.StatusOK, result)
}
