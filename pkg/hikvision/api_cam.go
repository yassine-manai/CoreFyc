package hikvision

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/config"
)

func HikvisionHandler(c *gin.Context) {

	now := time.Now()
	formattedTime := now.Format("2006/01/02 15:04:05")
	clientIP := c.ClientIP()
	var DataCapt Capture

	log.Info().Msgf("Capture Object Initialized%vat %v", DataCapt, formattedTime)

	// Check multipart form data
	form, err := c.MultipartForm()

	//log.Debug().Msgf("len anpr: %d", form.Value)

	if err != nil {
		log.Warn().Str("Camera IP", clientIP).Msg("Received message without multipart file")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Multipart form required",
			"code":  12,
		})
		return
	}
	//log.Debug().Msgf("len anpr: %d", len(form.File["anpr.xml"]))
	// log.Debug().Msgf("len vehiclePicture: %d", len(form.File["vehiclePicture.xml"]))

	////////////////////////////////	anpr.xml	//////// Get specific file   form.File["anpr.xml"]
	anprfil := form.File["anpr.xml"]
	if anprfil == nil {
		if config.Configvar.App.ExtraLog == "yes" {
			log.Warn().Str("Camera IP", clientIP).Msgf("anpr.xml not present file List: %v", form.File)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Multipart form required",
			"code":  12,
		})
		return
	}

	myfile, err := anprfil[0].Open()
	if myfile == nil {
		log.Error().Msgf("Error read anpr.xml: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Files not uploaded successfully",
			"code":    8,
		})
		return
	}

	file_as_string, err := ReadFromFileMulti(myfile)
	log.Debug().Msgf("------------------- # Reading ANPR File # -----------------")

	if err != nil {
		log.Err(err).Msgf("Error OpenFIle: %v", file_as_string)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  12,
		})
		return
	}

	DataCapt, err = Get_lpr_data(string(file_as_string))
	if err != nil {
		log.Error().Msgf("Error getting Data from eventNotification: %v", err)
	}

	if config.Configvar.App.SaveXml == "true" {
		go SaveFile(anprfil[0], "anpr_"+"_"+DataCapt.LicensePlate+".xml", clientIP)
		log.Debug().Msgf("Saving File Locaally assoc with Licence Plate Number:  %v", DataCapt.LicensePlate)
	}

	/*licensePlatePicture := form.File["licensePlatePicture.jpg"]
	vehiclePicture := form.File["vehiclePicture.jpg"]
	var Img1, Img2 string

	if licensePlatePicture != nil {
		Img1 = ManageFilePicture(licensePlatePicture[0], clientIP, "licensePlate", formattedTime)
	} else {
		Img1 = ""
	}

	if vehiclePicture != nil {
		Img2 = ManageFilePicture(vehiclePicture[0], clientIP, "vehicle", formattedTime)
	} else {
		Img2 = ""
	}

	 	Cdet := db.CarDetail{
	   		//CamBody: form.Value,
	   		Image1:  Img1,
	   		Image2:  Img2,
	   	}
	*/
	//log.Debug().Interface("* - * - *- * DATA *- *- *- *- ", Cdet).Send()

	go DataCapt.ProcessPresentCar(formattedTime, DataCapt)

	jsonData, err := json.Marshal(DataCapt)
	if err != nil {
		log.Error().Msgf("Error marshalling eventNotification to JSON: %v", err)
	}

	log.Debug().Msgf("------------------- # Converting Data As JSON # -----------------")
	log.Info().Msgf("ANPR event data as JSON: %s", string(jsonData))

	//go SaveDataInCamBody(jsonData) // do it later

	// camip / LPN / direction / confidance / country/ datetimenow / picturess ( dont log it) / datetimenow

	//c.JSON(http.StatusOK, gin.H{"message": "Files uploaded successfully"})
	log.Info().Msgf("Data Readed for the Licence plate %v Successfully at : %v", DataCapt.LicensePlate, DataCapt.CaptureTime)
	c.JSON(http.StatusOK, gin.H{
		//"time":    formattedTime,
		//"file":    file_as_string,
		"message": "Files uploaded successfully",
		"code":    8,
	})
}
