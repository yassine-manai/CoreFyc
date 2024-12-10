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

// GetZonesImages godoc
//
//	@Summary		Get all zones images
//	@Description	Get a list of all zones images
//	@Tags			Zones Image
//	@Produce		json
//	@Param			id			query	string	false	"Zone ID to fetch a specific zone's image"
//	@Param			extra		query	bool	false	"Include extra information if 'yes'"
//	@Param			typeImage	query	string	false	"choose the image type Small or Large (small or sm for Small Images / lg or large for Large)"
//	@Success		200			{array}	db.ImageZone
//	@Router			/fyc/zonesImages [get]
func GetAllImageZonesAPI(c *gin.Context) {
	ctx := context.Background()

	log.Info().Msg("GetAllImageZonesAPI called")

	// Get query parameters and log them
	extraReq := c.DefaultQuery("extra", "false")
	imageType := c.DefaultQuery("typeImage", "small")
	idParam := c.Query("id")
	log.Debug().Str("extra", extraReq).Str("typeImage", imageType).Str("idParam", idParam).Msg("Query parameters received")

	var id int
	var err error
	if idParam != "" {
		id, err = strconv.Atoi(idParam)
		if err != nil {
			log.Err(err).Msg("Invalid ID parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
			return
		}
		log.Debug().Int("id", id).Msg("ID parameter converted to int")
	}

	extra := strings.ToLower(extraReq) == "true"
	log.Debug().Bool("extra", extra).Msg("Extra parameter processed")

	if idParam != "" {
		if extra {
			log.Info().Int("id", id).Msg("Fetching zone image with extra data by ID")
			zoneImg, err := db.GetZoneImageByID(ctx, id)
			if err != nil {
				log.Err(err).Msg("Error getting zone image by ID with extra")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting zone image with extra data"})
				return
			}
			if zoneImg == nil {
				log.Info().Int("id", id).Msg("No zone found with specified ID")
				c.JSON(http.StatusNotFound, gin.H{"message": "No zone found with the specified ID"})
				return
			}
			//zoneImg.ImageSm, _ = functions.ByteaToBase64([]byte(zoneImg.ImageSm))
			//zoneImg.ImageLg, _ = functions.ByteaToBase64([]byte(zoneImg.ImageLg))
			log.Info().Int("id", id).Msg("Zone image with extra data fetched successfully")
			c.JSON(http.StatusOK, zoneImg)
		} else {
			log.Info().Int("id", id).Msg("Fetching zone image by ID")
			zoneImg, err := db.GetZoneImageByID(ctx, id)
			if err != nil {
				log.Err(err).Msg("Error getting zone image by ID")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting zone image"})
				return
			}
			if zoneImg == nil {
				log.Info().Int("id", id).Msg("No zone found with specified ID")
				c.JSON(http.StatusNotFound, gin.H{"message": "No zone found with the specified ID"})
				return
			}
			//zoneImg.ImageSm, _ = functions.ByteaToBase64([]byte(zoneImg.ImageSm))
			//zoneImg.ImageLg, _ = functions.ByteaToBase64([]byte(zoneImg.ImageLg))
			log.Info().Int("id", id).Msg("Zone image fetched successfully")
			c.JSON(http.StatusOK, zoneImg)
		}
		return
	}

	log.Debug().Str("imageType", imageType).Msg("Processing image type")

	switch strings.ToLower(imageType) {
	case "small", "sm":
		log.Info().Msg("Fetching small images")
		if idParam != "" {
			zoneImg, err := db.GetZoneImageByID(ctx, id)
			if err != nil {
				log.Err(err).Msg("Error getting small image by ID")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting small image by ID"})
				return
			}
			//zoneImg.ImageSm, _ = functions.ByteaToBase64([]byte(zoneImg.ImageSm))
			log.Info().Int("id", id).Msg("Small image fetched successfully by ID")
			c.JSON(http.StatusOK, zoneImg)
		} else {
			smallImages, err := db.GetAllZoneImageSm(ctx)
			if err != nil {
				log.Err(err).Msg("Error getting all small images")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting small images"})
				return
			}

			log.Info().Int("count", len(smallImages)).Msg("All small images fetched successfully")
			c.JSON(http.StatusOK, smallImages)
		}
		return

	case "large", "lg":
		log.Info().Msg("Fetching large images")
		if idParam != "" {
			zoneImg, err := db.GetZoneImageByID(ctx, id)
			if err != nil {
				log.Err(err).Msg("Error getting large image by ID")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting large image by ID"})
				return
			}
			//zoneImg.ImageLg, _ = functions.ByteaToBase64([]byte(zoneImg.ImageLg))
			log.Info().Int("id", id).Msg("Large image fetched successfully by ID")
			c.JSON(http.StatusOK, zoneImg)
		} else {
			largeImages, err := db.GetAllZoneImageLg(ctx)
			if err != nil {
				log.Err(err).Msg("Error getting all large images")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting large images"})
				return
			}

			log.Info().Int("count", len(largeImages)).Msg("All large images fetched successfully")
			c.JSON(http.StatusOK, largeImages)
		}
		return
	}

	log.Info().Msg("Fetching all zone images")
	smallImages, err := db.GetAllZoneImageSm(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all zone images")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting all zone images"})
		return
	}

	log.Info().Int("count", len(smallImages)).Msg("All zone images fetched successfully")
	c.JSON(http.StatusOK, smallImages)
}

// CreateZone godoc
//
//	@Summary		Add a new zone Image
//	@Description	Add a new zone image to the database
//	@Tags			Zones Image
//	@Accept			json
//	@Produce		json
//	@Param			ImageZone	body		db.ImageZone	true	"Zone image data"
//	@Success		201			{object}	db.ImageZone
//	@Router			/fyc/zonesImage [post]
func CreateZoneImageAPI(c *gin.Context) {
	var zoneImage db.ImageZone

	if err := c.ShouldBindJSON(&zoneImage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if !functions.Contains(db.Zonelist, *zoneImage.ZoneID) {
		log.Debug().Msg("Zone not found")

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", *zoneImage.ZoneID),
			"code":    9,
		})
		return
	}

	ImageSmEnc, err := functions.DecodeBase64ToByteArray(zoneImage.ImageSm)
	if err != nil {
		log.Err(err).Msg("Error converting image SM")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image SM",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ImageLgEnc, err := functions.DecodeBase64ToByteArray(zoneImage.ImageLg)
	if err != nil {
		log.Err(err).Msg("Error converting image LG")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	log.Debug().Msgf("Image Lenght %v", len(ImageLgEnc))
	log.Debug().Msgf("Image Lenght %v", len(ImageSmEnc))

	ctx := context.Background()

	if err := db.CreateZoneImage(ctx, &zoneImage); err != nil {
		log.Err(err).Msg("Error creating new zone image")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new zone image",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	c.JSON(http.StatusCreated, zoneImage)
}

// UpdateZoneImageId godoc
//
//	@Summary		Update a zone image by ID
//	@Description	Update the image data of an existing zone by its ID
//	@Tags			Zones Image
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Zone ID"
//	@Param			Image	body		db.ImageZone	true	"Updated zone image data"
//	@Success		200		{object}	db.Zone
//	@Router			/fyc/zonesImage/{id} [put]
func UpdateZoneImageByIdAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var updates db.ImageZoneNoBind
	ctx := context.Background()

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if !functions.Contains(db.Zonelist, updates.ZoneID) {
		log.Debug().Msg("Zone not found")

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", updates.ZoneID),
			"code":    9,
		})
		return
	}

	ImageSmEnc, err := functions.DecodeBase64ToByteArray(updates.ImageSm)
	if err != nil {
		log.Err(err).Msg("Error converting image SM")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image SM",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ImageLgEnc, err := functions.DecodeBase64ToByteArray(updates.ImageLg)
	if err != nil {
		log.Err(err).Msg("Error converting image LG")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	log.Debug().Int("ZoneID", id).Msgf("Image Lenght %v", len(ImageLgEnc))
	log.Debug().Int("ZoneID", id).Msgf("Image Lenght %v", len(ImageSmEnc))

	// Call the service to update the present car
	rowsAffected, err := db.UpdateZoneImage(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating zone image by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update zone image",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zone image found with the specified ID",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Zone Image modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		//"response":      updates,
	})
}

// DeleteZoneImage godoc
//
//	@Summary		Delete a zone image
//	@Description	Delete a zone image by ID
//	@Tags			Zones Image
//	@Param			id	path		int						true	"Zone image ID"
//	@Success		200	{object}	map[string]interface{}	"Zone image deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid request"
//	@Failure		404	{object}	map[string]interface{}	"Zone image not found"
//	@Router			/fyc/zonesImage/{id} [delete]
func DeleteZoneImageAPI(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := db.DeleteZoneImage(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting zone image")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete zone image",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Debug().Int("ZoneImage ID", id).Msg("Delete Operation ")

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zone image found with the specified ID ------  Affected rows 0 ",
			"code":    9,
		})
		return
	}

	log.Debug().Int("ZoneImage ID", id).Msg("Deleted Successfully")

	c.JSON(http.StatusOK, gin.H{
		"success":      "Zone Image deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
