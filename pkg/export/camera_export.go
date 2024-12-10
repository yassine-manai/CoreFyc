package export

import (
	"context"
	"fmt"
	"fyc/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary		Export Camera Data
// @Description	Export the Camera data in PDF or Excel format based on the `file_type` query parameter
// @Tags			BackOffice - Export
// @Accept			json
// @Produce		application/pdf, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuthBackOffice
// @Param			file_type	query		string	false	"The type of the export file"	Enums(pdf, excel)	default(pdf)
// @Param			camera_ids	body		[]int	true	"The list of camera IDs"
// @Success		200			{string}	string	"Export successful"
// @Failure		500			{string}	string	"Internal Server Error"
// @Router			/backoffice/export_camera [post]
func ExportCamera(c *gin.Context) {
	log.Debug().Msg(" / / / / # Exporting Camera Data # / / / /  ")
	ctx := context.Background()

	// Query params
	fileType := c.DefaultQuery("file_type", "pdf")

	var (
		cameraIDs []int
		cameras   []db.ResponseCamera
		err       error
	)

	// Fetch camera IDs from request body
	if err := c.ShouldBindJSON(&cameraIDs); err != nil {
		log.Warn().Msg("Error binding request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"code":    -1,
		})
		return
	}

	log.Debug().Interface("cameraIDs Passed in param", cameraIDs).Send()

	// Fetch cameras based on the camera IDs (if provided)
	if len(cameraIDs) > 0 {
		for _, camID := range cameraIDs {
			// Fetch camera by ID
			camera, err := db.GetCamByID(ctx, camID)
			if err != nil {
				log.Warn().Int("Camera ID", camID).Msg("Camera not found")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": fmt.Sprintf("Camera %v not found", camID),
					"code":    -4,
				})
				return
			}
			cameras = append(cameras, *camera)
		}
	} else {
		cameras, err = db.GetCameraByStatus(ctx, "")
		if err != nil {
			if err.Error() == "no rows found" {
				log.Warn().Str("Status", "All data").Msg("No cameras found, returning all the data")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "No cameras found",
					"code":    -4,
				})
				return
			}
			log.Debug().Msgf("Error fetching data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching camera data"})
			return
		}
	}

	headers := []string{"ID", "Camera Name", "Type", "Camera IP", "Port", "User", "Zone In", "Zone Out", "Direction", "Status", "Last Updated"}
	widths := []float64{15, 30, 25, 30, 20, 20, 20, 20, 30, 30, 40}
	var data [][]string

	for _, camera := range cameras {
		status := "Disabled"
		if camera.IsEnabled {
			status = "Enabled"
		}
		row := []string{
			fmt.Sprintf("%d", camera.CamID),
			camera.CamName,
			camera.CamType,
			camera.CamIP,
			fmt.Sprintf("%d", camera.CamPORT),
			camera.CamUser,
			fmt.Sprintf("%v", *camera.ZoneIdIn),
			fmt.Sprintf("%v", *camera.ZoneIdOut),
			camera.Direction,
			status,
			camera.LastUpdated,
		}
		data = append(data, row)
	}

	switch fileType {
	case "excel":
		ExportToExcel(c, data, headers, "Camera_Data")
	default:
		ExportToPDF(c, "L", data, headers, widths, "Clients Data Export", "Camera_Data", "./font/Cairo-Regular.ttf")
	}
}
