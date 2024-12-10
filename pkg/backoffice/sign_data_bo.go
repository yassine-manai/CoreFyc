package backoffice

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/functions"
	"fyc/pkg/db"
)

// GetSignAPI godoc
//
//	@Summary		Get signs or a specific sign by ID
//	@Description	Get a list of signs or a specific sign by ID with optional extra data
//	@Tags			Backoffice - Signs
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			sign_id	query		int							false	"sign ID"
//	@Success		200		{object}	[]map[string]interface{}	"List of signs or a single sign"
//	@Router			/backoffice/getSign [get]
func GetSignDataAPI(c *gin.Context) {
	log.Debug().Msg("Starting GetSignAPI request")
	ctx := context.Background()
	idStr := c.Query("sign_id")
	var response []map[string]interface{}
	//var responseSingle map[string]interface{}

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		log.Info().Str("sign_id", idStr).Msg("Received request to fetch a specific sign by ID")

		if err != nil {
			log.Error().Err(err).Str("id", idStr).Msg("Invalid sign ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		sign, err := db.GetSignById(ctx, id)
		if err != nil || sign == nil {
			log.Error().Err(err).Str("sign_id", idStr).Msg("Error retrieving sign by ID or sign not found")
			c.JSON(http.StatusOK, []db.Sign{})
			return
		}

		log.Info().Int("sign_id", sign.SignID).Msg("Successfully retrieved sign data")

		zone, err := db.GetZoneByID(ctx, sign.ZoneID)
		if err != nil || zone == nil {
			log.Error().Err(err).Int("zoneID", sign.ZoneID).Msg("Zone not found for the sign")
			c.JSON(http.StatusOK, []db.Zone{})
		}

		log.Info().Int("zoneID", zone.ZoneID).Msg("Successfully retrieved zone data for the sign")

		responseData := map[string]interface{}{
			"sign_id":       sign.SignID,
			"sign_name":     sign.SignName,
			"sign_ip":       sign.SignIP,
			"sign_port":     sign.SignPort,
			"sign_type":     sign.SignType,
			"sign_username": sign.SignUserName,
			"sign_password": sign.SignPassword,
			"zone_id":       zone.ZoneID,
			"zone_name":     zone.Name,
			"is_enabled":    sign.IsEnabled,
			"last_update":   sign.LastUpdated,
		}

		if len(responseData) == 0 {
			c.JSON(http.StatusOK, response)
			return
		}

		log.Info().Str("sign_id", idStr).Msg("Returning specific sign data in response")
		//response = append(response, responseData)
		c.JSON(http.StatusOK, responseData)
		return
	}

	// If no ID is provided, fetch all signs
	log.Info().Msg("No sign ID provided, fetching all signs")
	signs, err := db.GetAllSigns(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving all signs from the database")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve signs",
			"code":    -500,
		})
		return
	}

	if len(signs) == 0 {
		log.Info().Msg("No signs found in the database")
		c.JSON(http.StatusOK, []db.Zone{})
		return
	}

	log.Info().Int("sign_count", len(signs)).Msg("Successfully retrieved all signs")

	for _, sign := range signs {
		zone, err := db.GetZoneByID(ctx, sign.ZoneID)
		if err != nil || zone == nil {
			log.Err(err).Int("zoneID", sign.ZoneID).Msg("Zone not found for a sign")
			continue
		}

		responseData := map[string]interface{}{
			"sign_id":       sign.SignID,
			"sign_name":     sign.SignName,
			"sign_ip":       sign.SignIP,
			"sign_port":     sign.SignPort,
			"sign_type":     sign.SignType,
			"sign_username": sign.SignUserName,
			"sign_password": sign.SignPassword,
			"zone_id":       zone.ZoneID,
			"zone_name":     zone.Name,
			"is_enabled":    sign.IsEnabled,
			"last_update":   sign.LastUpdated,
		}

		//log.Info().Int("sign_id", sign.SignID).Msg("Appending sign data to response")
		response = append(response, responseData)
	}

	if len(response) == 0 {
		c.JSON(http.StatusOK, response)
		return
	}

	log.Info().Msg("Returning all sign data in response")
	c.JSON(http.StatusOK, response)
}

// Createsign godoc
//
//	@Summary		Add a new sign
//	@Description	Add a new sign to the database
//	@Tags			Backoffice - Signs
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			sign	body		models.AddSignModel	true	"sign data"
//	@Success		201		{object}	models.AddSignModel	"sign created successfully"
//	@Router			/backoffice/addSign [post]
func CreateSignDataAPI(c *gin.Context) {

	log.Info().Msg(" - - - - - # Creating new sign # - - - - - ")
	ctx := context.Background()
	var newSign db.Sign

	if err := c.ShouldBindJSON(&newSign); err != nil {
		log.Err(err).Msg("Invalid input for new sign")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    -5,
		})
		return
	}

	if functions.Contains(db.SignList, newSign.SignID) {
		log.Warn().Int("SignID", newSign.SignID).Msg("Sign ID alreay exist !")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Sign ID %v alreay exist !", newSign.SignID),
			"code":    -10,
		})
		return
	}

	if !functions.Contains(db.Zonelist, newSign.ZoneID) {
		log.Warn().Int("ZoneID", newSign.ZoneID).Msg("Zone Not Found !")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Zone ID %v Not Found !", newSign.ZoneID),
			"code":    -10,
		})
		return
	}

	if err := db.CreateSign(ctx, &newSign); err != nil {
		log.Err(err).Msg("Error creating new sign")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	log.Info().Int("sign_id", newSign.ID).Msg("Sign created successfully")
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Sign Added Successfully",
	})
}

// Updatesign godoc
//
//	@Summary		Update a sign by ID
//	@Description	Update an existing sign by ID
//	@Tags			Backoffice - Signs
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			sign_id	query	int						true	"sign ID"
//	@Param			sign	body	models.UpdateSignModel	true	"Updated sign data"
//	@Router			/backoffice/updateSign [put]
func UpdateSignDataAPI(c *gin.Context) {

	sign_id := c.Query("sign_id")
	var updates db.SignNoBind
	ctx := context.Background()
	log.Info().Str("sign_id", sign_id).Msg("Updating sign")

	if sign_id == "" {
		log.Warn().Msg("The Sign ID is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "The Sing ID is required",
			"code":    -5,
		})
		return
	}

	id, err := strconv.Atoi(sign_id)
	if err != nil {
		log.Error().Str("sign_id", sign_id).Msg("Invalid ID format for sign update")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for sign update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if !functions.Contains(db.SignList, id) {
		log.Warn().Int("sign_id", id).Msg("Sign ID not exist !")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("Sign ID %v not exist !", id),
			"code":    -10,
		})
		return
	}

	rowsAffected, err := db.UpdateSign(ctx, id, updates)
	if err != nil {
		log.Err(err).Msg("Error updating sign")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update sign",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("sign_id", sign_id).Int64("Rows Affected", rowsAffected).Msg("No sign found with the  ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No sign found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("sign_id", sign_id).Int("Rows Affected ", int(rowsAffected)).Msg("sign updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sign updated successfully",
	})
}

// DeletesignAPI godoc
//
//	@Summary		Soft delete a sign
//	@Description	Soft delete a sign by setting the is_deleted flag to true
//	@Tags			Backoffice - Signs
//	@Security		BearerAuthBackOffice
//	@Param			sign_id	query	int	true	"sign ID"
//	@Router			/backoffice/deleteSign [delete]
func DeleteSignDataAPI(c *gin.Context) {
	idStr := c.Query("sign_id")
	ctx := context.Background()

	if idStr == "" {
		log.Error().Msg("No sign ID provided for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "sign ID must be provided",
			"code":    12,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid sign ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	if !functions.Contains(db.SignList, id) {
		log.Warn().Int("sign_id", id).Msg("Sign ID not exist !")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("Sign ID %v not exist !", id),
			"code":    -10,
		})
		return
	}

	log.Info().Int("sign_id", id).Msg("Attempting to soft delete sign")

	rowsAffected, err := db.DeleteSign(ctx, id)
	if err != nil {
		log.Err(err).Int("sign_id", id).Msg("Failed to soft delete sign")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No Sign found with the specified ID ------  affected rows 0 ",
			"code":    9,
		})
		return
	}
	log.Info().Int("sign_id", id).Int("Rows Affected ", int(rowsAffected)).Msg("sign deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success":      "true",
		"message":      "Sign deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
