package backoffice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/functions"
	"fyc/pkg/db"
)

// GetClients godoc
//
//	@Summary		Get all Clients
//	@Description	Get a list of all Clients
//	@Tags			Backoffice - Clients
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			clientId	query	string		false	"ClientID"
//	@Success		200			{array}	db.ApiKey	"List of Clients"
//	@Router			/backoffice/get_clients [get]
func GetClients(c *gin.Context) {
	ctx := context.Background()
	idStr := c.Query("clientId")
	//var clientRes []*db.ApiKeyResponse

	if idStr != "" {

		log.Debug().Str("ClientID", idStr).Msg("Get Client by ID API request")
		client, err := db.GetClientCredByIdFalse(ctx, idStr)
		if err != nil {
			log.Warn().Err(err).Str("Client_ID", idStr).Msg("Error retrieving Client ID")
			c.JSON(http.StatusOK, []db.ApiKeyResponse{})
			return
		}

		//clientRes = append(clientRes, client)

		/* if len(clientRes) == 0 {
			log.Debug().Int("Clients", len(clientRes)).Msg("Not datat found ")
			c.JSON(http.StatusOK, []db.ApiKeyResponse{})
			return
		} */

		c.JSON(http.StatusOK, client)
		return
	}

	clients, err := db.GetAllClientCred(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all clients")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if len(clients) == 0 {
		log.Debug().Int("Client List", len(clients)).Msg("Not datat found ")
		c.JSON(http.StatusOK, []db.ApiKeyResponse{})
		return
	}

	c.JSON(http.StatusOK, clients)
}

// AddClient godoc
//
//	@Summary		Add a new client credential
//	@Description	Add a new client credential to the database
//	@Tags			Backoffice - Clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			clientCred	body		db.ApiKey	true	"Client credential data"
//	@Success		201			{object}	db.ApiKey
//	@Router			/backoffice/addClient [post]
func AddClientAPI(c *gin.Context) {

	var clientCred db.ApiKey
	log.Info().Msg("Adding new Client API KEY")

	if err := c.ShouldBindJSON(&clientCred); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for Client API KEY creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -5,
		})
		return
	}

	if clientCred.ClientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Client ID is Required !",
			"code":    -5,
		})
		return
	}

	if functions.ContainsStr(db.ClientListAPI, clientCred.ClientID) {
		log.Debug().Str("Client exist with ID ", clientCred.ClientID).Msg("Client exist with ID")

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Client ID %v already exist !", clientCred.ClientID),
			"code":    -12,
		})
		return
	}

	ctx := context.Background()
	if err := db.AddClientCred(ctx, &clientCred); err != nil {
		log.Error().Err(err).Msg("Error creating Client API KEY")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	log.Info().Str("client_id", clientCred.ClientID).Msg("Client created successfully")
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Client Added Successfully",
	})
}

// UpdateClient godoc
//
//	@Summary		Update client credential
//	@Description	Update client credential
//	@Tags			Backoffice - Clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			client_id	query		string		true	"Client ID"
//	@Param			Client		body		db.ApiKey	true	"Client credential data"
//	@Success		201			{object}	db.ApiKey
//	@Router			/backoffice/updateClient [put]
func UpdateClientAPI(c *gin.Context) {

	idStr := c.Query("client_id")
	log.Info().Str("client_id", idStr).Msg("Attempting to update Client API KEY")

	var clientCred db.ApiKeyNoBind
	log.Info().Msg("Updating Client API KEY")

	if err := c.ShouldBindJSON(&clientCred); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for Client API KEY creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -5,
		})
		return
	}

	if idStr == "" {
		log.Warn().Msg("The Client ID is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Client ID is Required !",
			"code":    -5,
		})
		return
	}

	if !functions.ContainsStr(db.ClientListAPI, idStr) {
		log.Debug().Str("ID ", idStr).Msg("Client not exist with ID")

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Client ID %v doesn't exist !", idStr),
			"code":    -12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := db.UpdateClientCred(ctx, idStr, &clientCred)
	if err != nil {
		log.Error().Err(err).Msg("Error updating Client API KEY")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("client_id", clientCred.ClientID).Msg("No Client API found to delete")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No Client API found with the specified ID",
			"code":    -4,
		})
		return
	}

	log.Info().Str("client_id", clientCred.ClientID).Msg("Client updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Client Updated successfully",
	})
}

// DeleteClient godoc
//
//	@Summary		Delete a client credential
//	@Description	Delete a client credential by ID
//	@Tags			Backoffice - Clients
//	@Security		BearerAuthBackOffice
//	@Param			client_id	query		string	true	"Client ID"
//	@Success		200			{string}	string	"Client API deleted successfully"
//	@Router			/backoffice/deleteClient [delete]
func DeleteClientAPI(c *gin.Context) {

	id := c.Query("client_id")
	log.Info().Str("client_id", id).Msg("Attempting to delete Client API")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request. ClientID parameter is required.",
			"code":    -5,
		})
		return
	}
	ctx := context.Background()
	rowsAffected, err := db.DeleteClientCred(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("client_id", id).Msg("Error deleting Client API KEY")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    -500,
			"message": "An unexpected error occurred. Please try again later.",
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("client_id", id).Msg("No Client API  found to delete")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No Client API found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("client_id", id).Msg("Client API  deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Client deleted successfully",
	})
}
