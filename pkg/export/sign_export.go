package export

import (
	"context"
	"fmt"
	"fyc/pkg/db"
	"net/http"

	arabic "github.com/abdullahdiaa/garabic"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary		Export Sign Data
// @Description	Export the Sign data in PDF or Excel format based on the `file_type` query parameter
// @Tags			BackOffice - Export
// @Accept			json
// @Produce		application/pdf, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuthBackOffice
// @Param			sign_ids	body		[]int	true	"The list of Sign IDs"
// @Param			file_type	query		string	false	"The type of the export file"	Enums(pdf, excel)	default(pdf)
// @Success		200			{string}	string	"Export successful"
// @Failure		500			{string}	string	"Internal Server Error"
// @Router			/backoffice/export_sign [post]
func ExportSign(c *gin.Context) {
	log.Debug().Msg(" / / / / # Exporting Sign Data # / / / /  ")
	ctx := context.Background()

	// Query params
	fileType := c.DefaultQuery("file_type", "pdf")
	var signIDs []int
	var signs []db.SignResp

	// Fetch sign_ids from body
	if err := c.ShouldBindJSON(&signIDs); err != nil {
		log.Warn().Msg("Error binding request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"code":    -1,
		})
		return
	}

	// Fetch signs based on the sign_ids (if provided)
	if len(signIDs) > 0 {
		for _, signID := range signIDs {
			sign, err := db.GetSignById(ctx, signID)
			if err != nil {
				log.Warn().Int("SignID", signID).Msg("Sign not found")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": fmt.Sprintf("Sign %v not found", signID),
					"code":    -4,
				})
				return
			}
			signs = append(signs, *sign)
		}
	} else {
		var err error
		signs, err = db.GetSignByStatus(ctx, "")
		if err != nil {
			if err.Error() == "no rows found" {
				log.Warn().Msg("No signs found")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "No Signs found",
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
	headers := []string{"Sign ID", "Sign Name EN", "Sign Name AR", "Sign Type", "Sign IP", "Sign Port", "Sign Username", "Zone ID", "Status", "Last Updated"}
	widths := []float64{20, 35, 35, 30, 30, 20, 30, 20, 25, 35}
	data := [][]string{}

	for _, sign := range signs {
		signNameAr, arExists := sign.SignName["ar"]
		if !arExists || signNameAr == "" {
			signNameAr = "No Arabic Sign Name"
		} else {
			signNameAr = arabic.Shape(fmt.Sprintf("%v", signNameAr))
		}

		signNameEn, enExists := sign.SignName["en"]
		if !enExists || signNameEn == "" {
			signNameEn = "No English Sign Name"
		}

		status := "Disabled"
		if sign.IsEnabled {
			status = "Enabled"
		}

		data = append(data, []string{
			fmt.Sprintf("%d", sign.SignID),
			fmt.Sprintf("%v", signNameEn),
			fmt.Sprintf("%v", signNameAr),
			sign.SignType,
			sign.SignIP,
			fmt.Sprintf("%d", sign.SignPort),
			sign.SignUserName,
			fmt.Sprintf("%d", sign.ZoneID),
			status,
			sign.LastUpdated,
		})
	}

	switch fileType {
	case "excel":
		ExportToExcel(c, data, headers, "sign_data")
	default:
		ExportToPDF(c, "L", data, headers, widths, "Signs Data Export", "sign_data", "./font/Cairo-Regular.ttf")
	}

}
