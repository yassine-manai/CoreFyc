package backoffice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/config"
	"fyc/middleware"
	"fyc/pkg/db"
)

type User struct {
	Username string
	Password string
}

//var jwtKey = []byte(config.Configvar.App.JSecret)

// LoginUser godoc
//
//	@Summary		User Login
//	@Description	Login for users to access the system
//	@Tags			Backoffice - Login
//	@Accept			json
//	@Produce		json
//	@Param			User	body	User	true	"User credentials"
//	@Router			/backoffice/login [post]
func Login(c *gin.Context) {

	ctx := context.Background()
	var tokenPref = config.Configvar.App.TokenPrefBackoffice

	var input struct {
		Username string `json:"username" binding:"required" example:"admin"`
		Password string `json:"password" binding:"required" example:"admin"`
	}

	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    -500,
				"message": "An unexpected error occurred. Please try again later.",
			})
		}
	}()

	//db.Token = ""

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn().Err(err).Msg("Error Getting Data")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    -5,
		})
		return
	}

	log.Debug().Str("Username", input.Username).Msg("Login Operation")

	userFound, err := db.GetUserByUsername(ctx, input.Username)
	if err != nil || userFound.UserName != input.Username || userFound.Password != input.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid username or password",
			"code":    -2,
			"message": "Invalid credentials",
		})
		return
	}

	if !userFound.IsEnabled {
		log.Debug().Str("User Disabled", input.Username).Msg("User Disabled")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User is disabled",
			"code":    -2,
		})
		return
	}

	token, hr, err := middleware.GenerateToken(userFound.UserName, userFound.Role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Could not generate token",
			"code":    -500,
		})
		return
	}

	log.Info().Str("User ", input.Username).Int("Time", hr).Msg("Connected successfully")
	//db.Token = token

	if tokenPref == "true" {
		log.Info().Msg("Passing token with BEARER Prefix")
		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"first_name": userFound.FirstName,
			"last_name":  userFound.LastName,
			"role":       userFound.Role,
			"token":      fmt.Sprint("Bearer " + token),
			"message":    "Connected successfully",
		})
		return
	}

	log.Info().Msg("Passing token without BEARER Prefix")
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"first_name": userFound.FirstName,
		"last_name":  userFound.LastName,
		"role":       userFound.Role,
		"token":      token,
		"message":    "Connected successfully",
	})
}
