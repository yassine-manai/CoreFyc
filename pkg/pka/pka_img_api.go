package pka

import (
	"context"
	"encoding/base64"
	"fyc/pkg/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// PkaImageAPI handles the request to get an image.
//
//	@Summary		PKA SYSTEM API - Search Image
//	@Description	Retrieve an image from the PKA system by image name.
//	@Tags			PKA - API
//	@Produce		image/jpeg, image/png
//	@Param			imagename	path	string	true	"Image Name"	"The name of the image to retrieve, without the '.png' or '.jpeg' extension."
//	@Success		200			{file}	string	"Image retrieved successfully."
//	@Router			/v2/maps/{imagename} [get]
func PkaImageAPI(c *gin.Context) {
	image_name := c.Param("imagename")
	ctx := context.Background()
	var imageData string

	log.Debug().Str(" Name", image_name).Msg("Requesting image")

	if image_name == "" {
		log.Err(nil).Str("Image Param Passed", image_name).Msg("Empty image name parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Empty image name parameter",
			"message": "Image name must not be empty",
			"code":    11,
		})
		return
	}

	// Remove the ".png" or ".jpeg" extension if present
	image_name = strings.TrimSuffix(image_name, ".png")
	image_name = strings.TrimSuffix(image_name, ".jpeg")

	// Parse the remaining part as an integer
	pictureName, err := strconv.Atoi(image_name)
	if err != nil {
		log.Err(err).Str("Image Param Passed", image_name).Msg("Invalid ID format for Image Param")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	zoneImg, err := db.GetZoneImageByZONEIDLang(ctx, pictureName, "en")
	if err != nil {
		log.Err(err).Int("Image ID", pictureName).Msg("Image not found for the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Image not found for the specified ID",
			"code":    -4,
		})
		return
	}

	settingsData, err := db.GetAllSettings(ctx)
	if err != nil || settingsData.PkaImageSize == "" {
		log.Warn().Int("Image ID", pictureName).Msg("Settings not found or empty -- using small image size")
		imageData = zoneImg.ImageLg
	} else {
		switch settingsData.PkaImageSize {
		case "large":
			imageData = zoneImg.ImageLg
		case "small":
			imageData = zoneImg.ImageSm
		default:
			imageData = zoneImg.ImageSm
		}
	}

	if len(imageData) == 0 || (!strings.HasPrefix(imageData, "data:image/jpeg;base64,") && !strings.HasPrefix(imageData, "data:image/png;base64,")) {
		log.Err(err).Int("Image Data for ZoneID", *zoneImg.ZoneID).Msg("Invalid image data format")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    -500,
			"message": "Image data is invalid.",
		})
		return
	}

	var decodedData []byte
	var contentType string

	if strings.HasPrefix(imageData, "data:image/jpeg;base64,") {
		contentType = "image/jpeg"
		decodedData, err = base64.StdEncoding.DecodeString(imageData[len("data:image/jpeg;base64,"):])
	} else if strings.HasPrefix(imageData, "data:image/png;base64,") {
		contentType = "image/png"
		decodedData, err = base64.StdEncoding.DecodeString(imageData[len("data:image/png;base64,"):])
	}

	if err != nil {
		log.Err(err).Msg("Error decoding the image data")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    -500,
			"message": "An unexpected error occurred while decoding the image data.",
		})
		return
	}

	log.Info().Str(" Name", image_name).Int("Zone ID", *zoneImg.ZoneID).Msg("Image Successfully Fetched")
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, decodedData)
}
