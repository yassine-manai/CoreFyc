package export

import (
	"context"
	"fmt"
	"fyc/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary		Export Cars Data
// @Description	Export the present cars data in PDF or Excel format based on the `file_type` query parameter
// @Tags			BackOffice - Export
// @Accept			json
// @Produce		application/pdf, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuthBackOffice
// @Param			license_plates	body		[]string	true	"The list of Lincense Plates"
// @Param			file_type		query		string		false	"The type of the export file"	Enums(pdf, excel)	default(pdf)
// @Success		200				{string}	string		"Export successful"
// @Failure		500				{string}	string		"Internal Server Error"
// @Router			/backoffice/export_cars [post]
func ExportCars(c *gin.Context) {
	log.Debug().Msg(" / / / / # Exporting Present Cars Data # / / / /  ")
	ctx := context.Background()

	// Query params
	fileType := c.DefaultQuery("file_type", "pdf")
	var CarsLpn []string
	var presents []db.PresentCar

	// Fetch sign_ids from body
	if err := c.ShouldBindJSON(&CarsLpn); err != nil {
		log.Warn().Msg("Error binding request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"code":    -1,
		})
		return
	}

	// Fetch signs based on the sign_ids (if provided)
	if len(CarsLpn) > 0 {
		for _, license := range CarsLpn {
			car, err := db.GetCarsLPN(ctx, license)
			if err != nil {
				log.Warn().Str("Licence Plate", license).Msg("Car not found")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": fmt.Sprintf("Car %v not found", license),
					"code":    -4,
				})
				return
			}
			presents = append(presents, *car)
		}
	} else {
		var err error
		presents, err = db.GetAllPresentExtra(ctx)
		if err != nil {
			if err.Error() == "no rows found" {
				log.Warn().Msg("No Car found")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "No Car found",
					"code":    -4,
				})
				return
			}
			log.Debug().Msgf("Error fetching data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching sign data"})
			return
		}
	}

	// Process data for export
	headers := []string{"License Plate", "Zone ID", "Zone Name", "Confidence", "Transaction Date"}
	widths := []float64{40, 30, 50, 25, 50}
	data := [][]string{}

	for _, prescar := range presents {

		CurZone, err := db.GetZoneByID(ctx, *prescar.CurrZoneID)
		if err != nil || CurZone == nil {
			log.Error().Err(err).Int("zoneID", *prescar.CurrZoneID).Msg("Current Zone not found")
			CurZone = &db.Zone{
				ZoneID: 0,
				Name:   make(map[string]interface{}),
			}
		}

		data = append(data, []string{
			prescar.LPN,
			fmt.Sprintf("%v", *prescar.CurrZoneID),
			fmt.Sprintf("%v", CurZone.Name["en"]),
			fmt.Sprintf("%d", *prescar.Confidence),
			prescar.TransactionDate,
		})
	}

	switch fileType {
	case "excel":
		ExportToExcel(c, data, headers, "cars_data")
	default:
		ExportToPDF(c, "P", data, headers, widths, "cars_data", "Present Cars Data Export", "./font/Cairo-Regular.ttf")
	}

}
