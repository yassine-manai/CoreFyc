package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/functions"
	"fyc/pkg/db"
)

// GetZonesAPI godoc
//
//	@Summary		Get all zones
//	@Description	Get a list of all zones, or a zone by ID if 'id' parameter is provided
//	@Tags			Zones
//	@Produce		json
//	@Param			id		query	int		false	"Zone ID"
//	@Param			extra	query	bool	false	"Include extra information if 'true'"
//	@Success		200		{array}	db.Zone	"List of zones or a single zone"
//	@Router			/fyc/zones [get]
func GetZonesAPI(c *gin.Context) {
	ctx := context.Background()
	extraReq := strings.ToLower(c.DefaultQuery("extra", "false"))
	idStr := c.Query("id")

	// Check if 'id' is provided
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid zone ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		// If extra info is requested
		if extraReq == "true" {
			log.Info().Int("Zone ID", id).Msg("Fetching zone by ID with extra information")
			zone, err := db.GetZoneByID(ctx, id)
			if err != nil {
				log.Err(err).Int("zone_id", id).Msg("Error retrieving zone by ID with extra info")
				c.JSON(http.StatusNotFound, gin.H{
					"error":   "Not Found",
					"message": "Zone not found",
					"code":    9,
				})
				return
			}
			c.JSON(http.StatusOK, zone)
			return
		}

		// If extra info is NOT requested
		log.Info().Int("Zone ID", id).Msg("Fetching zone by ID without extra information")
		zone, err := db.GetZoneByIDNoExtra(ctx, id)
		if err != nil {
			log.Err(err).Int("zone_id", id).Msg("Error retrieving zone by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Zone not found",
				"code":    9,
			})
			return
		}
		c.JSON(http.StatusOK, zone)
		return
	}

	// If 'id' is not provided, handle all zones case
	log.Info().Str("extra", extraReq).Msg("Fetching all zones")

	if extraReq == "true" {
		// Fetch all zones with extra info
		zones, err := db.GetAllZoneExtra(ctx)
		if err != nil {
			log.Err(err).Msg("Error getting all zones with extra data")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all zones with extra data",
				"code":    10,
			})
			return
		}
		if len(zones) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No zones found",
				"code":    9,
			})
			return
		}
		c.JSON(http.StatusOK, zones)
		return
	}

	// Default or "extra=false", fetch all zones without extra info
	zones, err := db.GetAllZoneNoExtra(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all zones")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all zones",
			"code":    10,
		})
		return
	}
	if len(zones) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zones found",
			"code":    9,
		})
		return
	}
	c.JSON(http.StatusOK, zones)
}

// CreateZoneAPI adds a new zone to the database
//
//	@Summary		Add a new zone
//	@Description	Add a new zone to the database
//	@Tags			Zones
//	@Accept			json
//	@Produce		json
//	@Param			zone	body		db.Zone	true	"Zone data"
//	@Success		201		{object}	db.Zone
//	@Router			/fyc/zones [post]
func CreateZoneAPI(c *gin.Context) {
	var zone db.Zone
	ctx := context.Background()

	// Bind the JSON request to the zone struct
	if err := c.ShouldBindJSON(&zone); err != nil {
		log.Err(err).Msg("Invalid request payload for zone creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if functions.Contains(db.Zonelist, zone.ZoneID) {
		log.Debug().Int("Zone already exists with ID ", zone.ZoneID).Msg("Error creating new zone")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Zone already exists",
			"message": fmt.Sprintf("Zone with ID %d already exists", zone.ZoneID),
			"code":    9,
		})
		return
	}

	if *zone.FreeCapacity > *zone.MaxCapacity {
		log.Debug().Int("Free capacity exceeds max capacity for zone ID", zone.ZoneID).Int("Free Capacity:", *zone.FreeCapacity).Int("Max Capacity:", *zone.MaxCapacity).Msg("Error creating new zone")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid free capacity",
			"message": fmt.Sprintf("Free capacity %d exceeds max capacity %d for zone ID %d", *zone.FreeCapacity, *zone.MaxCapacity, zone.ZoneID),
			"code":    9,
		})
		return
	}

	lowerCaseMap := make(map[string]interface{})
	for key, value := range zone.Name {
		lowerCaseMap[strings.ToLower(key)] = value
	}
	zone.Name = lowerCaseMap

	if err := db.CreateZone(ctx, &zone); err != nil {
		log.Err(err).Msg("Error creating new zone")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new zone",
			"message": err.Error(),
			"code":    10,
		})
		return
	}


	c.JSON(http.StatusCreated, zone)
}

// UpdateZoneIdAPI updates a zone by its ID
//
//	@Summary		Update a zone by ID
//	@Description	Update an existing zone by ID
//	@Tags			Zones
//	@Accept			json
//	@Produce		json
//	@Param			id		query		int		true	"Zone ID"
//	@Param			zone	body		db.Zone	true	"Updated zone data"
//	@Success		200		{object}	db.Zone
//	@Router			/fyc/zones [put]
func UpdateZoneIdAPI(c *gin.Context) {
	// Convert ID param to integer
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Msg("Invalid ID format for update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var updates db.ZoneNoBind
	ctx := context.Background()

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for zone update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	// Check if ZoneID exists
	if !functions.Contains(db.Zonelist, id) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", id),
			"code":    9,
		})
		return
	}

	lowerCaseMap := make(map[string]interface{})
	for key, value := range updates.Name {
		lowerCaseMap[strings.ToLower(key)] = value
	}
	updates.Name = lowerCaseMap

	if *updates.FreeCapacity > *updates.MaxCapacity {
		log.Debug().
			Int("ZoneID", updates.ZoneID).
			Int("FreeCapacity", *updates.FreeCapacity).
			Int("MaxCapacity", *updates.MaxCapacity).
			Msg("Error creating new zone: Free capacity exceeds max capacity")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid free capacity",
			"message": fmt.Sprintf("Free capacity %d exceeds max capacity %d for zone ID %d", updates.FreeCapacity, updates.MaxCapacity, updates.ZoneID),
			"code":    9,
		})
		return
	}

	// Call the service to update the zone
	rowsAffected, err := db.UpdateZone(ctx, id, updates)
	if err != nil {
		log.Err(err).Msg("Error updating zone by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update zone",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zone found with the specified ID",
			"code":    9,
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message":       "Zone modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		//"response":      updates,
	})
}

// DeleteZoneAPI deletes a zone by its ID
//
//	@Summary		Delete a zone
//	@Description	Delete a zone by ID
//	@Tags			Zones
//	@Param			id	query		int						true	"Zone ID"
//	@Success		200	{object}	map[string]interface{}	"Zone deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid request"
//	@Failure		404	{object}	map[string]interface{}	"Zone not found"
//	@Router			/fyc/zones [delete]
func DeleteZoneAPI(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Msg("Invalid ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := db.DeleteZone(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting Zone")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete Zone",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":        "Not Found",
			"message":      "No zone found with the specified ID",
			"rowsAffected": rowsAffected,
			"code":         9,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      "Zone deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}

/* // GetzonesAPI godoc
//
//	@Summary		Get enabled zones or a specific zone by ID
//	@Description	Get a list of enabled zones or a specific zone by ID with optional extra data
//	@Tags			Zones
//	@Produce		json
//	@Param			id		query		string	false	"Zone ID"
//	@Success		200		{array}		db.Zone		"List of enabled zones or a single zone"
//	@Router			/fyc/zonesEnabled [get]
func GetZoneEnabledAPI(c *gin.Context) {
	log.Debug().Msg("Get Zone EnabledAPI request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid Zone ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		Zone, err := GetZoneEnabledByID(ctx, id)
		if err != nil {
			log.Err(err).Str("Zone id", idStr).Msg("Error retrieving Zone by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Zone not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("Zone id", idStr).Msg("Enabled Zone fetched successfully")
		c.JSON(http.StatusOK, Zone)
		return
	}

	// Fetch all enabled zones
	Zone, err := GetZoneListEnabled(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving enabled Zones")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving enabled Zones",
			"code":    10,
		})
		return
	}

	if len(Zone) == 0 {
		log.Info().Msg("No enabled Zones found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No enabled Zone found",
			"code":    9,
		})
		return
	}

	log.Info().Int("Zone_count", len(Zone)).Msg("Enabled Zone fetched successfully")
	c.JSON(http.StatusOK, Zone)
}

// GetzonesAPI godoc
//
//	@Summary		Get Deleted zones or a specific zone by ID
//	@Description	Get a list of Deleted zones or a specific zone by ID with optional extra data
//	@Tags			Zones
//	@Produce		json
//	@Param			id		query		string	false	"Zone ID"
//	@Success		200		{object}	Zone		"List of Deleted zones or a single zone"
//	@Router			/fyc/zonesDeleted [get]
func GetZoneDeletedAPI(c *gin.Context) {
	log.Debug().Msg("Get Zone DeletedAPI request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid Zone ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		Zone, err := GetZoneDeletedByID(ctx, id)
		if err != nil {
			log.Err(err).Str("Zone_id", idStr).Msg("Error retrieving Zone by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Zone not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("Zone_id", idStr).Msg("Deleted Zone fetched successfully")
		c.JSON(http.StatusOK, Zone)
		return
	}

	// Fetch all deleted zones
	Zone, err := GetZoneListDeleted(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving deleted Zone")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving deleted Zone",
			"code":    10,
		})
		return
	}

	if len(Zone) == 0 {
		log.Info().Msg("No deleted Zone found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No deleted Zone found",
			"code":    9,
		})
		return
	}

	log.Info().Int("Zone_count", len(Zone)).Msg("Deleted Zone fetched successfully")
	c.JSON(http.StatusOK, Zone)
}

// ChangeStateAPI godoc
//
//	@Summary		Change Zone state or retrieve Zones by ID
//	@Description	Change the state of a Zone (e.g., enabled/disabled) or retrieve a Zone by ID
//	@Tags			Zones
//	@Produce		json
//	@Param			state	query		bool	false	"Zone State"
//	@Param			id		query		int 	false	"Zone ID"
//	@Success		200		{object}	int		"Number of rows affected by the state change"
//	@Router			/fyc/zoneState [put]
func ChangeZoneStateAPI(c *gin.Context) {
	log.Debug().Msg("Change State API request")
	ctx := context.Background()

	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid zone ID format")
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

	rowsAffected, err := ChangeZoneState(ctx, id, state)
	if err != nil {
		if err.Error() == fmt.Sprintf("zone with id %d is already enabled", id) {
			log.Info().Str("zone_id", idStr).Msg("zone is already enabled")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Error",
				"message": err.Error(),
				"code":    12,
			})
			return
		}

		log.Err(err).Str("zone_id", idStr).Msg("Error changing zone state")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("zone_id", idStr).Msg("Zone not found or state unchanged")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Zone not found or state unchanged",
			"code":    9,
		})
		return
	}

	log.Info().Str("zone_id", idStr).Bool("state", state).Msg("Zone state changed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":      "Zone state changed successfully",
		"rowsAffected": rowsAffected,
	})
}
*/
