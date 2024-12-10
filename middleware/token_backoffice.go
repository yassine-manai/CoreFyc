package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/config"
)

var JwtKey = []byte(config.Configvar.App.JSecret)

// Claims struct to store JWT claims
type ClaimsBackOffice struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(username string, role string) (string, int, error) {
	var tk string = config.Configvar.App.TkTime

	// Convert string to int
	hours, _ := strconv.Atoi(tk)
	expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)

	claims := &ClaimsBackOffice{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, hours, nil
}

// TokenMiddleware checks the token validity and expiration
func TokenMiddlewareBackOffice() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			log.Warn().Msg("Authorization required ")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -3,
				"error":   "Token is required",
				"success": false,
			})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &ClaimsBackOffice{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			log.Warn().Str("Client ID", claims.Username).Msg("Unauthorized, you need to connect first!")

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    -3,
				"message": "Unauthorized, you need to connect first !",
			})
			c.Abort()
			return
		}

		if claims.ExpiresAt < time.Now().Unix() {
			log.Warn().Msg("Unauthorized, Token has expired!")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    -3,
				"message": "Unauthorized, Token has expired!",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
