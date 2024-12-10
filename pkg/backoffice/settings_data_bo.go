package backoffice

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/functions"
	"fyc/pkg/cron"
	"fyc/pkg/db"
)

type ResponseData struct {
	General GeneralInfo `json:"general"`
	Cron    CronInfo    `json:"cron"`
	Kiosk   KioskInfo   `json:"Kiosk"`
}

type GeneralInfo struct {
	CarParkID    int                    `json:"carpark_id"`
	CarParkName  map[string]interface{} `json:"carpark_name" swaggertype:"object"`
	DefaultLang  string                 `json:"default_lang"`
	PkaImageSize string                 `json:"pka_image_size"`
}

type CronInfo struct {
	FycCleanCron      int  `json:"fyc_clean_cron"`
	CountingCleanCron int  `json:"counting_clean_cron"`
	IsFycEnabled      bool `json:"is_fyc_enabled"`
	IsCountingEnabled bool `json:"is_counting_enabled"`
}

type KioskInfo struct {
	TimeOutScreenKiosk int    `json:"timeout_screenKiosk"`
	AppLogo            string `json:"app_logo"`
	TC                 string `json:"TC"`
}

// GetSettings godoc
//
//	@Summary		Get settings
//	@Description	Get settings
//	@Tags			Backoffice - Settings
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Success		200	{object}	db.Settings
//	@Router			/backoffice/getSettings [get]
func GetSettingsDataAPI(c *gin.Context) {
	ctx := context.Background()

	log.Info().Msg("Fetching all settings")
	settings, err := db.GetAllSettings(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving settings")
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	responseData := []map[string]interface{}{
		{
			"general": map[string]interface{}{
				"carpark_id":     settings.CarParkID,
				"carpark_name":   settings.CarParkName,
				"default_lang":   settings.DefaultLang,
				"pka_image_size": settings.PkaImageSize,
			},
			"cron": map[string]interface{}{
				"fyc_clean_cron":      settings.FycCleanCron,
				"counting_clean_cron": settings.CountingCleanCron,
				"is_counting_enabled": settings.IsCountingEnabled,
				"is_fyc_enabled":      settings.IsFycEnabled,
			},
			"Kiosk": map[string]interface{}{
				"timeout_screenKiosk": settings.TimeOutScreenKisok,
				"app_logo":            settings.AppLogo,
				"tc":                  settings.TC,
			},
		},
	}

	log.Info().Msg("Settings fetched successfully")
	c.JSON(http.StatusOK, responseData)

}

// UpdateSettings godoc
//
//	@Summary		Update settings
//	@Description	Update an existing settings
//	@Tags			Backoffice - Settings
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			carpark_id	query		int				true	"CarPark ID"
//	@Param			settings	body		ResponseData	true	"Updated settings data"
//	@Success		200			{object}	ResponseData
//	@Router			/backoffice/updateSettings [put]
func UpdateSettingsDataAPI(c *gin.Context) {
	carpark_id := c.Query("carpark_id")

	if carpark_id == "" {
		log.Warn().Msg("CarPark ID is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "CarPark ID is required",
			"message": "Please provide an CarPark ID to update",
			"code":    12,
		})
		return
	}

	var set2update ResponseData
	var settings db.SettingsNoBind

	cp_id, err := strconv.Atoi(carpark_id)
	if err != nil {
		log.Err(err).Str("id", carpark_id).Msg("Invalid carpark ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Carpark ID must be a valid integer",
			"code":    -5,
		})
		return
	}

	if err := c.ShouldBindJSON(&set2update); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for settings update")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -5,
		})
		return
	}

	_, err = functions.DecodeBase64ToByteArray(settings.AppLogo)
	if err != nil {
		log.Err(err).Msg("Error converting App Logo")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -5,
		})
		return
	}

	dataSetting := db.SettingsNoBind{
		CarParkID:    set2update.General.CarParkID,
		CarParkName:  set2update.General.CarParkName,
		DefaultLang:  set2update.General.DefaultLang,
		PkaImageSize: set2update.General.PkaImageSize,

		FycCleanCron:      set2update.Cron.FycCleanCron,
		CountingCleanCron: set2update.Cron.CountingCleanCron,

		IsFycEnabled:      &set2update.Cron.IsFycEnabled,
		IsCountingEnabled: &set2update.Cron.IsCountingEnabled,

		TimeOutScreenKisok: set2update.Kiosk.TimeOutScreenKiosk,
		AppLogo:            set2update.Kiosk.AppLogo,
		TC:                 set2update.Kiosk.TC,
	}

	ctx := context.Background()
	err = db.UpdateSettings(ctx, dataSetting, cp_id)
	if err != nil {
		if err.Error() == "no rows updated" {
			log.Warn().Int("carpark_id", cp_id).Msg("No settings found to update")
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "No settings found with the specified CarPark ID",
				"code":    -4,
			})
			return
		}
		log.Error().Err(err).Int("carpark_id", cp_id).Msg("Error updating settings")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -500,
		})
		return
	}

	go cron.CronFyc()
	go cron.CronCounting()
	log.Info().Int("carpark_id", cp_id).Msg("Settings updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Settings updated successfully",
	})
}
