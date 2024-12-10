package hikvision

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"fyc/pkg/db"
)

//GetPresentFound

func ReadFromFileMulti(file io.Reader) (string, error) {
	// Create a buffer to store the file content
	var buf bytes.Buffer

	// Read the file content into the buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		// Log and return the error if any issues occurred while reading
		log.Printf("Error reading from file: %v", err)
		return "", err
	}

	// Convert the buffer to a string and return it
	return buf.String(), nil
}

// //////////////////////////////////////////////////////////////

func ReadFromFile(file *multipart.FileHeader) ([]byte, error) {
	// Open the file from the multipart.FileHeader
	src, err := file.Open()
	if err != nil {
		return []byte{}, err
	}
	defer src.Close()

	// Read the content of the file
	content, err := io.ReadAll(src)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

func Get_lpr_data(anprXml string) (Capture, error) {
	var eventNotification EventNotificationAlert
	var capture Capture
	now := time.Now().UTC()
	formattedTime := now.Format("2006-01-02 15:04:05")

	log.Debug().Msgf("------------------- # Getting LPR Data # -----------------")

	// Unmarshal XML to struct
	err := xml.Unmarshal([]byte(anprXml), &eventNotification)
	if err != nil {
		log.Error().Msgf("Error unmarshalling XML: %v", err)
		return capture, err
	}

	// Populate the capture object from eventNotification
	capture.State = eventNotification.ANPR.Country
	capture.LicensePlate = eventNotification.ANPR.LicensePlate
	capture.Direction = eventNotification.ANPR.Direction
	capture.Confidence = eventNotification.ANPR.ConfidenceLevel
	capture.CamIP = eventNotification.IpAddress
	capture.CaptureTime = formattedTime

	log.Info().Msgf("ANPR event data after refactor: %v", capture)

	return capture, nil
}

func SaveDataInCamBody(data []byte) (string, error) {
	ctx := context.Background()
	var CD db.CarDetail

	log.Info().Msgf("Camera Data to save in DB: %v", string(data))

	err := json.Unmarshal(data, &CD)
	if err != nil {
		log.Err(err).Msgf("Error unmarshalling JSON data: %v", string(data))
		return "Error", err
	}

	if CD.CamBody == nil {
		CD.CamBody = make(map[string]interface{})
	}

	CD.ID = 1
	CD.CamBody["camBody"] = data
	CD.Image1 = ""
	CD.Image2 = ""

	err = db.CreateCarDetail(ctx, &CD)
	if err != nil {
		log.Err(err).Msgf("Error saving data to db: %v", string(data))
		return "Error", err
	}

	log.Debug().Msgf("Data saved successfully: %v", string(data))
	return "Success", nil
}

func isCamExist(camList map[string]db.CameraStarter, capture string) (bool, *db.CameraStarter) {
	if cameraStr, exists := camList[capture]; exists {
		return true, &cameraStr
	}
	return false, nil
}

func SaveFile(file *multipart.FileHeader, fileName string, folder string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	destDir := "./logs/cameras/" + folder
	log.Debug().Msgf("Create Directory: %v", destDir)

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return err
	}

	dst, err := os.Create(destDir + "/" + fileName)
	if err != nil {
		fmt.Println(err)
		log.Err(err).Msgf("Error Create Directory")

		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func JsonData(jsonData string) {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Error().Msgf("Error unmarshalling JSON: %v", err)
	}

	for key, value := range data {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func ManageFilePicture(FileImage *multipart.FileHeader, clientIP string, filename string, formattedTime string) string {
	myfile, err := FileImage.Open()
	if err != nil {
		log.Error().Msgf("Error read %s: %v", filename, err)
		return ""
	}

	imageFile, err := ReadFromFileMulti(myfile)
	if err != nil {
		log.Err(err).Str("CameraIp", clientIP).Msg(" Error reading FileImage")
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return ""
	}
	//log.Debug().Bool("ResizePicture", global.AppConfig.ResizePicture).Int("MaxWidth", global.AppConfig.MaxWidth).Send()

	/* 	if global.AppConfig.SaveImage {
		go saveFile(FileImage, filename+"_"+formattedTime+".jpg", clientIP)
	} */

	return imageFile
}

/* func ResizeImage(imageBytes []byte, maxWidth, maxHeight uint) ([]byte, error) {
	// Decode the image from byte slice
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	// Calculate new dimensions while maintaining aspect ratio
	var newWidth, newHeight uint
	if img.Bounds().Dx() > img.Bounds().Dy() {
		// Landscape orientation
		newWidth = maxWidth
		newHeight = uint(float64(img.Bounds().Dy()) / float64(img.Bounds().Dx()) * float64(maxWidth))
	} else {
		// Portrait or square orientation
		newWidth = uint(float64(img.Bounds().Dx()) / float64(img.Bounds().Dy()) * float64(maxHeight))
		newHeight = maxHeight
	}

	// Resize the image
	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// Encode the resized image to bytes
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, resizedImg, nil); err != nil {
		return nil, err
	}

	// Return the resized image bytes
	return buffer.Bytes(), nil
} */
