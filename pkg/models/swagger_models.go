package models

////////// THIS FILE REPRESENT STRCTS FOR THE SWAGGER DATA MODELS //////////

// AddZone represents the data structure for a new zone
type AddZoneModel struct {
	ZoneID       int    `json:"zone_id"`
	Name         Name   `json:"name"`
	MaxCapacity  int    `json:"max_capacity"`
	FreeCapacity int    `json:"free_capacity"`
	Images       Images `json:"images"`
}

type UpdateZoneModel struct {
	ZoneID       int    `json:"-"`
	Name         Name   `json:"name"`
	MaxCapacity  int    `json:"max_capacity"`
	FreeCapacity int    `json:"free_capacity"`
	Images       Images `json:"images"`
	IsEnabled    bool   `json:"is_enabled"`
}

// Name represents the name of the zone in multiple languages
type Name struct {
	AR string `json:"ar" example:""`
	EN string `json:"en" example:""`
}

// Images represents the images in multiple languages
type Images struct {
	Ar ImageModel `json:"ar"`
	En ImageModel `json:"en"`
}

// Image represents image data with large and small versions
type ImageModel struct {
	ImageL string `json:"image_l" example:""`
	ImageS string `json:"image_s" example:""`
}

type ZoneNamesModel struct {
	ZoneID       int    `json:"zone_id"`
	Name         Name   `json:"name"`
	MaxCapacity  int    `json:"max_capacity"`
	FreeCapacity int    `json:"free_capacity"`
	Images       Images `json:"-"`
}

type ZoneDataModel struct {
	ZoneID       int    `json:"zone_id" `
	Name         Name   `json:"name"`
	MaxCapacity  int    `json:"max_capacity"`
	FreeCapacity int    `json:"free_capacity"`
	Images       Images `json:"-"`
	IsEnabled    bool   `json:"is_enabled"`
	LastUpdated  string `json:"last_update"`
}

type ZoneDataModel2 struct {
	ZoneID       int    `json:"zone_id" `
	Name         Name   `json:"name"`
	MaxCapacity  int    `json:"max_capacity"`
	FreeCapacity int    `json:"free_capacity"`
	Images       Images `json:"images"`
	IsEnabled    bool   `json:"is_enabled"`
	LastUpdated  string `json:"last_update"`
}

type AddSignModel struct {
	SignID       int    `json:"sign_id" example:"101"`
	SignName     Name   `json:"sign_name"`
	SignUserName string `json:"sign_username"`
	SignPassword string `json:"sign_password"`
	SignType     string `json:"sign_type"`
	SignIP       string `json:"sign_ip"`
	SignPort     int    `json:"sign_port"`
	ZoneID       int    `json:"zone_id"`
}

type UpdateSignModel struct {
	SignID       int    `json:"sign_id" example:"101"`
	SignName     Name   `json:"sign_name"`
	SignUserName string `json:"sign_username"`
	SignPassword string `json:"sign_password"`
	SignType     string `json:"sign_type"`
	SignIP       string `json:"sign_ip"`
	SignPort     int    `json:"sign_port"`
	ZoneID       int    `json:"zone_id"`
	IsEnabled    bool   `json:"is_enabled"`
	LastUpdated  string `json:"-"`
}
