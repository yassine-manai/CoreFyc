package export

import (
	"context"
	"fmt"
	"net/http"

	"fyc/pkg/db"

	arabic "github.com/abdullahdiaa/garabic"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary		Export Zone Data
// @Description	Export the zone data in PDF or Excel format based on the `file_type` query parameter
// @Tags			BackOffice - Export
// @Accept			json
// @Produce		application/pdf, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuthBackOffice
// @Param			file_type	query		string	false	"The type of the export file"	Enums(pdf, excel)	default(pdf)
// @Param			zone_ids	body		[]int	true	"The list of zone IDs"
// @Success		200			{string}	string	"Export successful"
// @Failure		500			{string}	string	"Internal Server Error"
// @Router			/backoffice/export_zone [post]
func ExportZone(c *gin.Context) {
	log.Debug().Msg(" / / / / # Exporting Zone Data # / / / /  ")
	ctx := context.Background()
	fileType := c.DefaultQuery("file_type", "pdf")

	var (
		zonesIDs []int
		zones    []db.ResponseZone
		err      error
	)

	if err := c.ShouldBindJSON(&zonesIDs); err != nil {
		log.Warn().Msg("Error binding request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"code":    -1,
		})
		return
	}

	log.Debug().Interface("ZoneIDs Passed in param", zonesIDs).Send()

	if len(zonesIDs) > 0 {
		for _, zoneID := range zonesIDs {
			zone, err := db.GetZoneByIDExport(ctx, zoneID)
			if err != nil {
				log.Warn().Int("Zone ID", zoneID).Msg("Zone not found")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": fmt.Sprintf("Zone %v not found", zoneID),
					"code":    -4,
				})
				return
			}
			zones = append(zones, *zone)
		}
	} else {
		zones, err = db.GetZoneByStatus(ctx, "")
		if err != nil {
			if err.Error() == "no rows found" {
				log.Warn().Str("Status", "All data").Msg("No zone found, returning all the data")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "No Zones found",
					"code":    -4,
				})
				return
			}
			log.Debug().Msgf("Error fetching data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching Zone data"})
			return
		}
	}

	var (
		data    [][]string
		headers = []string{"Zone ID", "Zone Name EN", "Zone Name AR", "Max Capacity", "Free Capacity", "Last Updated", "Status"}
		widths  = []float64{20, 50, 65, 30, 30, 60, 25}
	)

	for _, zone := range zones {
		zoneNameAr := zone.Name.Ar
		zoneNameEn := zone.Name.En

		if zoneNameAr == "" {
			zoneNameAr = "No Arabic Zone Name"
		} else {
			zoneNameAr = arabic.Shape(zoneNameAr)
		}
		if zoneNameEn == "" {
			zoneNameEn = "No English Zone Name"
		}

		status := "Disabled"
		if zone.IsEnabled {
			status = "Enabled"
		}

		data = append(data, []string{
			fmt.Sprintf("%d", *zone.ZoneID),
			zoneNameEn,
			zoneNameAr,
			fmt.Sprintf("%d", *zone.MaxCapacity),
			fmt.Sprintf("%d", *zone.FreeCapacity),
			zone.LastUpdated,
			status,
		})
	}

	switch fileType {
	case "excel":
		ExportToExcel(c, data, headers, "Zone_data")
	default:
		ExportToPDF(c, "L", data, headers, widths, "Zone Data Export", "Zone_data", "./font/Cairo-Regular.ttf")
	}
}
