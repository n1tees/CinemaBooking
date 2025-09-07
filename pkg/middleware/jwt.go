package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"CinemaBooking/config"
	"CinemaBooking/pkg/db"
	"CinemaBooking/pkg/models"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует токен"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTSecret()), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверная структура токена"})
			return
		}

		expUnix := int64(claims["exp"].(float64))
		if time.Now().Unix() > expUnix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Срок действия токена истёк"})
			return
		}

		userID := uint(claims["user_id"].(float64))
		userType := claims["user_type"].(string)

		// Проверка: есть ли такой пользователь в базе
		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
			return
		}

		c.Set("user_id", userID)
		c.Set("user_type", userType)
		c.Next()
	}
}
