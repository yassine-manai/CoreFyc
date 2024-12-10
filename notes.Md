------------------------------ NOTES ----------------------------
Step 1: 
export PATH=$(go env GOPATH)/bin:$PATH
Step 2:
swag init 
if error exe this :
swag init --parseDependency

swag init --exclude pkg/api,pkg/debug
swag init --exclude pkg/backoffice,pkg/pka,pkg/third_party

------------------------------- To DO --------------------------
databse/schema creation to take in consideration check it befre everything
 ---- del must be checked for later

------------------------------- DONE -------------------------
Setup Routes 
redirection from /docs to /docs/index.html 
List Return for all with id
add filter by date --- if no date (provide transation of the day) -- 09/10
autocreate of tables, double check it
check get all present cars with extra or without, and extra with lpn  -- 09/10
reader no map eng default --- filter 



------------------------------ NOTES ----------------------------
StatusMethodNotAllowed - 11 -- 405
StatusNotFound - 9  -- 404
StatusBadRequest  - 12 -- 400
StatusUnauthorized- 401
StatusInternalServerError - 10 - 500
StatusOK - 8 - 200 


StatusMethodNotAllowed - 11
StatusInternalServerError -10
StatusNotFound - 9 
StatusBadRequest  - 12
StatusUnauthorized- 13
StatusOK - 8 





	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	log.Debug().Str("Formatted date and time:", formattedTime)
	zone.LastUpdated = formattedTime



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
























check this i wanna make it as one api that contain extra and if and both of them passed as query and not required and if the extra not passed it return the data without extra elseif extra it return data with extra data and the id query if it passed it return the camerawith id if not passed it retun all data 