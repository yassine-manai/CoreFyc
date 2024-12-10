package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/functions"
	"fyc/pkg/db"
)

// GetCamera godoc
//
//	@Summary	Get camera data
//	@Tags		Cameras
//	@Param		state	query		string					false	"State of cameras: enabled or deleted"
//	@Success	200		{object}	db.Camera				"List of cameras or a single camera"
//	@Failure	500		{object}	map[string]interface{}	"Internal server error"
//	@Failure	404		{object}	map[string]interface{}	"No cameras found"
//	@Failure	400		{object}	map[string]interface{}	"Bad request: Invalid camera ID"
//	@Param		extra	query		string					false	"Include extra information if 'yes'"
//	@Router		/fyc/cameras [get]
func GetCameraAPI(c *gin.Context) {
	log.Debug().Msg("Get Camera API request")
	ctx := context.Background()

	// Extract query parameters
	idStr := c.Query("id")
	extraReq := c.Query("extra")
	stateReq := c.Query("state")

	if stateReq != "" && stateReq != "enabled" && stateReq != "deleted" {
		log.Error().Str("state", stateReq).Msg("Invalid state parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid state parameter. Must be 'enabled' or 'deleted'",
			"code":    400,
		})
		return
	}

	handleError := func(err error, message string, code int, status int) {
		log.Err(err).Msg(message)
		c.JSON(status, gin.H{
			"error":   http.StatusText(status),
			"message": message,
			"code":    code,
		})
	}

	// Handle retrieving camera by ID
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			handleError(err, "Invalid camera ID format", 12, http.StatusBadRequest)
			return
		}
		// Fetch camera without state filtering
		if extraReq == "yes" {
			camera, err := db.GetCameraByIDExtra(ctx, id)
			if err != nil {
				handleError(err, "Camera with extra data not found", 9, http.StatusNotFound)
				return
			}
			log.Info().Int("Camera ID", id).Msg("Camera with extra data fetched successfully")
			c.JSON(http.StatusOK, camera)
		} else {
			camera, err := db.GetCameraByID(ctx, id)
			if err != nil {
				handleError(err, "Camera by ID not found", 9, http.StatusNotFound)
				return
			}
			log.Info().Int("Camera ID", id).Msg("Camera fetched successfully")
			c.JSON(http.StatusOK, camera)
		}

		return
	}

	if stateReq == "enabled" {
		if extraReq == "yes" {
			cameras, err := db.GetCameraListEnabledExtra(ctx)
			if err != nil {
				handleError(err, "Error retrieving Enabled Cameras with extra data", 9, http.StatusNotFound)
				return
			}
			c.JSON(http.StatusOK, cameras)
		} else {
			cameras, err := db.GetCameraListEnabled(ctx)
			if err != nil {
				handleError(err, "Error retrieving Enabled Cameras", 9, http.StatusNotFound)
				return
			}
			c.JSON(http.StatusOK, cameras)
		}
		return
	}

	// Fetch all cameras without filtering
	if extraReq == "yes" {
		cameras, err := db.GetAllCameraExtra(ctx)
		if err != nil {
			handleError(err, "Error retrieving all cameras with extra data", 9, http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, cameras)
	} else {
		cameras, err := db.GetAllCamera(ctx)
		if err != nil {
			handleError(err, "Error retrieving all cameras", 9, http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, cameras)
	}
}

// CreateCamera godoc
//
//	@Summary		Add a new camera
//	@Description	Add a new camera to the database
//	@Tags			Cameras
//	@Accept			json
//	@Produce		json
//	@Param			Camera	body		db.Camera				true	"Camera data"
//	@Success		201		{object}	db.Camera				"Camera created successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload"
//	@Failure		500		{object}	map[string]interface{}	"Failed to create a new camera"
//	@Router			/fyc/cameras [post]
func CreateCameraAPI(c *gin.Context) {
	log.Debug().Msg("Create Camera API request")
	ctx := context.Background()
	var newCam db.Camera

	if err := c.ShouldBindJSON(&newCam); err != nil {
		log.Err(err).Msg("Invalid input for new camera")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	log.Info().Msg("Creating new camera")
	log.Debug().Msg("Checkin zones")

	if !functions.Contains(db.Zonelist, *newCam.ZoneIdIn) {
		*newCam.ZoneIdIn = 0
	}

	if !functions.Contains(db.Zonelist, *newCam.ZoneIdOut) {
		*newCam.ZoneIdOut = 0
	}

	if err := db.CreateCamera(ctx, &newCam); err != nil {
		log.Err(err).Msg("Error creating new camera")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new camera",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	db.LoadCameralist()
	db.CamStartup()
	log.Info().Int("camera_id", newCam.ID).Msg("Camera created successfully")
	c.JSON(http.StatusCreated, newCam)
}

// UpdateCamera godoc
//
//	@Summary		Update a camera by ID
//	@Description	Update an existing camera by ID
//	@Tags			Cameras
//	@Accept			json
//	@Produce		json
//	@Param			id		query		int						true	"Camera ID"
//	@Param			Camera	body		db.Camera				true	"Updated camera data"
//	@Success		200		{object}	map[string]interface{}	"Camera updated successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload or ID mismatch"
//	@Failure		404		{object}	map[string]interface{}	"Camera not found"
//	@Failure		500		{object}	map[string]interface{}	"Failed to update camera"
//	@Router			/fyc/cameras [put]
func UpdateCameraAPI(c *gin.Context) {
	idStr := c.Query("id")
	var updates db.CameraNoBind
	ctx := context.Background()
	log.Info().Str("camera_id", idStr).Msg("Updating camera in progress")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Str("camera_id", idStr).Msg("Invalid ID format for camera update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for camera update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if updates.CamID != id {
		log.Warn().Msg("The ID in the request body does not match the query ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the query ID",
			"code":    13,
		})
		return
	}

	rowsAffected, err := db.UpdateCamera(ctx, id, &updates)
	if err != nil {
		log.Err(err).Int("Camera ID :", id).Msgf("Error updating camera data : \n Body : %v", &updates)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update camera",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("camera_id", idStr).Int64("Rows Affected :", rowsAffected).Msg("No Camera found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No camera found with the specified ID",
			"code":    9,
		})
		return
	}

	db.LoadCameralist()
	db.CamStartup()
	log.Info().Str("camera_id", idStr).Msg("Camera updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":       "Camera updated successfully",
		"rows_affected": rowsAffected,
		"response":      updates,
		"code":          8,
	})
}

// DeleteCameraAPI godoc
//
//	@Summary		Delete a camera
//	@Description	Delete a camera by setting the is_deleted flag to true
//	@Tags			Cameras
//	@Param			id	query		string					true	"Camera ID"
//	@Success		200	{object}	map[string]interface{}	"Camera deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid camera ID"
//	@Failure		500	{object}	map[string]interface{}	"Failed to delete camera"
//	@Router			/fyc/cameras [delete]
func DeleteCameraAPI(c *gin.Context) {
	idStr := c.Query("id")
	ctx := context.Background()

	if idStr == "" {
		log.Error().Msg("No camera ID provided for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "Camera ID must be provided",
			"code":    12,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid camera ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	log.Info().Int("camera_id", id).Msg("Attempting to soft delete camera")

	rowsAffected, err := db.DeleteCamera(ctx, id)
	if err != nil {
		log.Err(err).Int("camera_id", id).Msg("Failed to soft delete camera")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete camera",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("camera_id", idStr).Int64("Rows Affected :", rowsAffected).Msg("No camera found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No camera found with the specified ID",
			"code":    9,
		})
		return
	}

	db.LoadCameralist()
	db.CamStartup()
	log.Info().Str("camera_id", idStr).Msg("Camera deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":       "Camera deleted successfully",
		"rows_affected": rowsAffected,
		"code":          8,
	})
}

/* // ChangeStateAPI godoc
//
//	@Summary		Change camera state or retrieve cameras by ID
//	@Description	Change the state of a camera (e.g., enabled/deleted) or retrieve a camera by ID
//	@Tags			Cameras
//	@Produce		json
//	@Param			state	query		bool	false	"Camera State"
//	@Param			id		query		int 	false	"Camera ID"
//	@Success		200		{object}	int64		"Number of rows affected by the state change"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No cameras found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid camera ID or state"
//	@Router			/fyc/cameraState [put]
func ChangeCameraStateAPI(c *gin.Context) {
	log.Debug().Msg("ChangeStateAPI request")
	ctx := context.Background()

	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid camera ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	stateStr := c.Query("state")
	state, err := strconv.ParseBool(stateStr)
	if err != nil {
		log.Err(err).Str("state", stateStr).Msg("Invalid state format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid state format",
			"message": "State must be a boolean value (true/false)",
			"code":    13,
		})
		return
	}

	rowsAffected, err := ChangeCameraState(ctx, id, state)
	if err != nil {
		if err.Error() == fmt.Sprintf("camera with id %d is already enabled", id) {
			log.Info().Str("camera_id", idStr).Msg("Camera is already enabled")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Conflict",
				"message": err.Error(),
				"code":    12,
			})
			return
		}

		log.Err(err).Str("camera_id", idStr).Msg("Error changing camera state")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error changing camera state",
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("camera_id", idStr).Msg("Camera not found or state unchanged")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Camera not found or state unchanged",
			"code":    9,
		})
		return
	}

	LoadCameralist()
	log.Info().Str("camera_id", idStr).Bool("state", state).Msg("Camera state changed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":      "Camera state changed successfully",
		"rowsAffected": rowsAffected,
	})
} */
