package backoffice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/functions"
	"fyc/pkg/db"
)

type AddZone struct {
	ID           int                    `json:"-"`
	ZoneID       int                    `json:"zone_id" binding:"required"`
	Name         map[string]interface{} `json:"name" binding:"required" swaggertype:"object"`
	MaxCapacity  *int                   `json:"max_capacity" default:"0"`
	FreeCapacity *int                   `json:"free_capacity" default:"9999"`
	LastUpdated  string                 `json:"-"`
	Images       json.RawMessage        `json:"images" binding:"required" swaggertype:"object"`
	Extra        map[string]interface{} `json:"extra" swaggertype:"object"`
}

type Image struct {
	EN Imagels `json:"en"`
	AR Imagels `json:"ar"`
}

type Imagels struct {
	Image_s string `json:"image_s"`
	Image_l string `json:"image_l"`
}

type UpdateZone struct {
	ID           int                    `json:"-"`
	ZoneID       int                    `json:"zone_id"`
	Name         map[string]interface{} `json:"name" swaggertype:"object"`
	MaxCapacity  *int                   `json:"max_capacity"`
	FreeCapacity *int                   `json:"free_capacity"`
	LastUpdated  string                 `json:"-"`
	Images       json.RawMessage        `json:"images" swaggertype:"object"`
	Extra        map[string]interface{} `json:"extra" swaggertype:"object"`
	IsDeleted    *bool                  `json:"-"`
	IsEnabled    *bool                  `json:"is_enabled"`
}

// GetZonesAPI godoc
//
//	@Summary		Get all zones or a zone by ID
//	@Description	Get a list of all zones, or a zone by ID if 'zone_id' parameter is provided
//	@Tags			Backoffice - Zone
//	@Produce		json
//	@Param			zone_id	query	int		false	"Zone ID"
//	@Param			extra	query	bool	false	"Include extra information if 'true'"
//	@Security		BearerAuthBackOffice
//	@Router			/backoffice/get_zones [get]
//	@Success		200	{array}	models.ZoneDataModel2	"List of zones or a single zone"
func GetZonesAPI(c *gin.Context) {
	ctx := context.Background()
	extraReq := strings.ToLower(c.DefaultQuery("extra", "false"))
	idStr := c.Query("zone_id")

	// If zone_id is provided, fetch the specific zone
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid zone ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID must be a valid integer",
				"code":    -5,
			})
			return
		}

		if extraReq == "true" {
			log.Info().Int("Zone ID", id).Msg("Fetching zone by ID with extra information")
			zone, err := db.GetZoneByID(ctx, id)

			if err != nil {
				log.Err(err).Int("Zone ID", id).Msg("No data found ")
				c.JSON(http.StatusOK, []db.Zone{})
				return
			}

			log.Info().Int("Zone ID", id).Msg("Fetching zoneImage data")
			zoneImg, err := db.GetAllImagesbyZoneID(ctx, id)
			if err != nil {
				log.Err(err).Int("Zone ID", id).Msg("No data found ")
				c.JSON(http.StatusOK, []db.Zone{})
				return
			}

			images := gin.H{}
			for _, img := range zoneImg {
				if img.Language != "" {
					images[img.Language] = gin.H{
						"image_l": img.ImageLg,
						"image_s": img.ImageSm,
					}
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"free_capacity": zone.FreeCapacity,
				"images":        images,
				"is_enabled":    zone.IsEnabled,
				"last_update":   zone.LastUpdated,
				"max_capacity":  zone.MaxCapacity,
				"name":          zone.Name,
				"zone_id":       zone.ZoneID,
				"extra":         zone.Extra,
			})
			return

		}
		if extraReq == "false" {
			var zone *db.Zone
			log.Info().Int("Zone ID", id).Msg("Fetching zone by ID without extra information")
			zone, err := db.GetZoneByID(ctx, id)
			if err != nil {
				log.Err(err).Int("Zone ID", id).Msg("Not Found data for ")
				c.JSON(http.StatusOK, []db.Zone{})
				return
			}

			log.Info().Int("Zone ID", id).Msg("Fetching zoneImage data")
			zoneImg, err := db.GetAllImagesbyZoneID(ctx, id)
			if err != nil {
				log.Err(err).Int("Zone ID", id).Msg("Not Found data for ")
				c.JSON(http.StatusOK, []db.Zone{})
				return
			}

			// Structure the images dynamically by language
			images := gin.H{}
			for _, img := range zoneImg {
				if img.Language != "" {
					images[img.Language] = gin.H{
						"image_l": img.ImageLg,
						"image_s": img.ImageSm,
					}
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"free_capacity": zone.FreeCapacity,
				"images":        images,
				"is_enabled":    zone.IsEnabled,
				"last_update":   zone.LastUpdated,
				"max_capacity":  zone.MaxCapacity,
				"name":          zone.Name,
				"zone_id":       zone.ZoneID,
			})
			return
		}
	}

	// If no zone_id is provided, fetch all zones
	log.Info().Str("Extra Param", extraReq).Msg("Fetching all zones")

	if extraReq == "true" {
		// Fetch all zones with extra info
		zones, err := db.GetAllZoneExtra(ctx)
		if err != nil {
			log.Err(err).Str("Extra Requst", extraReq).Msg("Not Found data for ")
			c.JSON(http.StatusOK, []db.Zone{})
			return
		}

		if len(zones) == 0 {
			log.Debug().Str("Extra Requst", extraReq).Msg("Not Found data for ")
			c.JSON(http.StatusOK, []db.Zone{})
			return
		}
		c.JSON(http.StatusOK, zones)
		return
	}

	if extraReq == "false" {

		zones, err := db.GetAllZoneNoExtra(ctx)
		if err != nil {
			log.Debug().Str("Extra Requst", extraReq).Msg("Not Found data for ")
			c.JSON(http.StatusOK, []db.Zone{})
			return
		}

		if len(zones) == 0 {
			c.JSON(http.StatusOK, []db.Zone{})
			return
		}

		c.JSON(http.StatusOK, zones)
	}
}

// GetZonesNames godoc
//
//	@Summary		Get all zones Names
//	@Description	Get a list of all zones names
//	@Tags			Backoffice - Zone
//	@Security		BearerAuthBackOffice
//	@Produce		json
//	@Success		200	{array}	models.ZoneNamesModel	"List of zones or a single zone"
//	@Router			/backoffice/get_zones_names [get]
func GetZonesNames(c *gin.Context) {
	ctx := context.Background()

	log.Debug().Msg("Fetching all zone names")

	defer func() {
		if r := recover(); r != nil {
			log.Warn().Interface("Error ", r).Msg("An unexpected error occurred.")

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    -500,
				"message": "An unexpected error occurred. Please try again later.",
			})
		}
	}()

	zones, err := db.GetAllZoneExtra(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all zones")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if len(zones) == 0 {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	var respZoneNames []gin.H
	for _, zone := range zones {
		respZoneNames = append(respZoneNames, gin.H{
			"zone_id":   &zone.ZoneID,
			"zone_name": zone.Name,
		})
	}
	log.Debug().Int("Zone Length", len(zones))
	log.Debug().Int("Zone Length", len(respZoneNames))

	// Return the filtered response
	c.JSON(http.StatusOK, respZoneNames)
}

// CreateZoneAPI adds a new zone to the database
//
//	@Summary		Add a new zone
//	@Description	Add a new zone to the database
//	@Tags			Backoffice - Zone
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			zone	body	models.AddZoneModel	true	"Zone data"
//	@Router			/backoffice/add_zone [post]
func CreateZone(c *gin.Context) {
	var addZone AddZone
	ctx := context.Background()

	if err := c.ShouldBindJSON(&addZone); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for zone creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -5,
		})
		return
	}
	if functions.Contains(db.AllZonelist, addZone.ZoneID) {
		log.Warn().Int("Zone ID", addZone.ZoneID).Msg("Zone already exists")
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": fmt.Sprintf("Zone ID %d already exists", addZone.ZoneID),
			"code":    -5,
		})
		return
	}

	if addZone.MaxCapacity != nil && addZone.FreeCapacity != nil {
		if *addZone.MaxCapacity < *addZone.FreeCapacity {
			log.Warn().Int("MaxCapacity", *addZone.MaxCapacity).
				Int("FreeCapacity", *addZone.FreeCapacity).
				Msg("Invalid capacity values")

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Free capacity cannot exceed or equal maximum capacity",
				"code":    -5,
			})
			return
		}
	}

	if addZone.Extra == nil {
		addZone.Extra = make(map[string]interface{})
	}

	// Validate and process images before creating the zone
	//log.Debug().Interface("img", addZone.Images).Send()
	fmt.Println(reflect.TypeOf(addZone.Images))

	// Create the zone in the database
	zone := db.Zone{
		ZoneID:       addZone.ZoneID,
		Name:         addZone.Name,
		MaxCapacity:  addZone.MaxCapacity,
		FreeCapacity: addZone.FreeCapacity,
		Extra:        addZone.Extra,
	}

	var zoneImages map[string]map[string]string
	err := json.Unmarshal(addZone.Images, &zoneImages)
	if err != nil {
		log.Err(err).Msg("Error unmarshalling to map:")
		// db.DeleteZoneError(ctx, zone.ZoneID)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if err := db.CreateZone(ctx, &zone); err != nil {
		log.Error().Err(err).Int("Zone ID", addZone.ZoneID).Msg("Error creating zone")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	// for lang, images := range zoneImages {
	// 	log.Debug().Msgf("lang: %v,images: %v", lang, len(images))
	// 	for key, value := range images {
	// 		log.Warn().Msgf("key: %v,value: %v", key, len(value))
	// 	}
	// }

	for lang, images := range zoneImages {
		log.Debug().Str("Lang", lang).Int("ZoneID", addZone.ZoneID).Msg("Working on image to add for ")
		ImageZone := db.ImageZone{
			ZoneID:   &addZone.ZoneID,
			Language: strings.ToLower(lang),
			ImageSm:  images["image_s"],
			ImageLg:  images["image_l"],
		}

		if err := db.CreateZoneImage(ctx, &ImageZone); err != nil {
			log.Error().Err(err).Int("Zone ID", *ImageZone.ZoneID).Msg("Error adding zone image for language  Delete inserted data" + ImageZone.Language)
			db.DeleteZoneError(ctx, zone.ZoneID)
			db.DeleteZoneImage(ctx, zone.ZoneID) //must change it to delete image zone id
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Error occurred while adding zone image for language: %s", ImageZone.Language),
				"code":    -500,
			})

			return
		}
	}

	log.Debug().Int("ZoneID", addZone.ZoneID).Msg("Zone successfully created")
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Zone successfully created",
	})
}

// UpdateZoneAPI updates an existing zone in the database
//
//	@Summary		Update an existing zone
//	@Description	Update an existing zone in the database
//	@Tags			Backoffice - Zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id	query	int	true	"Zone ID"
//	@Security		BearerAuthBackOffice
//	@Param			zone	body	models.UpdateZoneModel	true	"Updated Zone data"
//	@Router			/backoffice/update_zone [put]
func UpdateZoneDataAPI(c *gin.Context) {
	var Zone2update UpdateZone
	ctx := context.Background()

	zone_id := c.Query("zone_id")
	log.Info().Str("id", zone_id).Msg("Updating Zone ")

	if zone_id == "" {
		log.Warn().Msg("The Zone ID is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "The Zone ID is required",
			"code":    -5,
		})
		return
	}

	id, err := strconv.Atoi(zone_id)
	if err != nil {
		log.Error().Str("zone_id", zone_id).Msg("Invalid ID format for zone update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    -5,
		})
		return
	}

	if !functions.Contains(db.AllZonelist, id) {
		log.Warn().Int("Zone ID", id).Msg("Zone not found")
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": fmt.Sprintf("Zone ID %d not found", id),
			"code":    -5,
		})
		return
	}

	log.Debug().Interface("Zone Up", Zone2update).Send()

	if err := c.ShouldBindJSON(&Zone2update); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for zone creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -5,
		})
		return
	}

	if Zone2update.MaxCapacity != nil && Zone2update.FreeCapacity != nil {
		if *Zone2update.MaxCapacity < *Zone2update.FreeCapacity {
			log.Warn().
				Int("MaxCapacity", *Zone2update.MaxCapacity).
				Int("FreeCapacity", *Zone2update.FreeCapacity).
				Msg("Invalid capacity values")

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Free capacity cannot exceed or equal maximum capacity",
				"code":    -5,
			})
			return
		}
	}

	//log.Debug().Interface("Data", Zone2update).Int.Send()

	zone := db.ZoneNoBind{
		//ZoneID:       id,
		Name: Zone2update.Name,
		//MaxCapacity: &updateZone.MaxCapacity,
		MaxCapacity:  Zone2update.MaxCapacity,
		FreeCapacity: Zone2update.FreeCapacity,
		Extra:        Zone2update.Extra,
		IsEnabled:    Zone2update.IsEnabled,
		IsDeleted:    Zone2update.IsDeleted,
	}

	//log.Debug().Interface("Data -*-*-*-* ", updateZone).Send()
	//log.Debug().Msgf("Data --- %v ***************", Zone2update)

	rows_affected, err := db.UpdateZone(ctx, id, zone)
	if err != nil {
		log.Err(err).Int("Zone ID", id).Msg("Error updating zone")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "An unexpected error occurred. Please try again later.",
			"code":    -500,
		})
		return
	}

	if rows_affected == 0 {
		log.Info().Int64("Rows Affected", rows_affected).Msg("Error Updating zone ")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No Row updated",
			"code":    -7,
		})
		return
	}

	if len(Zone2update.Images) != 0 {
		log.Warn().Msg("Body request has Images , update ZoneImage")

		var zoneImages map[string]map[string]string
		err = json.Unmarshal(Zone2update.Images, &zoneImages)
		if err != nil {
			log.Err(err).Msg("Error unmarshalling to map:")
		}

		log.Debug().Int("Image Length", len(Zone2update.Images)).Int("Zone ID ", Zone2update.ZoneID)
		for lang, images := range zoneImages {
			log.Debug().Msgf("lang: %v", lang)
			for key := range images {
				log.Warn().Msgf("key: %v ", key)
			}
		}

		for lang, images := range zoneImages {
			log.Debug().Str("Lang", lang).Int("ZoneID", id).Msg("Working on image to update for ")

			ImageZone := db.ImageZoneNoBind{
				Language: strings.ToLower(lang),
				ImageSm:  images["image_s"],
				ImageLg:  images["image_l"],
			}

			//log.Debug().Int("Zone ID", id).Str("Lang", lang).Msg("----- #Zone image updated success# ----")

			rows_affected, err := db.UpdateZoneImage(ctx, id, &ImageZone)
			if err != nil {
				log.Error().Err(err).Int("Zone ID", ImageZone.ZoneID).Msg("Error adding zone image for language " + ImageZone.Language)

				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": fmt.Sprintf("Error occurred while adding zone image for language: (%s) -- Please Add All image Types", ImageZone.Language),
					"code":    -500,
				})
				return
			}

			log.Debug().Int("Zone ID", id).Str("Lang", lang).Msg("**** ---- Zone image updated success ---- ****")

			if rows_affected == 0 {
				log.Info().Int64("Rows Affected", rows_affected).Msg("Error Updating zone ")
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Error Updating zone",
					"code":    -7,
				})
				return
			}
		}
	}

	log.Debug().Int("ZoneID", id).Msg("Zone successfully updated")
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Zone successfully updated",
	})
}

// DeleteZoneAPI deletes a zone by its ID
//
//	@Summary		Delete a zone
//	@Description	Delete a zone by ID
//	@Tags			Backoffice - Zone
//	@Security		BearerAuthBackOffice
//	@Param			zone_id	query	int	true	"Zone ID"
//	@Router			/backoffice/delete_zone [delete]
func DeleteZoneDataAPI(c *gin.Context) {

	idStr := c.Query("zone_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Msg("Invalid ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID must be a valid integer",
			"code":    -5,
		})
		return
	}

	if !functions.Contains(db.AllZonelist, id) {
		log.Warn().Int("Zone ID", id).Msg("Zone not exists")
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": fmt.Sprintf("Zone ID %d not exists", id),
			"code":    -5,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := db.DeleteZone(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting Zone")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -500,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No zone found with the specified ID",
			"code":    -7,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Zone deleted successfully",
	})
}
