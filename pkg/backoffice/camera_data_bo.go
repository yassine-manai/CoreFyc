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

// GetCameras godoc
//
//	@Summary		Get cameras
//	@Description	Get a list of all cameras or a specific camera
//	@Tags			Backoffice - Camera
//	@Success		200	{object}	db.Camera	"List of cameras or a single camera"
//	@Security		BearerAuthBackOffice
//	@Param			camera_id	query	string	false	"Camera ID"
//	@Param			extra		query	string	false	"Include extra information if 'true'"
//	@Router			/backoffice/getCameras [get]
func GetCameraDataAPI(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Query("camera_id")
	extraData := c.DefaultQuery("extra", "false")
	handler := functions.NewResponseHandler()
	var (
		cameras  []db.Camera
		response []map[string]interface{}
		err      error
	)

	log.Debug().Str("Extra Param", extraData).Str("Camera Id", idStr).Msg("Get Camera API request")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("Invalid camera ID format")
			handler.RespondWithStatus(c, -5)
			return
		}

		camera, err := db.GetCameraByIDExtra(ctx, id)
		if err != nil {
			log.Error().Err(err).Int("cameraID", id).Msg("Camera data not found")
			c.JSON(http.StatusOK, []map[string]interface{}{})
			return
		}
		cameras = append(cameras, *camera)
	} else {
		cameras, err = db.GetAllCameraExtra(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Camera data not found")
			c.JSON(http.StatusOK, []map[string]interface{}{})
			return
		}
	}

	for _, camera := range cameras {
		zoneIn, err := fetchZoneData(ctx, camera.ZoneIdIn, "Zone IN")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error(), "code": 404})
			return
		}

		zoneOut, err := fetchZoneData(ctx, camera.ZoneIdOut, "Zone OUT")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error(), "code": 404})
			return
		}

		responseData := buildCameraResponse(camera, zoneIn, zoneOut, extraData)
		response = append(response, responseData)
	}

	if len(response) == 0 {
		c.JSON(http.StatusOK, []map[string]interface{}{})
	} else if len(response) == 1 {
		log.Info().Msg("One data found")
		c.JSON(http.StatusOK, response[0])
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// Helper function to fetch zone data
func fetchZoneData(ctx context.Context, zoneID *int, zoneType string) (*db.Zone, error) {
	if zoneID == nil {
		return &db.Zone{ZoneID: 0, Name: make(map[string]interface{})}, nil
	}

	zone, err := db.GetZoneByID(ctx, *zoneID)
	if err != nil || zone == nil {
		log.Error().Err(err).Int("zoneID", *zoneID).Msg(fmt.Sprintf("%s not found", zoneType))
		return nil, fmt.Errorf("%s %v data not found", zoneType, *zoneID)
	}

	if !zone.IsEnabled || zone.IsDeleted {
		zone.ZoneID = 0
		zone.Name = make(map[string]interface{})
	}
	return zone, nil
}

// Helper function to build camera response data
func buildCameraResponse(camera db.Camera, zoneIn *db.Zone, zoneOut *db.Zone, extraData string) map[string]interface{} {
	responseData := map[string]interface{}{
		"cam_id":        camera.CamID,
		"cam_name":      camera.CamName,
		"cam_type":      camera.CamType,
		"cam_ip":        camera.CamIP,
		"cam_port":      camera.CamPORT,
		"cam_user":      camera.CamUser,
		"cam_password":  camera.CamPass,
		"zone_in_id":    zoneIn.ZoneID,
		"zone_in_name":  zoneIn.Name,
		"zone_out_id":   zoneOut.ZoneID,
		"zone_out_name": zoneOut.Name,
		"direction":     camera.Direction,
		"is_enabled":    camera.IsEnabled,
		"last_update":   camera.LastUpdated,
	}

	if extraData == "true" {
		responseData["extra"] = camera.Extra
	}

	return responseData
}

// CreateCamera godoc
//
//	@Summary		Add a new camera
//	@Description	Add a new camera to the database
//	@Tags			Backoffice - Camera
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			Camera	body		db.Camera	true	"Camera data"
//	@Success		201		{object}	db.Camera	"Camera created successfully"
//	@Router			/backoffice/addCamera [post]
func AddCameraDataAPI(c *gin.Context) {
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

	////////////////////////////////	log.Info().Msg(" -- -- -- -- -- Creating new camera -- -- -- -- --")
	log.Info().Msg(" -- -- -- -- -- Checking zones -- -- -- -- --")

	if !functions.Contains(db.Zonelist, *newCam.ZoneIdIn) {
		log.Warn().Int("Zone IN ID", *newCam.ZoneIdIn).Msg("Zone doesn't exists")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    -8,
			"message": fmt.Sprintf("Zone IN %v not found !", *newCam.ZoneIdIn),
		})
		return
	}

	if !functions.Contains(db.Zonelist, *newCam.ZoneIdOut) {
		log.Warn().Int("Zone OUT ID", *newCam.ZoneIdOut).Msg("Zone doesn't exists")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    -8,
			"message": fmt.Sprintf("Zone OUT %v not found !", *newCam.ZoneIdOut),
		})
		return
	}

	if functions.Contains(db.CameraList, newCam.CamID) {
		log.Warn().Int("Camera ID", newCam.CamID).Msg("Camera Already Exist !")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    -8,
			"message": fmt.Sprintf("Camera %v Already Exist !", newCam.CamID),
		})
		return
	}

	if err := db.CreateCamera(ctx, &newCam); err != nil {
		log.Err(err).Msg("Error creating new camera")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    -500,
			"message": "An unexpected error occurred. Please try again later.",
			"error":   err,
		})
		return
	}

	log.Info().Int("camera_id", newCam.CamID).Msg("Camera created successfully")
	c.JSON(204, gin.H{
		"success": true,
		"message": "Camera Added Successfully",
	})
}

// UpdateCamera godoc
//
//	@Summary		Update a camera by ID
//	@Description	Update an existing camera by ID
//	@Tags			Backoffice - Camera
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			camera_id	query	int			true	"Camera ID"
//	@Param			Camera		body	db.Camera	true	"Updated camera data"
//	@Router			/backoffice/updateCamera [put]
func UpdateCameraDataAPI(c *gin.Context) {
	idStr := c.Query("camera_id")
	var updates db.CameraNoBind
	ctx := context.Background()

	log.Info().Msg(" -- -- -- -- -- Updating new camera -- -- -- -- --")
	if idStr == "" {
		log.Warn().Msg("The Camera ID is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "The Camera ID is required",
			"code":    -5,
		})
		return
	}

	log.Info().Str("camera_id", idStr).Msg("Updating camera in progress")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Str("camera_id", idStr).Msg("Invalid ID format for camera update")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID must be a valid integer",
			"code":    -5,
		})
		return
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for camera update")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -5,
		})
		return
	}

	if !functions.Contains(db.CameraList, id) {
		log.Warn().Int("Camera ID", id).Msg("Camera doesn't exists")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    -8,
			"message": fmt.Sprintf("Camera %v not found !", id),
		})
		return
	}

	rowsAffected, err := db.UpdateCamera(ctx, id, &updates)
	if err != nil {
		log.Err(err).Int("Camera ID :", id).Msgf("Error updating camera data : \n Body : %v", &updates)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -5,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("camera_id", idStr).Int64("Rows Affected :", rowsAffected).Msg("No Camera found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No camera found with the specified ID",
			"code":    -9,
		})
		return
	}

	log.Info().Str("camera_id", idStr).Msg("Camera updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Camera Updated successfully",
	})
}

// DeleteCameraAPI godoc
//
//	@Summary		Delete a camera
//	@Description	Delete a camera by setting the is_deleted flag to true
//	@Tags			Backoffice - Camera
//	@Security		BearerAuthBackOffice
//	@Param			camera_id	query	string	true	"Camera ID"
//	@Router			/backoffice/deleteCameras [delete]
func DeleteCameraDataAPI(c *gin.Context) {
	idStr := c.Query("camera_id")
	ctx := context.Background()

	if idStr == "" {
		log.Error().Str("Id provided", idStr).Msg("No camera ID provided for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request. 'Camera id' parameter is required.",
			"code":    -5,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid camera ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID must be a valid integer",
			"code":    -5,
		})
		return
	}

	log.Info().Int("camera_id", id).Msg("Attempting to delete camera")

	rowsAffected, err := db.DeleteCamera(ctx, id)
	if err != nil {
		log.Err(err).Int("camera_id", id).Msg("Failed to soft delete camera")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("camera_id", idStr).Int64("Rows Affected :", rowsAffected).Msg("No camera found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Camera ID not found !",
			"code":    -9,
		})
		return
	}

	log.Info().Str("camera_id", idStr).Msg("Camera deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Camera deleted successfully",
	})
}
