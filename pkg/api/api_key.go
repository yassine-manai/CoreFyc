package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/pkg/db"
)

// GetCred godoc
//
//	@Summary		Get a Client by ClientID, or all Client
//	@Description	Get a Client by ClientID, or all Client
//	@Tags			Client API
//	@Produce		json
//	@Param			clientId	query		string					false	"ClientID"
//	@Success		200			{object}	db.ApiKey				"List of clients or a single client"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Failure		404			{object}	map[string]interface{}	"No Client found"
//	@Failure		400			{object}	map[string]interface{}	"Bad request: Invalid client ID"
//	@Router			/fyc/apikey [get]
func GetAllClientCredsApi(c *gin.Context) {
	log.Debug().Msg("Get Client API request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr == "" {
		api, err := db.GetAllClientCred(ctx)
		if err != nil || len(api) == 0 {
			log.Err(err).Msg("Error retrieving Client API")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No Client API found",
				"code":    9,
			})
			return
		}
		c.JSON(http.StatusOK, api)
		return
	}

	if idStr != "" {
		client, err := db.GetClientCredById(ctx, idStr)
		if err != nil {
			log.Err(err).Str("Client_ID", idStr).Msg("Error retrieving Client ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Client ID not found",
				"code":    9,
			})
			return
		}
		c.JSON(http.StatusOK, client)
		return
	}
}

// AddClientCred godoc
//
//	@Summary		Add a new client credential
//	@Description	Add a new client credential to the database
//	@Tags			Client API
//	@Accept			json
//	@Produce		json
//	@Param			clientCred	body		db.ApiKey	true	"Client credential data"
//	@Success		201			{object}	db.ApiKey
//	@Router			/fyc/apikey [post]
func AddClientCredAPI(c *gin.Context) {
	var clientCred db.ApiKey
	log.Info().Msg("Attempting to add new Client API KEY")

	if err := c.ShouldBindJSON(&clientCred); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for Client API KEY creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	if err := db.AddClientCred(ctx, &clientCred); err != nil {
		log.Error().Err(err).Msg("Error creating Client API KEY")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create Client API KEY",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Str("client_id", clientCred.ClientID).Msg("Client API KEY created successfully")
	c.JSON(http.StatusCreated, clientCred)
}

// UpdateClientCred godoc
//
//	@Summary		Update a client credential
//	@Description	Update an existing client credential by ID
//	@Tags			Client API
//	@Accept			json
//	@Produce		json
//	@Param			client_id	query		string		true	"Client ID"
//	@Param			clientCred	body		db.ApiKey	true	"Updated client credential data"
//	@Success		200			{object}	db.ApiKey
//	@Router			/fyc/apikey [put]
func UpdateClientCredAPI(c *gin.Context) {
	idStr := c.Query("client_id")

	log.Info().Str("client_id", idStr).Msg("Attempting to update Client API KEY")

	var clientCred db.ApiKeyNoBind
	if err := c.ShouldBindJSON(&clientCred); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for Client API KEY update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": "Invalid Client API KEY data",
			"code":    12,
		})
		return
	}

	if clientCred.ClientID != idStr {
		log.Warn().Str("id_param", idStr).Str("id_body", clientCred.ClientID).Msg("ID mismatch between Query and body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the ID in the query parameter",
			"code":    13,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := db.UpdateClientCred(ctx, idStr, &clientCred)
	if err != nil {
		log.Error().Err(err).Str("client_id", idStr).Msg("Error updating Client API KEY")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update Client API KEY",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("client_id", idStr).Msg("No Client API KEY found to update")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No Client API KEY found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("client_id", idStr).Msg("Client API KEY updated successfully")
	c.JSON(http.StatusOK, clientCred)
}

// DeleteClientCred godoc
//
//	@Summary		Delete a client credential
//	@Description	Delete a client credential by ID
//	@Tags			Client API
//	@Param			id	query		string	true	"Client ID"
//	@Success		200	{string}	string	"Client credential deleted successfully"
//	@Router			/fyc/apikey [delete]
func DeleteClientCredAPI(c *gin.Context) {
	idStr := c.Query("id")

	log.Info().Str("client_id", idStr).Msg("Attempting to delete Client API KEY")

	ctx := context.Background()
	rowsAffected, err := db.DeleteClientCred(ctx, idStr)
	if err != nil {
		log.Error().Err(err).Str("client_id", idStr).Msg("Error deleting Client API KEY")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete Client API KEY",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("client_id", idStr).Msg("No Client API KEY found to delete")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No Client API KEY found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("client_id", idStr).Msg("Client API KEY deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": "Client API KEY deleted successfully",
		"code":    8,
	})
}

/*
// GetClientAPI godoc
//
//	@Summary		Get enabled clients or a specific client by clientID
//	@Description	Get a list of enabled clients or a specific client by ID with optional extra data
//	@Tags			Client API
//	@Produce		json
//	@Param			id		query		string	false	"Client ID"
//	@Success		200		{object}	db.ApiKey		"List of enabled clients or a single client"
//	@Router			/fyc/clientEnabled [get]
func GetClientEnabledAPI(c *gin.Context) {
	log.Debug().Msg("Get Enabled API request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		Client, err := db.GetClientEnabledByID(ctx, idStr)
		if err != nil {
			log.Err(err).Str("Client_id", idStr).Msg("Error retrieving Client by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Client not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("Client_id", idStr).Msg("Enabled Client fetched successfully")
		c.JSON(http.StatusOK, Client)
		return
	}

	// Fetch all enabled Clients
	Clients, err := db.GetClientListEnabled(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving enabled Clients")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving enabled Clients",
			"code":    10,
		})
		return
	}

	if len(Clients) == 0 {
		log.Info().Msg("No enabled Clients found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No enabled Clients found",
			"code":    9,
		})
		return
	}

	log.Info().Int("Client_count", len(Clients)).Msg("Enabled Clients fetched successfully")
	c.JSON(http.StatusOK, Clients)
}

// GetClientAPI godoc
//
//	@Summary		Get deleted Clients or a specific Client by ID
//	@Description	Get a list of deleted Client or a specific Client by ID with optional extra data
//	@Tags			Client API
//	@Produce		json
//	@Param			id		query		string	false	"Client ID"
//	@Success		200		{object}	ApiKey		"List of deleted Clients or a Client Client"
//	@Router			/fyc/clientsDeleted [get]
func GetClientDeletedAPI(c *gin.Context) {
	log.Debug().Msg("Get Client Deleted API request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		Client, err := GetClientDeletedByID(ctx, idStr)
		if err != nil {
			log.Err(err).Str("Client_ID", idStr).Msg("Error retrieving Client by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Client not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("Client_ID", idStr).Msg("Deleted Client fetched successfully")
		c.JSON(http.StatusOK, Client)
		return
	}

	// Fetch all deleted Clients
	Clients, err := GetClientListDeleted(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving deleted Clients")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving deleted Clients",
			"code":    10,
		})
		return
	}

	if len(Clients) == 0 {
		log.Info().Msg("No deleted Clients found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No deleted Clients found",
			"code":    9,
		})
		return
	}

	log.Info().Int("Client_count", len(Clients)).Msg("Deleted Clients fetched successfully")
	c.JSON(http.StatusOK, Clients)
} */

/* // ChangeStateAPI godoc
//
//	@Summary		Change Client state or retrieve Client by ID
//	@Description	Change the state of a Client (e.g., enabled/disabled) or retrieve a client by ID
//	@Tags			Client API
//	@Produce		json
//	@Param			state	query		bool	false	"Client State"
//	@Param			id		query		int 	false	"Client ID"
//	@Success		200		{object}	int64		"Number of rows affected by the state change"
//	@Router			/fyc/clientState [put]
func ChangeClientStateAPI(c *gin.Context) {
	log.Debug().Msg("ChangeStateAPI request")
	ctx := context.Background()
	id := c.Query("id")

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

	rowsAffected, err := ChangeApiKeyState(ctx, id, state)
	if err != nil {
		if err.Error() == fmt.Sprintf("client with id %s is already enabled", id) {
			log.Info().Str("client_id", id).Msg("client is already enabled")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Conflict",
				"message": err.Error(),
				"code":    12,
			})
			return
		}

		log.Err(err).Str("client_id", id).Msg("Error changing client state")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "An unexpected error occurred",
			"message": fmt.Sprintf("client with id %s is already enabled", id),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("client_id", id).Msg("client not found or state unchanged")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Client not found or state unchanged",
			"code":    9,
		})
		return
	}

	log.Info().Str("client_id", id).Bool("state", state).Msg("Client state changed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":      "Client state changed successfully",
		"rowsAffected": rowsAffected,
	})
}
*/
