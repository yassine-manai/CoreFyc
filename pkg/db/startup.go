package db

import (
	"context"
	"fmt"

	_ "github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

// var CarParkList []int
var Zonelist []int
var AllZonelist []int
var CameraList []int
var SignList []int
var ClientListAPI []string
var SettingsList = make(map[int]string)
var Token string = ""
var Db_GlobalVar *bun.DB
var ClientList = make(map[string]string)

type ClientDetails struct {
	ClientID        string
	ClientSecret    string
	ClientGrantType string
	ClientActive    bool
	FuzzyLogic      bool
}

func LoadSignlist() {
	//Zonelist = nil
	//log.Debug().Msgf("Prepare Sign List ...")
	ctx := context.Background()
	SignService, _ := GetSigns(ctx)
	for _, v := range SignService {
		SignList = append(SignList, v.SignID)
	}
}

/*
	 func LoadCarparklist() {
		log.Debug().Msgf("Prepare CarPark List \n")
		ctx := context.Background()

		CarService, _ := GetAllCarparks(ctx)
		for _, v := range CarService {
			CarParkList = append(CarParkList, v.ID)
		}
	}
*/

func LoadCameralist() {
	//CameraList = nil
	//log.Debug().Msgf("Camera List Loading ...")
	ctx := context.Background()

	CameraService, _ := GetCameras(ctx)
	for _, v := range CameraService {
		CameraList = append(CameraList, v.CamID)
	}
}

func LoadClientsApi() {
	ctx := context.Background()

	apikey, err := GetAllClientCred(ctx)
	if err != nil {
		fmt.Println("Error fetching client from DB:", err)
		return
	}
	for _, row := range apikey {
		ClientList[row.ClientID] = row.ClientSecret
	}
}

func LoadClientlist() {
	//CameraList = nil
	//log.Debug().Msgf("ClientAPI List Loading ...")
	ctx := context.Background()

	ClientService, _ := GetAllClientCred(ctx)
	for _, v := range ClientService {
		ClientListAPI = append(ClientListAPI, v.ClientID)
	}
}

var ClientDataList = make(map[string]ClientDetails)

func LoadClientDataList() {
	ctx := context.Background()

	apikey, err := GetAllClientDatas(ctx)
	if err != nil {
		fmt.Println("Error fetching client from DB:", err)
		return
	}

	for _, row := range apikey {
		ClientDataList[row.ClientID] = ClientDetails{
			ClientID:        row.ClientID,
			ClientSecret:    row.ClientSecret,
			ClientGrantType: row.GrantType,
			ClientActive:    row.IsEnabled,
			FuzzyLogic:      *row.FuzzyLogic,
		}
	}
}

func LoadzoneList() {
	ctx := context.Background()
	ZoneService, _ := GetZones(ctx)
	for _, v := range ZoneService {
		AllZonelist = append(AllZonelist, *v.ZoneID)
	}
}

func LoadAllZonelist() {
	ctx := context.Background()
	ZoneSer, _ := GetAllZone(ctx)
	for _, v := range ZoneSer {
		Zonelist = append(Zonelist, *v.ZoneID)
	}
}
