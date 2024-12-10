package middleware

import (
	"context"
	"fyc/config"
	"fyc/pkg/db"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ClaimsThirdParty struct {
	ClientID        string `json:"client_id"`
	ClientSecret    string `json:"client_secret"`
	ClientGrantType string `json:"grant_type"`
	FuzzyLogic      bool   `json:"fuzzy_logic"`
	jwt.StandardClaims
}

func GenerateTokenThirdParty(client_id, client_secret, client_grantType string, client_fuzzy_logic bool) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &ClaimsThirdParty{
		ClientID:        client_id,
		ClientSecret:    client_secret,
		ClientGrantType: client_grantType,
		FuzzyLogic:      client_fuzzy_logic,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckTokenInDB(ctx context.Context, client_id, tokenString string) bool {
	client, err := db.GetClientCredById(ctx, client_id)
	if err != nil {
		log.Warn().Err(err).Str("Client_ID", client_id).Msg("Error retrieving Client ID")
		return false
	}

	if client.ApiKey != nil && *client.ApiKey == tokenString {
		log.Info().Str("Client_ID", client_id).Msg("API KEY CORRECT ")
		return true
	}

	return false
}

func TokenMiddlewareThirdParty() gin.HandlerFunc {
	tokenCheck := config.Configvar.App.TokenCheck
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			log.Warn().Msg("Authorization required")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -3,
				"error":   "Token is required",
				"success": false,
			})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &ClaimsThirdParty{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			log.Warn().Str("Client ID", claims.ClientID).Msg("Unauthorized, you need to connect first!")
			//fmt.Printf("pre %v", tokenPref)

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    -3,
				"message": "Unauthorized, you need to connect first!",
			})
			c.Abort()
			return
		}

		if tokenCheck == "true" {
			log.Info().Str("Token Check ", tokenCheck).Msg("Making Token Check")
			ctx := c.Request.Context()
			if !CheckTokenInDB(ctx, claims.ClientID, tokenString) {
				log.Warn().Str("Client ID", claims.ClientID).Msg("Unauthorized, invalid token in DB!")
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"code":    -3,
					"message": "Unauthorized, invalid token!",
				})
				c.Abort()
				return
			}
		}

		if claims.ExpiresAt < time.Now().Unix() {
			log.Warn().Str("Client ID", claims.ClientID).Msg("Unauthorized, Token has expired!")
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
