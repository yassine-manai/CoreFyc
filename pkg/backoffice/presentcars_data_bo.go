package backoffice

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/pkg/db"
)

// ///////////////////////////////////////////////////////////////  DATA BY ID //////////////////////////////////////////////////////////////

// GetAllPresentCars godoc
//
//	@Summary		Get present cars by ID
//	@Tags			Backoffice - PresentCars
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			car_id	query	int	true	"Transaction ID"
//	@Router			/backoffice/get_present_car_id [get]
//	@Description	Get a list of present cars filtered by ID and date range. If no date is provided, it will return all present cars for the current day.
//	@Description	Date format must be YYYY-MM-DD.
//	@Success		200	{object}	db.PresentCar
func GetAllPresentTransactionsDataIDAPI(c *gin.Context) {
	ctx := context.Background()
	ID := c.Query("car_id")

	log.Info().Str("ID", ID).Msg("Received request for present cars by ID")

	if ID == "" {
		log.Warn().Msg("ID is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Please provide an ID to filter the cars",
			"code":    12,
		})
		return
	}

	cars, err := db.GetPresentCarsID(ctx, ID)
	if err != nil || cars == nil {
		log.Err(err).Str("ID", ID).Msg("Error fetching present cars by ID")
		c.JSON(http.StatusOK, []db.PresentCar{})
	}

	CurZone, err := db.GetZoneByID(ctx, *cars.CurrZoneID)
	if err != nil || CurZone == nil {
		log.Error().Err(err).Int("zoneID", *cars.CurrZoneID).Msg("Current Zone not found")
		CurZone = &db.Zone{
			ZoneID: 0,
			Name:   make(map[string]interface{}),
		}
	}

	if !CurZone.IsEnabled || CurZone.IsDeleted {
		log.Warn().Int("Zone ID", CurZone.ZoneID).Msg("Zone is disabled or deleted")
		CurZone.ZoneID = 0
		CurZone.Name = make(map[string]interface{})
	}

	// Fetch and validate camera data
	camera, err := db.GetCameraByID(ctx, cars.CameraID)
	if err != nil || camera == nil {
		log.Error().Err(err).Int("Camera ID", cars.CameraID).Msg("Camera not found")
		camera = &db.ResponseCamera{
			CamID:   0,
			CamName: "",
		}
	} else if !camera.IsEnabled || camera.IsDeleted {
		log.Warn().Int("Camera ID", camera.CamID).Msg("Camera is disabled or deleted")
		camera.CamID = 0
		camera.CamName = ""
	}

	// Fetch and validate zone image data
	ImageZone, err := db.GetZoneImageByZoneID(ctx, *cars.CurrZoneID)
	if err != nil || ImageZone == nil {
		log.Error().Err(err).Int("Zone ID", *cars.CurrZoneID).Msg("Zone image not found")
		ImageZone = &db.ImageZone{
			ImageSm: "",
			ImageLg: "",
		}
	} else if !CurZone.IsEnabled || CurZone.IsDeleted {
		log.Warn().Int("Zone ID", CurZone.ZoneID).Msg("Zone is disabled or deleted")
		ImageZone.ImageSm = ""
		ImageZone.ImageLg = ""
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"id":               cars.ID,
		"license_plate":    cars.LPN,
		"zone_id":          CurZone.ZoneID,
		"zone_name_ar":     CurZone.Name["ar"],
		"zone_name_en":     CurZone.Name["en"],
		"direction":        cars.Direction,
		"transaction_date": cars.TransactionDate,
		"confidence":       cars.Confidence,
		"camera_id":        camera.CamID,
		"camera_name":      camera.CamName,
		"image1":           ImageZone.ImageSm,
		"image2":           ImageZone.ImageLg,
		"cam_event":        "", // to check
	}

	log.Info().Msg("Successfully fetched present cars by ID")
	c.JSON(http.StatusOK, responseData)
}

// GetAllPresentTransactionsDataAPI godoc
//
//	@Summary		Get all present cars
//	@Description	Get a list of all present cars with pagination. If no date is provided, it will return all present cars for the current day.
//	@Description	Date format must be YYYY-MM-DD.
//	@Tags			Backoffice - PresentCars
//	@Produce		json
//	@Security		BearerAuthBackOffice
//	@Param			start			query		string	false	"Include StartDate in the format YYYY-MM-DD"
//	@Param			end				query		string	false	"Include EndDate in the format YYYY-MM-DD"
//	@Param			zoneID			query		string	false	"Zone ID"
//	@Param			licensePlate	query		string	false	"License Plate"
//	@Param			fuzzy_logic		query		bool	false	"Fuzzy Logic"				default(false)
//	@Param			page			query		int		false	"Page number"				default(1)
//	@Param			items_per_page	query		int		false	"Number of items per page"	default(10)
//	@Param			is_present		query		string	false	"Present Car in Parking"	Enums(all, yes, no) default(all)
//	@Success		200				{object}	PaginatedResponse
//	@Router			/backoffice/get_present_car [get]
func GetAllPresentTransactionsDataAPI(c *gin.Context) {
	ctx := context.Background()

	// Retrieve query parameters
	startDate := c.DefaultQuery("start", time.Now().Format("2006-01-02"))
	endDate := c.DefaultQuery("end", time.Now().Format("2006-01-02"))
	licencePlate := c.Query("licensePlate")
	zoneID := c.Query("zoneID")

	// Pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("items_per_page", "10"))
	if err != nil || perPage < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid items_per_page parameter"})
		return
	}

	var response []map[string]interface{}
	var singleResponse map[string]interface{}
	var cars []db.PresentCar

	log.Debug().
		Str("startDate", startDate).
		Str("endDate", endDate).
		Int("page", page).
		Int("perPage", perPage).
		Msg("Received request for present cars")

	// Validate date format
	if err := validateDateFormat(startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}
	if err := validateDateFormat(endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	// Helper function to fetch zone details
	fetchZoneDetails := func(zoneID int) (*db.Zone, error) {
		zone, err := db.GetZoneByID(ctx, zoneID)
		if err != nil || zone == nil || !zone.IsEnabled || zone.IsDeleted {
			if err != nil {
				log.Error().Err(err).Int("zoneID", zoneID).Msg("Current Zone not found")
			}
			return nil, fmt.Errorf("current Zone %v not found", zoneID)
		}
		return zone, nil
	}

	// Fetch data based on provided parameters
	switch {
	case zoneID != "" && licencePlate != "":
		zoneIDInt, _ := strconv.Atoi(zoneID)
		cars, err = db.GetAllLPNZONE(ctx, startDate, endDate, licencePlate, zoneIDInt)
		response, singleResponse = processCars(cars, fetchZoneDetails)
	case zoneID != "":
		zoneIDInt, _ := strconv.Atoi(zoneID)
		cars, err = db.GetZLSE(ctx, startDate, endDate, zoneIDInt)
		response, singleResponse = processCars(cars, fetchZoneDetails)
	case licencePlate != "":
		cars, err = db.GetAllPCLpnBiDateExtra(ctx, startDate, endDate, licencePlate)
		response, singleResponse = processCars(cars, fetchZoneDetails)
	default:
		cars, err = db.GetAllPresentCarsBiDateNoExtra(ctx, startDate, endDate)
		response, singleResponse = processCars(cars, fetchZoneDetails)
	}

	if err != nil {
		handleError(c, err, "Error getting present cars")
		return
	}

	// Pagination calculation
	totalItems := len(response)
	totalPages := (totalItems + perPage - 1) / perPage
	start := (page - 1) * perPage
	end := start + perPage
	if end > totalItems {
		end = totalItems
	}

	var paginatedData []map[string]interface{}
	if start < totalItems {
		paginatedData = response[start:end]
	}

	// Generate links for pagination
	var links []PageLink
	for i := 1; i <= totalPages; i++ {
		url := fmt.Sprintf("/page=%d", i)
		links = append(links, PageLink{
			URL:    url,
			Page:   i,
			Active: i == page,
		})
	}

	// Prepare pagination object
	pagination := Pagination{
		PrevPageURL:  "",
		NextPageURL:  "",
		FirstPageURL: "page=1",
		LastPageURL:  fmt.Sprintf("page=%d", totalPages),
		TotalPages:   totalPages,
		Links:        links,
	}

	if page > 1 {
		pagination.PrevPageURL = fmt.Sprintf("/page=%d", page-1)
	}
	if page < totalPages {
		pagination.NextPageURL = fmt.Sprintf("/page=%d", page+1)
	}

	responsePayload := PaginatedResponse{
		Payload: struct {
			Pagination Pagination `json:"pagination"`
		}{
			Pagination: pagination,
		},
		Data: paginatedData,
	}

	if totalItems == 0 {
		log.Info().Msg("No present cars found")
		responsePayload.Data = []map[string]interface{}{}
		c.JSON(http.StatusOK, responsePayload)
		return
	}

	if totalItems == 1 {
		log.Info().Msg("One data found")
		responsePayload.Data = []map[string]interface{}{singleResponse}
		c.JSON(http.StatusOK, responsePayload)
		return
	}

	log.Info().Msg("Successfully fetched present cars")
	c.JSON(http.StatusOK, responsePayload)
}

// Supporting Structs
type PaginatedResponse struct {
	Payload struct {
		Pagination Pagination `json:"pagination"`
	} `json:"payload"`
	Data []map[string]interface{} `json:"data"`
}

type Pagination struct {
	PrevPageURL  string     `json:"prev_page_url"`
	NextPageURL  string     `json:"next_page_url"`
	FirstPageURL string     `json:"first_page_url"`
	LastPageURL  string     `json:"last_page_url"`
	TotalPages   int        `json:"total_pages"`
	Links        []PageLink `json:"links"`
}

type PageLink struct {
	URL    string `json:"url,omitempty"`
	Page   int    `json:"page,omitempty"`
	Active bool   `json:"active"`
}

// Helper function to validate date format
func validateDateFormat(date string) *gin.H {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Warn().Err(err).Str("date", date).Msg("Invalid date format")
		return &gin.H{
			"error":   "Invalid date format",
			"message": "Date must be in the format YYYY-MM-DD",
			"code":    11,
		}
	}
	return nil
}

// Helper function to handle errors
func handleError(c *gin.Context, err error, message string) {
	log.Err(err).Msg(message)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "An unexpected error occurred",
		"message": message,
		"code":    10,
	})
}

// Helper function to process car data
func processCars(cars []db.PresentCar, fetchZoneDetails func(int) (*db.Zone, error)) ([]map[string]interface{}, map[string]interface{}) {
	var response []map[string]interface{}
	var singleResponse map[string]interface{}
	for _, car := range cars {
		if car.CurrZoneID == nil {
			*car.CurrZoneID = 0
		}

		zone, err := fetchZoneDetails(*car.CurrZoneID)
		if err != nil {
			log.Error().Err(err).Int("zoneID", *car.CurrZoneID).Msg("Error fetching zone details")
			continue
		}

		carData := map[string]interface{}{
			"id":               car.ID,
			"license_plate":    car.LPN,
			"zone_id":          zone.ZoneID,
			"zone_name_ar":     zone.Name["ar"],
			"zone_name_en":     zone.Name["en"],
			"transaction_date": car.TransactionDate,
			"confidence":       car.Confidence,
		}
		response = append(response, carData)

		if len(response) == 1 {
			singleResponse = carData
		}
	}
	return response, singleResponse
}
