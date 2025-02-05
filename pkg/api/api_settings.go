package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/functions"
	"fyc/pkg/db"
)

// GetSettings godoc
//
//	@Summary		Get settings by CarPark ID
//	@Description	Get settings by CarPark ID
//	@Tags			Settings
//	@Produce		json
//	@Param			carpark_id	query		int	false	"CarPark ID"
//	@Success		200			{object}	db.Settings
//	@Router			/fyc/settings [get]
func GetSettingsAPI(c *gin.Context) {
	carParkIDStr := c.Query("carpark_id")
	ctx := context.Background()

	if carParkIDStr != "" {
		carParkID, err := strconv.Atoi(carParkIDStr)
		if err != nil {
			log.Error().Err(err).Msg("Invalid CarPark ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid CarPark ID format",
				"message": "CarPark ID must be an integer",
				"code":    13,
			})
			return
		}
		log.Info().Int("carpark_id", carParkID).Msg("Fetching settings by CarPark ID")
		ctx := context.Background()
		settings, err := db.GetSettings(ctx, carParkID)
		if err != nil {
			log.Error().Err(err).Int("carpark_id", carParkID).Msg("Error retrieving settings")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Settings not found for the specified CarPark ID",
				"code":    9,
			})
			return
		}
		log.Info().Int("carpark_id", carParkID).Msg("Settings fetched successfully")
		c.JSON(http.StatusOK, settings)
		return
	}

	log.Info().Msg("Fetching all settings")
	settings, err := db.GetAllSettings(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving settings")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Settings not found",
			"code":    9,
		})
		return
	}

	log.Info().Msg("Settings fetched successfully")
	c.JSON(http.StatusOK, settings)
}

// AddSettings godoc
//
//	@Summary		Add new settings
//	@Description	Add new settings for a CarPark
//	@Tags			Settings
//	@Accept			json
//	@Produce		json
//	@Param			settings	body		db.Settings	true	"Settings data"
//	@Success		201			{object}	db.Settings
//	@Router			/fyc/settings [post]
func AddSettingsAPI(c *gin.Context) {
	var settings db.Settings
	ctx := context.Background()

	log.Info().Msg("Attempting to add new settings")

	if err := c.ShouldBindJSON(&settings); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for settings creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	settings.DefaultLang = strings.ToLower(settings.DefaultLang)
	_, err := functions.DecodeBase64ToByteArray(settings.AppLogo)
	if err != nil {
		log.Err(err).Msg("Error converting App Logo")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	//fmt.Println(len(App_Logo))

	if err := db.CreateSettings(ctx, &settings); err != nil {
		log.Error().Err(err).Msg("Error creating settings")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create settings",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Int("carpark_id", settings.CarParkID).Msg("Settings created successfully")
	c.JSON(http.StatusCreated, settings)
}

// UpdateSettings godoc
//
//	@Summary		Update settings by CarPark ID
//	@Description	Update an existing settings by CarPark ID
//	@Tags			Settings
//	@Accept			json
//	@Produce		json
//	@Param			settings	body		db.SettingsNoBind	true	"Updated settings data"
//	@Success		200			{object}	db.SettingsNoBind
//	@Router			/fyc/settings [put]
func UpdateSettingsAPI(c *gin.Context) {

	var settings db.SettingsNoBind
	if err := c.ShouldBindJSON(&settings); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for settings update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	_, err := functions.DecodeBase64ToByteArray(settings.AppLogo)
	if err != nil {
		log.Err(err).Msg("Error converting App Logo")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	settings.DefaultLang = strings.ToLower(settings.DefaultLang)
	ctx := context.Background()
	err = db.UpdateSettings(ctx, settings, settings.CarParkID)
	if err != nil {
		if err.Error() == "no rows updated" {
			log.Warn().Int("carpark_id", settings.CarParkID).Msg("No settings found to update")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No settings found with the specified CarPark ID",
				"code":    9,
			})
			return
		}
		log.Error().Err(err).Int("carpark_id", settings.CarParkID).Msg("Error updating settings")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update settings",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Int("carpark_id", settings.CarParkID).Msg("Settings updated successfully")
	c.JSON(http.StatusOK, settings)
}
