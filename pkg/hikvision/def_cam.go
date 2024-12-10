package hikvision

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"fyc/functions"
	"fyc/pkg/counting"
	"fyc/pkg/db"
)

type EventNotificationAlert struct {
	ChannelName string `xml:"channelName" json:"channelName,omitempty"`
	MacAddress  string `xml:"macAddress" json:"macAddress,omitempty"`
	IpAddress   string `xml:"ipAddress" json:"ipAddress,omitempty"`
	PicNum      int    `xml:"picNum" json:"picNum,omitempty"`
	PortNo      string `xml:"portNo" json:"portNo,omitempty"`
	Protocol    string `xml:"protocol" json:"protocol,omitempty"`
	ChannelID   string `xml:"channelID" json:"channelID,omitempty"`
	//DateTime             string   `xml:"dateTime" json:"-"`
	ActivePostCount      int      `xml:"activePostCount" json:"activePostCount,omitempty"`
	EventType            string   `xml:"eventType" json:"eventType,omitempty"`
	EventState           string   `xml:"eventState" json:"eventState,omitempty"`
	EventDescription     string   `xml:"eventDescription" json:"eventDescription,omitempty"`
	UUID                 string   `xml:"UUID" json:"UUID,omitempty"`
	IsDataRetransmission bool     `xml:"isDataRetransmission" json:"isDataRetransmission,omitempty"`
	ANPR                 ANPRInfo `xml:"ANPR" json:"ANPR,omitempty"`
}

type ANPRInfo struct {
	Country         string `xml:"country" json:"country"`
	Province        string `xml:"province" json:"province"`
	LicensePlate    string `xml:"licensePlate" json:"licensePlate,omitempty"`
	Line            int    `xml:"line" json:"line,omitempty"`
	Direction       string `xml:"direction" json:"direction,omitempty"` // can be "reverse" or forward or uknown
	ConfidenceLevel int    `xml:"confidenceLevel" json:"confidenceLevel,omitempty"`
	//PlateType         string    `xml:"plateType" json:"-"`
	//PlateColor        string    `xml:"plateColor" json:"-"`
	//LicenseBright     int       `xml:"licenseBright" json:"-"`
	//DangMark          string    `xml:"dangmark" json:"-"`
	//TwoWheelVehicle   string    `xml:"twoWheelVehicle" json:"-"`
	//ThreeWheelVehicle string    `xml:"threeWheelVehicle" json:"-"`
	//PlateCharBelieve  string    `xml:"plateCharBelieve" json:"-"`
	//SpeedLimit        int       `xml:"speedLimit" json:"-"`
	//Region            string    `xml:"region" json:"-"`
	//Area              string    `xml:"area" json:"-"`
	//PlateSize         string    `xml:"plateSize" json:"-"`
	//Category          string    `xml:"category" json:"-"`
	VehicleType string `xml:"vehicleType" json:"vehicleType,omitempty"`
	DetectDir   int    `xml:"detectDir" json:"detectDir"`
	//DetectType           int         `xml:"detectType" json:"detectType"`
	//AlarmDataType        int         `xml:"alarmDataType" json:"-"`
	VehicleInfo          VehicleInfo   `xml:"vehicleInfo" json:"vehicleInfo,omitempty"`
	PictureInfoList      []PictureInfo `xml:"pictureInfo" json:"pictureInfoList,omitempty"`
	OriginalLicensePlate string        `xml:"originalLicensePlate" json:"originalLicensePlate,omitempty"`
	//CRIndex              int         `xml:"CRIndex" json:"CRIndex"`
	//VehicleListName string `xml:"vehicleListName" json:"vehicleListName"`
}

type VehicleInfo struct {
	Index              int    `xml:"index" json:"index,omitempty"`
	ColorDepth         int    `xml:"colorDepth" json:"colorDepth,omitempty"`
	Color              string `xml:"color" json:"color,omitempty"`
	Length             int    `xml:"length" json:"length,omitempty"`
	VehicleLogoRecog   int    `xml:"vehicleLogoRecog" json:"vehicleLogoRecog,omitempty"`
	VehileSubLogoRecog int    `xml:"vehileSubLogoRecog" json:"vehileSubLogoRecog,omitempty"`
	VehileModel        int    `xml:"vehileModel" json:"vehileModel,omitempty"`
}

type PictureInfo struct {
	FileName string `xml:"fileName" json:"fileName,omitempty"`
	Type     string `xml:"type" json:"type,omitempty"`
	DataType int    `xml:"dataType" json:"dataType,omitempty"`
	AbsTime  string `xml:"absTime" json:"absTime,omitempty"`
	PId      string `xml:"pId" json:"pId,omitempty"`
}

type Capture struct {
	State        string `json:"country"`
	LicensePlate string `json:"licensePlate,omitempty"`
	Direction    string `json:"direction,omitempty"` // can be "reverse" or "forward" or "unknown"
	Confidence   int    `json:"confidenceLevel,omitempty"`
	CamIP        string `json:"ipAddress,omitempty"`
	CaptureTime  string `json:"cap_time,omitempty"`
	//PortNo          string `xml:"portNo" json:"portNo,omitempty"`
	//PictureInfoList []PictureInfo `xml:"pictureInfo" json:"pictureInfoList,omitempty"`
}

func Proces_PrCar(captr Capture, camData db.CameraStarter) db.PresentCar {

	var direction = captr.Direction
	var currZone, lastZone *int
	ctx := context.Background()
	//var direction = camData.Direction

	switch captr.Direction {

	case "forward":
		log.Debug().Str("Direction", direction).Msg("/ / / Camera detect car  / / / ")
		currZone = &camData.ZoneIdIn
		lastZone = &camData.ZoneIdOut
		direction = "forward"

		capacity := counting.Decrease_Zone_Capacity(ctx, *currZone, captr.LicensePlate)
		counting.Sign_Data_Values(*currZone, "dec", fmt.Sprintf("%d", capacity))

	case "reverse":
		log.Debug().Str("Direction", direction).Msg("* * * Camera detect car  * * *  ")
		currZone = &camData.ZoneIdOut
		lastZone = &camData.ZoneIdIn
		direction = "reverse"

		capacity := counting.Increase_Zone_Capacity(ctx, *currZone, captr.LicensePlate)
		counting.Sign_Data_Values(*currZone, "inc", fmt.Sprintf("%d", capacity))

	default:
		log.Debug().Str("Direction", direction).Msg("- - -  Camera detect car  - - -  ")

		currZone = &camData.ZoneIdIn
		lastZone = &camData.ZoneIdOut
		direction = "unknown"
	}

	PCam := db.PresentCar{
		CameraID:        camData.CamID,
		CurrZoneID:      currZone,
		LastZoneID:      lastZone,
		LPN:             captr.LicensePlate,
		Confidence:      &captr.Confidence,
		TransactionDate: functions.GetFormatedLocalTime(),
		Direction:       direction,
	}

	return PCam
}

func (c *Capture) ProcessPresentCar(captTime string, dataCapture Capture) {
	ctx := context.Background()

	if exists, cameras := isCamExist(db.CamList, dataCapture.CamIP); exists {
		log.Debug().Str("CamIP", dataCapture.CamIP).Str("LPN", dataCapture.LicensePlate).Str("Direction", dataCapture.Direction).Msg("Camera Existed")

		// Process DATA CAPTURE
		ProcessCar := Proces_PrCar(dataCapture, *cameras)

		// Check if present car already exists or not
		exists, _ := db.GetPresentFound(ctx, ProcessCar.LPN)
		if exists {
			log.Info().Str("Licence Plate", ProcessCar.LPN).Str("CamIP", dataCapture.CamIP).Msg("Present car data already EXIST")

			rows_affected, err := db.UpdatePresentCarByLpn(ctx, ProcessCar.LPN, &ProcessCar)
			if rows_affected == 0 {
				log.Error().Str("Licence Plate", ProcessCar.LPN).Str("CamIP", dataCapture.CamIP).Int("ROWS AFF", int(rows_affected)).Msg("No rows affected")
			}

			if err != nil {
				log.Err(err).Str("LPN", ProcessCar.LPN).Str("Direction", ProcessCar.Direction).Str("CamIP", dataCapture.CamIP).Msgf("Error updating present car")
			} else {
				log.Debug().Str("Licence Plate", ProcessCar.LPN).Str("CamIP", dataCapture.CamIP).Msg("Present car data successfully updated")
			}

		} else {
			log.Info().Str("Licence Plate", ProcessCar.LPN).Msg("Present car data NOT EXIST Creating new one")

			if err := db.CreatePresentCar(ctx, &ProcessCar); err != nil {
				log.Error().Msgf("Error creating present car: %v", err)
			} else {
				log.Debug().Str("Licence Plate", ProcessCar.LPN).Str("Direction", ProcessCar.Direction).Msg("Present car data successfully created")
			}

		}
		c.ProcessHistory(ctx, ProcessCar)

	} else {
		log.Warn().Str("CamIP", dataCapture.CamIP).Str("LPN", dataCapture.LicensePlate).Str("Direction", dataCapture.Direction).Msg("Camera NOT Existed")
	}
}

func (c *Capture) ProcessHistory(ctx context.Context, car2add db.PresentCar) {
	// Process History
	log.Debug().Str("Licence Plate", car2add.LPN).Str("Direction", car2add.Direction).Msg("Adding Present Car to history")

	PresentCarFormatted := db.PresentCarHistory{
		CameraID:        car2add.CameraID,
		CurrZoneID:      car2add.CurrZoneID,
		LastZoneID:      car2add.LastZoneID,
		LPN:             car2add.LPN,
		Confidence:      car2add.Confidence,
		Direction:       car2add.Direction,
		TransactionDate: car2add.TransactionDate,
		Extra:           car2add.Extra,
		CarDetailsID:    car2add.CarDetailsID,
	}

	if err := db.CreatePresentCarHistory(ctx, &PresentCarFormatted); err != nil {
		log.Err(err).Str("Licence Plate", car2add.LPN).Msgf("Error creating History Car")
	} else {
		log.Info().Str("Licence Plate", car2add.LPN).Str("Direction", car2add.Direction).Msg("Present car data successfully Added to history")
	}
}
