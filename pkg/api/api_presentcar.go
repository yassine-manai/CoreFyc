package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/config"
	"fyc/functions"
	"fyc/pkg/db"
)

// GetAllPresentCars godoc
//
//	@Summary		Get all present cars
//	@Description	Get a list of all present cars if not date provided it will return all present cars that entered this day
//	@Tags			PresentCars
//	@Produce		json
//	@Param			extra	query	string	false	"Include extra information if 'yes'"
//	@Param			lpn		query	string	false	"License Plate Number"
//	@Success		200		{array}	db.PresentCar
//	@Router			/fyc/presentcars [get]
func GetPresentCarsAPI(c *gin.Context) {
	ctx := context.Background()
	extra_req := c.DefaultQuery("extra", "false")
	lpn := c.DefaultQuery("lpn", "")

	//lang := c.DefaultQuery("lang", "EN")
	/////////////////// must check extra query with lpn
	if strings.ToLower(extra_req) == "true" || strings.ToLower(extra_req) == "1" || strings.ToLower(extra_req) == "yes" {
		cars, err := db.GetAllPresentExtra(ctx)
		if err != nil {
			log.Err(err).Msg("Error getting all present cars with extra data ")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all present cars with extra data ",
				"code":    10,
			})
			return
		}

		if len(cars) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No present cars found",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, cars)
		return
	}

	if strings.ToLower(lpn) != "" {
		cars, err := db.GetPresentCarByLPNs(ctx, lpn)
		if err != nil {
			log.Err(err).Msg("Error getting present cars with licence Plate ")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting present cars with licence Plate  ",
				"code":    10,
			})
			return
		}
		if len(cars) == 0 {
			log.Debug().Str("Licence Plate", lpn).Msg("No Data found for licence plate")

			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No present cars found",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, cars)
		return
	}

	Pcars, err := db.GetAllPresentCars(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all present cars")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all present cars",
			"code":    10,
		})
		return
	}

	if len(Pcars) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No present cars found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, Pcars)
}

// CreatePresentCar godoc
//
//	@Summary		Add a new present car
//	@Description	Add a new present car to the database
//	@Description	Default Date Format :YYYY-MM-DD HH:mm:ss
//	@Tags			PresentCars
//	@Accept			json
//	@Produce		json
//	@Param			presentCar	body		db.PresentCar	true	"Present Car data"
//	@Success		201			{object}	db.PresentCar
//	@Router			/fyc/presentcars [post]
func CreatePresentCarAPI(c *gin.Context) {
	var car db.PresentCar
	ctx := context.Background()

	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if !functions.Contains(db.Zonelist, *car.CurrZoneID) {
		log.Debug().Int("LastZoneID does not exist with ID ", *car.CurrZoneID).Msg("LastZoneID not found")

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "CurrZoneID not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", *car.CurrZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(db.Zonelist, *car.LastZoneID) {
		log.Debug().Int("LastZoneID does not exist with ID ", *car.LastZoneID).Msg("LastZoneID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "LastZoneID not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", *car.LastZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(db.CameraList, car.CameraID) {
		log.Debug().Int("Camera does not exist with ID ", car.CameraID).Msg("Camera ID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Camera ID not found ",
			"message": fmt.Sprintf("Camera with ID %d does not exist", car.CameraID),
			"code":    9,
		})
		return
	}

	log.Debug().Str("Body Date&Time:", car.TransactionDate)
	formattedDateTime, err := config.FormatDateTime(car.TransactionDate)
	log.Debug().Str("Formatted date and time:", formattedDateTime)

	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}
	car.TransactionDate = formattedDateTime

	// Insert the new car into the database
	if err := db.CreatePresentCar(ctx, &car); err != nil {
		log.Err(err).Msg("Error creating present car")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create present car",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	c.JSON(http.StatusCreated, car)
}

// UpdatePresentCarById godoc
//
//	@Summary		Update a present car by ID
//	@Description	Update an existing present car by ID
//	@Description	Default Date Format :YYYY-MM-DD HH:mm:ss
//	@Tags			PresentCars
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int				true	"Present Car ID"
//	@Param			presentCar	body		db.PresentCar	true	"Updated present car data"
//	@Success		200			{object}	db.PresentCar
//	@Router			/fyc/presentcars/{id} [put]
func UpdatePresentCarByIdAPI(c *gin.Context) {
	// Convert ID param to integer
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var updates db.PresentCar
	// Bind JSON payload
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if *updates.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the ID in the param",
			"code":    13,
		})
		return
	}

	if !functions.Contains(db.Zonelist, *updates.CurrZoneID) {
		log.Debug().Msg("CurrZoneID not found")

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "CurrZoneID not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", *updates.CurrZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(db.Zonelist, *updates.LastZoneID) {
		log.Debug().Msg("LastZoneID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "LastZoneID not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", *updates.LastZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(db.CameraList, updates.CameraID) {
		log.Debug().Msg("Camera ID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Camera ID not found ",
			"message": fmt.Sprintf("Camera with ID %d does not exist", updates.CameraID),
			"code":    9,
		})
		return
	}

	log.Debug().Str("Body Date&Time:", updates.TransactionDate)
	formattedDateTime, err := config.FormatDateTime(updates.TransactionDate)
	log.Debug().Str("Formatted date and time:", formattedDateTime)

	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	updates.TransactionDate = formattedDateTime

	ctx := context.Background()
	rowsAffected, err := db.UpdatePresentCar(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating present car by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update present car",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No present car found with the specified ID",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Present car modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// UpdatePresentCarByLpn godoc
//
//	@Summary		Update a present car by LPN
//	@Description	Update an existing present car by lpn
//	@Description	Default Date Format :YYYY-MM-DD HH:mm:ss
//	@Tags			PresentCars
//	@Accept			json
//	@Produce		json
//	@Param			lpn			query		string			true	"string default"	default(A)
//	@Param			presentCar	body		db.PresentCar	true	"Updated present car data by lpn"
//	@Success		200			{object}	db.PresentCar
//	@Router			/fyc/presentcars [put]
func UpdatePresentCarBylpnAPI(c *gin.Context) {

	lpn := c.Query("lpn")
	log.Info().Str("lpn", lpn).Msg("provided lpn parameters ")
	var updates db.PresentCar

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if !functions.Contains(db.Zonelist, *updates.CurrZoneID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Current ZoneID not found",
			"message": fmt.Sprintf("Zone ID %v does not exist", *updates.CurrZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(db.Zonelist, *updates.LastZoneID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Last ZoneID not found",
			"message": fmt.Sprintf("Zone with ID %v does not exist", *updates.LastZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(db.CameraList, updates.CameraID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Camera not found",
			"message": fmt.Sprintf("Camera with ID %v does not exist", updates.CameraID),
			"code":    9,
		})
		return
	}

	log.Debug().Str("Body Date&Time:", updates.TransactionDate)
	formattedDateTime, err := config.FormatDateTime(updates.TransactionDate)
	log.Debug().Str("Formatted date and time:", formattedDateTime)

	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := db.UpdatePresentCarByLpn(ctx, lpn, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating present car")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update present car",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No present car found with the specified Licence Plate",
			"code":    9,
		})
		return
	}

	response := db.ResponsePC{
		ID:              updates.ID,
		CarDetailsID:    updates.CarDetailsID,
		CameraID:        updates.CameraID,
		Confidence:      updates.Confidence,
		CurrZoneID:      updates.CurrZoneID,
		LastZoneID:      updates.LastZoneID,
		Direction:       updates.Direction,
		LPN:             lpn,
		TransactionDate: updates.TransactionDate,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Present car modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      response,
	})

}

// DeletePresentCar godoc
//
//	@Summary		Delete a present car
//	@Description	Delete a present car by ID
//	@Tags			PresentCars
//	@Param			id	path		int		true	"Present Car ID"
//	@Success		200	{string}	string	"Present car deleted successfully"
//	@Success		200	{object}	string	"Success"
//	@Failure		400	{object}	string	"Bad Request"
//	@Failure		404	{object}	string	"Not Found"
//	@Router			/fyc/presentcars/{id} [delete]
func DeletePresentCarAPI(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := db.DeletePresentCar(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting present car")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete present car",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No present car found with the specified ID ------ affected rows 0 ",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Present car deleted successfully",
		"message": rowsAffected,
		"code":    8,
	})
}
