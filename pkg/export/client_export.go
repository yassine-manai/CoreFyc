package export

import (
	"context"
	"fmt"
	"fyc/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary		Export Client Data
// @Description	Export the Client data in PDF or Excel format based on the `file_type` query parameter
// @Tags			BackOffice - Export
// @Accept			json
// @Produce		application/pdf, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuthBackOffice
// @Param			file_type	query		string		false	"The type of the export file"	Enums(pdf, excel)	default(pdf)
// @Param			client_ids	body		[]string	true	"The list of client IDs"
// @Success		200			{string}	string		"Export successful"
// @Failure		500			{string}	string		"Internal Server Error"
// @Router			/backoffice/export_client [post]
func ExportClient(c *gin.Context) {
	log.Debug().Msg(" / / / / # Exporting Client Data # / / / / ")
	ctx := context.Background()

	fileType := c.DefaultQuery("file_type", "pdf")
	var (
		clientIDs []string
		clients   []db.ApiKeyResponse
		err       error
	)

	if err := c.ShouldBindJSON(&clientIDs); err != nil {
		log.Warn().Msg("Error binding request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"code":    -1,
		})
		return
	}

	if len(clientIDs) > 0 {
		for _, clientID := range clientIDs {
			client, err := db.GetClientById(ctx, clientID)
			if err != nil {
				log.Warn().Str("Client ID", clientID).Msg("Client not found")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": fmt.Sprintf("Client %v not found", clientID),
					"code":    -4,
				})
				return
			}
			clients = append(clients, *client)
		}

	} else {
		clients, err = db.GetClientByStatus(ctx, "")
		if err != nil {
			if err.Error() == "no rows found" {
				log.Warn().Msg("No clients found")
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "No clients found",
					"code":    -4,
				})
				return
			}
			log.Debug().Msgf("Error fetching clients: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching client data"})
			return
		}
	}

	headers := []string{"Client ID", "Client Name", "Grant Type", "Fuzzy Logic", "Last Updated", "Status"}
	widths := []float64{40, 60, 40, 40, 60, 35}
	var data [][]string

	for _, client := range clients {
		status := "Disabled"
		if client.IsEnabled {
			status = "Enabled"
		}
		data = append(data, []string{
			client.ClientID,
			client.ClientName,
			client.GrantType,
			fmt.Sprintf("%v", *client.FuzzyLogic),
			client.LastUpdated,
			status,
		})
	}
	switch fileType {
	case "excel":
		ExportToExcel(c, data, headers, "Client_Data")
	default:
		ExportToPDF(c, "L", data, headers, widths, "Clients Export", "Client_Data", "./font/Cairo-Regular.ttf")
	}
}
