package third_party

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/config"
	"fyc/middleware"
	"fyc/pkg/db"
)

type TokenRequester struct {
	ClientID     string `form:"client_id" `
	ClientSecret string `form:"client_secret" `
	GrantType    string `form:"grant_type" default:"client_credentials"`
}
type TokenRespose struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
type CarLocation struct {
	ZoneName     string `json:"zone_name"`
	SpotID       string `json:"spot_id"`
	PictureName  string `json:"picture_name"`
	LicensePlate string `json:"license_plate"`
}
type FindMyCarResponse struct {
	ResponseCode int           `json:"response_code"`
	Locations    []CarLocation `json:"locations"`
}
type PictureResponse struct {
	ImageData string `json:"picture_data"`
}

// @Summary		Get an access token
//
// @Description	Get an access token using client credentials
// @Tags			Third Party
// @Produce		json
// @Param			client_id		formData	string	true	"Client ID"
// @Param			client_secret	formData	string	true	"Client Secret"
// @Param			grant_type		formData	string	true	"Grant Type"
// @Success		200				{object}	TokenRespose
// @Router			/token [post]
func GetToken(c *gin.Context) {
	ctx := context.Background()
	var tokenPref = config.Configvar.App.TokenPref3rdParty
	var TokenRequester TokenRequester
	var (
		ClientID     string
		ClientSecret string
		GrantType    string
		IsEnabled    bool
		FuzzyLogic   bool
	)

	if err := c.ShouldBind(&TokenRequester); err != nil {
		log.Err(err).Msg("Error binding token")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/// function that returns 500 in case of panic error
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    -500,
				"message": "An unexpected error occurred. Please try again later.",
			})
		}
	}()

	exists, clientDetails := isClientExist(db.ClientDataList, TokenRequester.ClientID)
	log.Debug().Str("Client ID", clientDetails.ClientID).Bool("Fuzzic Logic", clientDetails.FuzzyLogic).Bool("Client Status", clientDetails.ClientActive).Msg("Fetch Client Data ")

	if !exists {
		log.Warn().Str("ClientID", ClientID).Msg("Client Not Found")
	} else {
		ClientID = clientDetails.ClientID
		ClientSecret = clientDetails.ClientSecret
		GrantType = clientDetails.ClientGrantType
		IsEnabled = clientDetails.ClientActive
		FuzzyLogic = clientDetails.FuzzyLogic
	}

	var missingFields []string
	if TokenRequester.ClientID == "" {
		missingFields = append(missingFields, "ClientID")
	}
	if TokenRequester.ClientSecret == "" {
		missingFields = append(missingFields, "ClientSecret")
	}
	if TokenRequester.GrantType == "" {
		missingFields = append(missingFields, "GrantType")
	}

	// 400 Done - StatusBadRequest
	if len(missingFields) > 0 {
		log.Warn().Strs("Missing Fields", missingFields).Msg("Missing required fields")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    -5,
			"message": "Invalid request. " + strings.Join(missingFields, ", ") + " parameter are required",
		})
		return
	}

	// 401 Done - Unauthorized
	if TokenRequester.ClientID != ClientID || TokenRequester.ClientSecret != ClientSecret {
		log.Warn().Str("ClientID", TokenRequester.ClientID).Str("ClientSecret", TokenRequester.ClientSecret).Msg("Invalid ClientID or ClientSecret")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    -1,
			"message": "Invalid ClientID or ClientSecret",
		})
		return
	}

	// 401 Done - Unauthorized
	if TokenRequester.GrantType != GrantType {
		log.Warn().Str("GrantType", TokenRequester.GrantType).Msg("Invalid GrantType")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    -1,
			"message": "Invalid GrantType",
		})
		return
	}

	// 403 Done - UserDisabled
	if !IsEnabled {
		log.Warn().Str("ClientID", TokenRequester.ClientID).Msg("User is disabled")
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    -2,
			"message": "User is disabled",
		})
		return
	}

	token, err := middleware.GenerateTokenThirdParty(TokenRequester.ClientID, TokenRequester.ClientSecret, TokenRequester.GrantType, FuzzyLogic)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Could not generate token",
			"code":    -500,
		})
		return
	}

	stored, err := db.StoreToken(ctx, ClientID, token)
	if stored == 0 {
		log.Warn().Int("Rows Affeected :", int(stored)).Str("ClientID", ClientID).Msg("Failed to store token")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Could not store token",
			"code":    -7,
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Str("ClientID", ClientID).Msg("Error storing token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    -500,
		})
	}
	if tokenPref == "true" {
		// 200 status - Done
		log.Info().Str("ClientID", ClientID).Msg("Success Generting Token with prefix")
		c.JSON(http.StatusOK, TokenRespose{
			AccessToken: fmt.Sprint("Bearer " + token),
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		})
		return
	}

	log.Info().Str("ClientID", ClientID).Msg("Success Generting Token without Prefix")
	c.JSON(http.StatusOK, TokenRespose{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	})
}

// @Summary		Find a car by license plate
// @Description	Find a car using the license plate number
// @Tags			Third Party
// @Accept			json
// @Produce		json
// @Param			license_plate	query	string	true	"License Plate"
// @Param			language		query	string	false	"Language"	default(en)
// @Security		BearerAuth3rdParty
// @Success		200	{object}	CarLocation
// @Router			/findmycar [get]
func FindMyCar(c *gin.Context) {
	var carResponses []CarLocation
	authHeader := c.GetHeader("Authorization")

	fuzzy_logic, ClientId, err := Extract_token_data(authHeader)
	if err != nil {
		log.Err(err).Msg("Error Getting Fuzzy Logic")

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid authorization header format",
			"code":    -1,
		})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    -500,
				"message": "An unexpected error occurred. Please try again later.",
			})
		}
	}()

	//fuzzy_logic := c.DefaultQuery("fuzzy_logic", "false")
	licensePlate := c.Query("license_plate")
	lang := c.DefaultQuery("language", "en")
	var language = strings.ToLower(lang)
	ctx := context.Background()

	log.Info().Str("Language Provided ", language).Msg("Request Get Picture ")
	log.Debug().Str("Licence Plate", licensePlate).Msg("Find Licence Plate data in progress")

	if licensePlate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Please provide a license plate number",
			"code":    -5,
		})
		return
	}

	////////////////////////////////////////////////////// FALSE
	if !fuzzy_logic {
		log.Info().Bool("Fuzzy Logic", fuzzy_logic).Str("Client ID", ClientId).Msg("Accepetd Request with ")

		car, err := db.GetPresentCarByLPN(ctx, licensePlate)
		if err != nil {

			log.Warn().Str("Error : ", err.Error()).Str("license_plate", licensePlate).Msg("Error retrieving car by LPN")
			c.JSON(http.StatusOK, []CarLocation{})
			return
		}

		log.Info().Str("License Plate", licensePlate).Msg("Car found with license plate")
		log.Debug().Int("zone", *car.CurrZoneID).Msg("Last Zone ID")
		log.Debug().Interface("zone", car).Msg("Present Data")

		spotID := *car.CurrZoneID
		licensePlate = car.LPN

		log.Debug().Str("Spot ID", fmt.Sprint(spotID)).Str("License Plate", licensePlate).Msg("Zone Found")

		// Retrieve the zone image by ID and language
		zoneImage, err := db.GetZoneImageByZONEIDLang(ctx, spotID, language)
		if err != nil {
			log.Warn().Str("Error : ", err.Error()).Str("Language Provided", language).Int("Car Detail ID", *car.CarDetailsID).Msg("Error retrieving zone image")
			c.JSON(http.StatusOK, []CarLocation{})
			return
		}

		// Get Zone name
		zoneData, err := db.GetZoneByID(ctx, spotID)
		if err != nil {
			log.Warn().Str("Error: ", err.Error()).Int("Zone ID", *car.CurrZoneID).Msg("Error retrieving zone")
			c.JSON(http.StatusOK, []CarLocation{})
			return
		}

		// Debug log for zone names
		log.Debug().Interface("Zone Names", zoneData.Name)

		if zoneName, ok := zoneData.Name[language]; ok {
			log.Info().Str("Language", language).Interface("Name", zoneName).Interface("Name", zoneData.Name[language]).Msg("Zone Name")
		} else {
			log.Warn().Msg("Zone name not found for the specified language")
		}

		log.Debug().Str("Found Picture for Car with license plate", licensePlate).Str("Picture Name", fmt.Sprint(zoneImage.ID))

		// Prepare the successful response

		response := CarLocation{
			ZoneName:     fmt.Sprint(zoneData.Name[language]),
			LicensePlate: licensePlate,
			SpotID:       fmt.Sprint(spotID),
			PictureName:  fmt.Sprint(zoneImage.ID),
		}

		carResponses = append(carResponses, response)

		c.JSON(http.StatusOK, carResponses)
	}

	////////////////////////////////////////////////////// TRUE
	if fuzzy_logic {
		log.Info().Bool("Fuzzy Logic", fuzzy_logic).Str("Client ID", ClientId).Msg("Accepetd Request with ")
		var cars []db.PresentCar
		cars, err := db.GetPresentCarByLPNs(ctx, licensePlate)
		if err != nil {
			log.Warn().Str("Error", err.Error()).Str("license_plate", licensePlate).Msg("Error retrieving car by LPN")
			c.JSON(http.StatusOK, []CarLocation{})
			return
		}

		// Check if any car was found
		if len(cars) == 0 {
			log.Warn().Str("license_plate", licensePlate).Msg("No car found with the provided license plate")
			c.JSON(http.StatusOK, []CarLocation{})
			return
		}

		var carResponses []CarLocation

		for _, car := range cars {
			log.Debug().Int("zone", *car.CurrZoneID).Msg("Last Zone ID")
			log.Debug().Interface("car_data", car).Msg("Present Data")

			spotID := *car.CurrZoneID
			licensePlate = car.LPN

			log.Debug().Str("Spot ID", fmt.Sprint(spotID)).Str("License Plate", licensePlate).Msg("Zone Found")

			zoneImages, err := db.GetZoneImageByZONEIDLangs(ctx, spotID, language)
			if err != nil {
				log.Warn().Str("Error", err.Error()).Str("Language Provided", language).Int("Car Detail ID", *car.CarDetailsID).Msg("Error retrieving zone images")
				continue
			}

			// Get Zone name
			zoneData, err := db.GetZoneByID(ctx, spotID)
			if err != nil {
				log.Warn().Str("Error", err.Error()).Int("Zone ID", *car.CurrZoneID).Msg("Error retrieving zone")
				continue
			}

			// Debug log for zone names
			log.Debug().Interface("Zone Names", zoneData.Name)

			var zoneName string
			if name, ok := zoneData.Name[language]; ok {
				zoneName = fmt.Sprint(name)
				log.Info().Str("Language", language).Interface("Name", zoneData.Name[language]).Msg("Zone Name")
			} else {
				log.Warn().Msg("Zone name not found for the specified language")
				zoneName = "Unknown"
			}

			for _, zoneImage := range zoneImages {
				log.Debug().Str("Found Picture for Car with license plate", licensePlate).Str("Picture Name", fmt.Sprint(zoneImage.ID))

				// Prepare the response for each image
				response := CarLocation{
					ZoneName:     zoneName,
					LicensePlate: licensePlate,
					SpotID:       fmt.Sprint(spotID),
					PictureName:  fmt.Sprint(zoneImage.ID),
				}

				// Append the response for each image to the list of car responses
				carResponses = append(carResponses, response)
			}
		}

		// Return all responses for all cars and their images
		c.JSON(http.StatusOK, carResponses)
	}
}

// @Summary		Get a picture by picture name
// @Description	Get an image using the picture name
// @Tags			Third Party
// @Produce		json
// @Param			picture_name	query	string	true	"Picture Name"
// @Param			picture_size	query	string	false	"Picture Size small is default size 'big or small'"	default(small)
// @Security		BearerAuth3rdParty
// @Success		200	{object}	PictureResponse
// @Router			/getpicture [get]
func GetPicture(c *gin.Context) {
	pictureName := c.Query("picture_name")
	imageSize := c.DefaultQuery("picture_size", "small")
	//language := c.DefaultQuery("language", "en")
	//var lang = strings.ToLower(language)

	//log.Info().Str("Language Provided ", lang).Msg("Request Get Picture ")

	ctx := context.Background()

	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    -500,
				"message": "An unexpected error occurred. Please try again later.",
			})
		}
	}()

	if pictureName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Please provide a picture Name",
			"code":    -5,
		})
		return
	}

	if pictureName != "" {
		// Convert pictureName to ID
		id, err := strconv.Atoi(pictureName)
		log.Info().Str("Picture ID ", pictureName).Msg("Fetching Picture by ID")

		if err != nil {
			log.Err(err).Str("id", pictureName).Msg("Invalid Picture ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID must be a valid integer",
				"code":    -5,
			})
			return
		}

		// Fetch the image for the given zone ID and language
		zoneImg, err := db.GetZoneImageByID(ctx, id)
		if err != nil {
			log.Warn().Str("zoneImg ID ", pictureName).Str("language", zoneImg.Language).Msgf("Error retrieving zone image by ID and language: %v", err)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Image not found for the specified language",
				"code":    -4,
			})
			return
		}

		// Check the image size and return the appropriate image
		switch imageSize {
		case "small":
			if zoneImg.ImageSm != "" {
				log.Info().Str("Small Image", pictureName).Str("Language", zoneImg.Language).Int("Zone ID", id).Msg("Small Image fetched successfully")
				c.JSON(http.StatusOK, PictureResponse{
					ImageData: zoneImg.ImageSm,
				})
				return
			}
		case "big":
			if zoneImg.ImageLg != "" {
				log.Info().Str("Big Image", pictureName).Str("Language", zoneImg.Language).Int("Zone ID", id).Msg("Big Image fetched successfully")
				c.JSON(http.StatusOK, PictureResponse{
					ImageData: zoneImg.ImageLg,
				})
				return
			}
		default:
			// Default to small image if size is not specified correctly
			if zoneImg.ImageSm != "" {
				log.Info().Str("Small Image", pictureName).Str("Language", zoneImg.Language).Int("Zone ID", id).Msg("Default Small Image fetched successfully")
				c.JSON(http.StatusOK, PictureResponse{
					ImageData: zoneImg.ImageSm,
				})
				return
			}
		}

		// If the image size doesn't match or is missing
		log.Warn().Str("zoneImg", pictureName).Str("language", zoneImg.Language).Msg("Image not found for the requested size")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Image not found for the requested size",
			"code":    -4,
		})
	}
}

// @Summary	Get Settings
// @Tags		Third Party
// @Accept		application/x-www-form-urlencoded
// @Produce	json
// @Success	200	{object}	db.Settings
// @Security	BearerAuth3rdParty
// @Router		/getsettings [get]
func Getsettings(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	FuzzyLogicValue, _, err := Extract_token_data(authHeader)
	if err != nil {
		log.Err(err).Msg("Error Getting Fuzzy Logic")

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid authorization header format",
			"code":    -1,
		})
		return
	}

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

	ctx := context.Background()
	settings, err := db.GetAllSettingsThirdParty(ctx)
	if err != nil {
		log.Warn().Str("Error ", err.Error()).Msg("Error retrieving Settings fromqsdq db")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Settings not found",
			"code":    -4,
		})
		return
	}

	type SettingsResponse3rdParty struct {
		CarParkID          int                    `json:"carpark_id"`
		CarParkName        map[string]interface{} `json:"carpark_name"`
		AppLogo            string                 `json:"app_logo"`
		DefaultLang        string                 `json:"default_lang"`
		TimeOutScreenKisok *int                   `json:"timeout_screenKiosk"`
		FuzzyLogic         bool                   `json:"fuzzy_logic"`
	}

	log.Info().Int("CarPark ID", settings.CarParkID).Bool("FuzzyLogic", FuzzyLogicValue).Msg("Settings fetched successfully")

	c.JSON(http.StatusOK, SettingsResponse3rdParty{
		CarParkID:          settings.CarParkID,
		CarParkName:        settings.CarParkName,
		AppLogo:            settings.AppLogo,
		DefaultLang:        settings.DefaultLang,
		TimeOutScreenKisok: settings.TimeOutScreenKisok,
		FuzzyLogic:         FuzzyLogicValue,
	})
}
