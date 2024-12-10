package third_party

import (
	"fmt"
	"fyc/pkg/db"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

func isClientExist(clientList map[string]db.ClientDetails, clientID string) (bool, *db.ClientDetails) {
	if clientData, exists := clientList[clientID]; exists {
		return true, &clientData
	}
	return false, nil
}

func Extract_token_data(authHeader string) (bool, string, error) {
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return false, "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fuzzyLogic, fuzzyExists := claims["fuzzy_logic"].(bool)
		clientID, clientIDExists := claims["client_id"].(string)

		// Check fields are found
		if !fuzzyExists {
			return false, "", fmt.Errorf("fuzzy_logic not found in token")
		}
		if !clientIDExists {
			return false, "", fmt.Errorf("client_id not found in token")
		}

		log.Debug().Str("client_id", clientID).Bool("fuzzy_logic", fuzzyLogic).Msg("Client ID and Fuzzy Logic found")
		return fuzzyLogic, clientID, nil
	}
	return false, "", fmt.Errorf("could not parse claims from token")
}

/* func CheckToken(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    -5,
			"message": "Invalid request. 'Token' parameter is required.",
		})
		return false
	}

	// Extract the token from the "Bearer <token>" format
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != AccessTokenFake {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    -3,
			"message": "Unauthorized, you need to connect first !",
		})
		return false
	}

	return true
} */
