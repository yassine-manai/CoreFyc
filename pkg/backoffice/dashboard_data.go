package backoffice

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/functions"
	"fyc/pkg/db"
)

// @Summary		Get Dashboard Data
// @Description	Retrieve dashboard data including total cameras, zones, capacity, free spaces, signs, and present cars.
// @Tags			Backoffice - Dashboard
// @Security		BearerAuthBackOffice
// @Produce		json
// @Router			/backoffice/get_dashboard_data [get]
func GetDashboardData(c *gin.Context) {
	handler := functions.NewResponseHandler()
	ctx := context.Background()

	totalCameras := 0
	totalCapacity := 0
	totalFreeCapacity := 0
	totalPresentCars := 0
	totalSigns := 0
	totalZones := 0

	// Retrieve data from the database
	Camera, err := db.GetAllCameraExtra(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all cameras with extra data")
		handler.RespondWithStatus(c, -500)
		return
	}
	Zones, err := db.GetAllZoneNoExtra(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all zones with extra data")
		handler.RespondWithStatus(c, -500)
		return
	}
	Signs, err := db.GetAllSigns(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all signs with extra data")
		handler.RespondWithStatus(c, -500)
		return
	}
	PresentCars, err := db.GetAllPresentCars(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all present cars")
		handler.RespondWithStatus(c, -500)
		return
	}

	// Aggregate data
	for _, zone := range Zones {
		if zone.MaxCapacity != nil {
			totalCapacity += *zone.MaxCapacity
		}
		if zone.FreeCapacity != nil {
			totalFreeCapacity += *zone.FreeCapacity
		}
	}
	totalCameras = len(Camera)
	totalZones = len(Zones)
	totalSigns = len(Signs)
	totalPresentCars = len(PresentCars)

	log.Info().
		Int("TotalCameras", totalCameras).
		Int("TotalCapacity", totalCapacity).
		Int("TotalFreeSpaces", totalFreeCapacity).
		Int("TotalPresentCars", totalPresentCars).
		Int("TotalSigns", totalSigns).
		Int("TotalZones", totalZones).
		Msg("Dashboard data retrieved successfully")

	// Prepare response as an array
	response := []map[string]interface{}{
		{
			"success":            true,
			"total_cameras":      totalCameras,
			"total_capacity":     totalCapacity,
			"total_free_spaces":  totalFreeCapacity,
			"total_present_cars": totalPresentCars,
			"total_signs":        totalSigns,
			"total_zones":        totalZones,
		},
	}

	// Send response
	c.JSON(http.StatusOK, response)

}
