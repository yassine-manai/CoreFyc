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

// GetCarDetails godoc
//
//	@Summary		Get car details
//	@Description	Get a list of all car details or a specific car detail by ID
//	@Tags			Car Details
//	@Produce		json
//	@Param			id		query	int		false	"CarDetail ID"
//	@Param			extra	query	string	false	"Include extra information if 'yes'"
//	@Success		200		{array}	db.CarDetail
//	@Router			/fyc/carDetails [get]
func GetCarDetailsAPI(c *gin.Context) {
	log.Debug().Msg("GetCarDetailsAPI request")

	ctx := context.Background()
	idStr := c.Query("id")
	extraReq := strings.ToLower(c.DefaultQuery("extra", "")) // Normalize extra request

	log.Info().Str("extra", extraReq).Msg("Extra request parameter received")

	// Case 1: ID is not empty, extra is "yes" → Return data with ID and extra
	if idStr != "" && extraReq == "yes" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid carDetail ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "carDetail ID must be a valid integer",
				"code":    12,
			})
			return
		}

		log.Info().Str("carDetail", idStr).Msg("Fetching carDetail by ID with extra data")
		carDetail, err := db.GetCarDetailByIDExtra(ctx, id)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Error retrieving carDetail by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "CarDetail not found",
				"code":    9,
			})
			return
		}

		for i := range carDetail {
			carDetail[i].Image1, _ = functions.ByteaToBase64([]byte(carDetail[i].Image1))
			carDetail[i].Image2, _ = functions.ByteaToBase64([]byte(carDetail[i].Image2))
		}
		c.JSON(http.StatusOK, carDetail)
		return
	}

	// Case 2: ID is not empty, extra is empty → Return data with ID, no extra
	if idStr != "" && extraReq == "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid carDetail ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "carDetail ID must be a valid integer",
				"code":    12,
			})
			return
		}

		log.Info().Str("carDetail", idStr).Msg("Fetching carDetail by ID")
		carDetail, err := db.GetCarDetailByID(ctx, id)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Error retrieving carDetail by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "CarDetail not found",
				"code":    9,
			})
			return
		}

		for i := range carDetail {
			carDetail[i].Image1, _ = functions.ByteaToBase64([]byte(carDetail[i].Image1))
			carDetail[i].Image2, _ = functions.ByteaToBase64([]byte(carDetail[i].Image2))
		}
		c.JSON(http.StatusOK, carDetail)
		return
	}

	// Case 3: ID is empty, extra is "yes" → Return all data with extra
	if idStr == "" && (extraReq == "yes" || extraReq == "true" || extraReq == "1") {
		log.Info().Msg("Fetching all car details with extra data")
		carDetailExtra, err := db.GetAllCarDetailExtra(ctx)
		if err != nil {
			log.Err(err).Msg("Error getting all car details with extra data")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all car details with extra data",
				"code":    10,
			})
			return
		}

		if len(carDetailExtra) == 0 {
			log.Info().Msg("No car details found with extra data")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No car details found",
				"code":    9,
			})
			return
		}

		for i := range carDetailExtra {
			carDetailExtra[i].Image1, _ = functions.ByteaToBase64([]byte(carDetailExtra[i].Image1))
			carDetailExtra[i].Image2, _ = functions.ByteaToBase64([]byte(carDetailExtra[i].Image2))
		}
		c.JSON(http.StatusOK, carDetailExtra)
		return
	}

	// Case 4: ID is empty, extra is empty → Return all data without extra
	if idStr == "" && extraReq == "" {
		log.Info().Msg("Fetching all car details without extra data")
		carDet, err := db.GetAllCarDetail(ctx)
		if err != nil {
			log.Err(err).Msg("Error getting all car details")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all car details",
				"code":    10,
			})
			return
		}

		if len(carDet) == 0 {
			log.Info().Msg("No car details found")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No car details found",
				"code":    9,
			})
			return
		}

		for i := range carDet {
			carDet[i].Image1, _ = functions.ByteaToBase64([]byte(carDet[i].Image1))
			carDet[i].Image2, _ = functions.ByteaToBase64([]byte(carDet[i].Image2))
		}
		c.JSON(http.StatusOK, carDet)
		return
	}

	// If none of the conditions match, return bad request
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Invalid request",
		"message": "Invalid combination of ID and extra parameters",
		"code":    11,
	})
}

// CreateCarDetail godoc
//
//	@Summary		Add a new car detail
//	@Description	Add a new car detail to the database
//	@Tags			Car Details
//	@Accept			json
//	@Produce		json
//	@Param			CarDetail	body		db.CarDetail	true	"Car detail data"
//	@Success		201			{object}	db.CarDetail
//	@Router			/fyc/carDetails [post]
func CreateCarDetailAPI(c *gin.Context) {
	var carDetail db.CarDetail
	log.Debug().Msg("Creating CarDetail")

	if err := c.ShouldBindJSON(&carDetail); err != nil {
		log.Err(err).Msg("Invalid request payload for car detail creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	Image1Enc, err := functions.DecodeBase64ToByteArray(carDetail.Image1)
	if err != nil {
		log.Err(err).Msg("Error converting image 1")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 1",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	Image2Enc, err := functions.DecodeBase64ToByteArray(carDetail.Image2)
	if err != nil {
		log.Err(err).Msg("Error converting image 2")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 2",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	log.Info().Int("image1_size", len(Image1Enc)).Int("image2_size", len(Image2Enc)).Msg("Images decoded successfully")

	ctx := context.Background()
	if err := db.CreateCarDetail(ctx, &carDetail); err != nil {
		log.Err(err).Msg("Error creating new car detail")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new car detail",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Int("CarDetail", carDetail.ID).Msg("Car detail created successfully")
	c.JSON(http.StatusCreated, carDetail)
}

// UpdateCarDetailById godoc
//
//	@Summary		Update a car detail by ID
//	@Description	Update an existing car detail by ID
//	@Tags			Car Details
//	@Accept			json
//	@Produce		json
//	@Param			id			query		int				true	"Car ID"
//	@Param			CarDetail	body		db.CarDetail	true	"Updated car detail data"
//	@Success		200			{object}	db.CarDetail
//	@Failure		400			{object}	map[string]interface{}	"Invalid request"
//	@Failure		404			{object}	map[string]interface{}	"Car detail not found"
//	@Router			/fyc/carDetails [put]
func UpdateCarDetailByIdAPI(c *gin.Context) {
	log.Debug().Msg("Updating CarDetail")

	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid ID format for car detail update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var updates db.CarDetail
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for car detail update")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if updates.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the query param ID",
			"code":    13,
		})
		return
	}

	Image1Enc, err := functions.DecodeBase64ToByteArray(updates.Image1)
	if err != nil {
		log.Err(err).Msg("Error converting image 1")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 1",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	Image2Enc, err := functions.DecodeBase64ToByteArray(updates.Image2)
	if err != nil {
		log.Err(err).Msg("Error converting image 2")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 2",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	log.Info().Int("image1_size", len(Image1Enc)).Int("image2_size", len(Image2Enc)).Msg("Images decoded successfully")

	ctx := context.Background()
	rowsAffected, err := db.UpdateCarDetail(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating car detail by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update car detail",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Int64("rowsAffected", rowsAffected).Msg("No car detail found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No car detail found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Int64("rowsAffected", rowsAffected).Msg("Car detail modified successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":       "Car detail modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// DeleteCarDetail godoc
//
//	@Summary		Delete a car detail
//	@Description	Delete a car detail by ID
//	@Tags			Car Details
//	@Param			id	query		int						true	"Car detail ID"
//	@Success		200	{object}	map[string]interface{}	"Car detail deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid request"
//	@Failure		404	{object}	map[string]interface{}	"Car detail not found"
//	@Router			/fyc/carDetails [delete]
func DeleteCarDetailAPI(c *gin.Context) {

	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Msg("Error ID Format")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := db.DeleteCarDetail(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting car detail")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete car detail",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Int64("rowsAffected", rowsAffected).Msg("No car detail found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No car detail found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Int64("rowsAffected", rowsAffected).Msg("Car detail deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success":      "Car detail deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
