package db

import (
	"context"

	"github.com/rs/zerolog/log"
)

type CameraStarter struct {
	CamID     int    `json:"cam_id"`
	CamIP     string `json:"cam_ip"`
	CamPORT   int    `json:"cam_port"`
	CamUser   string `json:"cam_user"`
	CamPass   string `json:"cam_password"`
	ZoneIdIn  int    `json:"zone_in_id"`
	ZoneIdOut int    `json:"zone_out_id"`
	Direction string `json:"direction"`
}

var CamList = make(map[string]CameraStarter)

func CamStartup() {
	ctx := context.Background()

	/* Camera.CamID = CameraInit.ID
	Camera.CamIP = CameraInit.CamIP
	Camera.CamPORT = CameraInit.CamPORT
	Camera.CamUser = CameraInit.CamUser
	Camera.CamPass = CameraInit.CamPass
	Camera.ZoneIdIn = *CameraInit.ZoneIdIn
	Camera.ZoneIdOut = *CameraInit.ZoneIdOut
	Camera.Direction = CameraInit.Direction */

	Cam, err := GetAllCamera(ctx)
	if err != nil {
		log.Err(err).Msg("Error Getting Cameras")
		return
	}

	for _, camera := range Cam {
		CamList[camera.CamIP] = CameraStarter{
			CamID:     camera.CamID,
			CamIP:     camera.CamIP,
			CamPORT:   camera.CamPORT,
			CamUser:   camera.CamUser,
			CamPass:   camera.CamPass,
			ZoneIdIn:  *camera.ZoneIdIn,
			ZoneIdOut: *camera.ZoneIdOut,
			Direction: camera.Direction,
		}
	}
}
