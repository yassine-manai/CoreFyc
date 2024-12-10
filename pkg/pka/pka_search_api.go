package pka

import (
	"context"
	"fmt"
	"fyc/pkg/db"
	"fyc/pkg/third_party"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary		PKA SYSTEM API - Search Car
//
// @Description	Get ZONE with PKA SYSTEM
// @Tags			PKA - API
// @Produce		json
// @Param			visit.plate.text	query	string	true	"visit.plate.text"
// @Router			/v2/bays.json [get]
func PkaSearchAPI(c *gin.Context) {
	ctx := context.Background()
	visitPlate := c.Query("visit.plate.text")

	if visitPlate == "" {
		log.Error().Msg("No visit plate provided.")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No visit plate provided.",
			"success": false,
			"code":    -4,
		})
		return
	}

	// Fetch data
	car, err := db.GetPresentCarByLPN(ctx, visitPlate)
	if err != nil {
		log.Warn().Str("Error : ", err.Error()).Str("license_plate", visitPlate).Msg("Error retrieving car by LPN")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    -4,
			"message": "No data found !",
		})
		return
	}

	log.Info().Str("License Plate", visitPlate).Msg("Car found with license plate")
	log.Debug().Int("zone", *car.CurrZoneID).Msg("Last Zone ID")
	//log.Debug().Interface("zone", car).Msg("Present Data")

	spotID := *car.CurrZoneID
	visitPlate = car.LPN
	language := "en"

	log.Debug().Str("Spot ID", fmt.Sprint(spotID)).Str("License Plate", visitPlate).Msg("Zone Found")

	// Retrieve the zone image by ID and language
	zoneImage, err := db.GetZoneImageByZONEIDLang(ctx, spotID, language)
	if err != nil {
		log.Warn().Str("Error : ", err.Error()).Str("Language Provided", language).Int("Car Detail ID", *car.CarDetailsID).Msg("Error retrieving zone image")
		c.JSON(http.StatusOK, []third_party.CarLocation{})
		return
	}

	// Get Zone name
	zoneData, err := db.GetZoneByID(ctx, spotID)
	if err != nil {
		log.Warn().Str("Error: ", err.Error()).Int("Zone ID", *car.CurrZoneID).Msg("Error retrieving zone")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    -4,
			"message": "No data found !",
		})
		return
	}

	// Debug log for zone names
	log.Debug().Interface("Zone Names", zoneData.Name)

	if zoneName, ok := zoneData.Name[language]; ok {
		log.Info().Str("Language", language).Interface("Name", zoneName).Interface("Name", zoneData.Name[language]).Msg("Zone Name")
	} else {
		log.Warn().Msg("Zone name not found for the specified language")
	}

	log.Debug().Str("Found Picture for Car with license plate", visitPlate).Str("Picture Name", fmt.Sprint(zoneImage.ID))

	// Prepare the successful response

	response := third_party.CarLocation{
		ZoneName:     fmt.Sprint(zoneData.Name[language]),
		LicensePlate: visitPlate,
		SpotID:       fmt.Sprint(spotID),
		PictureName:  fmt.Sprint(zoneImage.ID),
	}

	log.Info().Str("License Plate", visitPlate).Msg("Data Found in PKA API")
	log.Debug().Interface("Data Found in PKA API", response).Send()

	c.JSON(http.StatusOK, gin.H{
		"id":                24,
		"is_in_violation":   false,
		"is_occupied":       true,
		"is_out_of_service": false,
		"is_reserved":       false,
		"map": gin.H{
			"id":   zoneData.ZoneID,         // to change
			"name": zoneData.Name[language], // to change
		},
		"position": gin.H{
			"x": 340.025390625,
			"y": 215.331237792969,
		},
		"visit": gin.H{
			"dwell":           "15:29:08.3940000",
			"entry_timestamp": car.TransactionDate, // to change
			"id":              5260304,
			"plate": gin.H{
				"confidence": 99,                  // to change
				"text":       visitPlate,          // to change
				"timestamp":  car.TransactionDate, // to change
			},
		},
		"zone": gin.H{
			"id":   zoneData.ZoneID,         // to change
			"name": zoneData.Name[language], // to change
		},
	})

}
