package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"fyc/functions"
)

type Camera struct {
	bun.BaseModel `json:"-" bun:"table:camera"`
	ID            int                    `bun:"id,autoincrement" json:"-"`
	CamID         int                    `bun:"cam_id,pk" json:"cam_id" binding:"required"`
	CamName       string                 `bun:"cam_name" json:"cam_name" binding:"required"`
	CamType       string                 `bun:"cam_type" json:"cam_type" binding:"required"`
	CamIP         string                 `bun:"cam_ip" json:"cam_ip" binding:"required"`
	CamPORT       int                    `bun:"cam_port" json:"cam_port" binding:"required"`
	CamUser       string                 `bun:"cam_user" json:"cam_user" binding:"required"`
	CamPass       string                 `bun:"cam_password" json:"cam_password" binding:"required"`
	ZoneIdIn      *int                   `bun:"zone_in_id" json:"zone_in_id" binding:"required"`
	ZoneIdOut     *int                   `bun:"zone_out_id" json:"zone_out_id" binding:"required"`
	Direction     string                 `bun:"direction" json:"direction" binding:"required"`
	IsEnabled     bool                   `bun:"is_enabled,type:bool" json:"is_enabled"`
	IsDeleted     bool                   `bun:"is_deleted,type:bool" json:"-"`
	LastUpdated   string                 `bun:"last_update,type:timestamp" json:"-"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" swaggertype:"object"`
}

type ResponseCamera struct {
	bun.BaseModel `json:"-" bun:"table:camera"`
	ID            int    `bun:"id" json:"-"`
	CamID         int    `bun:"cam_id" json:"cam_id"`
	CamName       string `bun:"cam_name" json:"cam_name"`
	CamType       string `bun:"cam_type" json:"cam_type"`
	CamIP         string `bun:"cam_ip" json:"cam_ip"`
	CamPORT       int    `bun:"cam_port" json:"cam_port" `
	CamUser       string `bun:"cam_user"  json:"cam_user"`
	CamPass       string `bun:"cam_password" json:"cam_password"`
	ZoneIdIn      *int   `bun:"zone_in_id"  json:"zone_in_id"`
	ZoneIdOut     *int   `bun:"zone_out_id" json:"zone_out_id"`
	Direction     string `bun:"direction" json:"direction"`
	IsEnabled     bool   `bun:"is_enabled" json:"-"`
	LastUpdated   string `bun:"last_update" json:"last_update"`
	IsDeleted     bool   `bun:"is_deleted" json:"-"`
}

type CameraNoBind struct {
	bun.BaseModel `json:"-" bun:"table:camera"`
	ID            int                    `bun:"id,autoincrement" json:"-"`
	CamID         int                    `bun:"cam_id,pk" json:"cam_id"`
	CamName       string                 `bun:"cam_name" json:"cam_name"`
	CamType       string                 `bun:"cam_type" json:"cam_type"`
	CamIP         string                 `bun:"cam_ip" json:"cam_ip" `
	CamPORT       int                    `bun:"cam_port" json:"cam_port"`
	CamUser       string                 `bun:"cam_user" json:"cam_user" `
	CamPass       string                 `bun:"cam_password" json:"cam_password"`
	ZoneIdIn      *int                   `bun:"zone_in_id" json:"zone_in_id"`
	ZoneIdOut     *int                   `bun:"zone_out_id" json:"zone_out_id"`
	Direction     string                 `bun:"direction" json:"direction" `
	IsEnabled     *bool                  `bun:"is_enabled,type:bool" json:"is_enabled"`
	IsDeleted     *bool                  `bun:"is_deleted,type:bool" json:"-"`
	LastUpdated   string                 `bun:"last_update,type:timestamp" json:"-"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" swaggertype:"object"`
}

// Get all camera Data
func GetDataCamera(ctx context.Context) ([]Camera, error) {
	var camData []Camera
	err := Db_GlobalVar.NewSelect().Model(&camData).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all cameras : %w", err)
	}

	for i := range camData {
		camData[i].LastUpdated, _ = functions.ParseTimeData(camData[i].LastUpdated)
	}
	return camData, nil
}

// Get all camera with extra data
func GetAllCameraExtra(ctx context.Context) ([]Camera, error) {
	var camera []Camera
	err := Db_GlobalVar.NewSelect().Model(&camera).
		Where("is_deleted = ?", false).
		Column().
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Camera with Extra Data: %w", err)
	}
	for i := range camera {
		camera[i].LastUpdated, _ = functions.ParseTimeData(camera[i].LastUpdated)
	}
	return camera, nil
}

func GetCameraByStatus(ctx context.Context, status string) ([]ResponseCamera, error) {
	var camera []ResponseCamera

	query := Db_GlobalVar.NewSelect().
		Model(&camera).
		Where("is_deleted = ?", false)

	switch status {
	case "enabled":
		query.Where("is_enabled = ?", true)
	case "disabled":
		query.Where("is_enabled = ?", false)
	case "":
	default:
		return nil, fmt.Errorf("invalid status: %s, expected 'enabled', 'disabled', or 'all'", status)
	}

	// Execute the query
	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("error getting client -- status: %s, err: %w", status, err)
	}

	// Parse `LastUpdated` for each zone
	for i := range camera {
		camera[i].LastUpdated, _ = functions.ParseTimeData(camera[i].LastUpdated)
	}

	return camera, nil
}

func GetCameraByStatusZone(ctx context.Context, status string, zone int) ([]ResponseCamera, error) {
	var camera []ResponseCamera

	query := Db_GlobalVar.NewSelect().
		Model(&camera).
		Where("zone_in_id = ?", zone).
		Where("is_deleted = ?", false)

	switch status {
	case "enabled":
		query.Where("is_enabled = ?", true)
	case "disabled":
		query.Where("is_enabled = ?", false)
	case "":
	default:
		return nil, fmt.Errorf("invalid status: %s, expected 'enabled', 'disabled', or 'all'", status)
	}

	// Execute the query
	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("error getting client -- status: %s, err: %w", status, err)
	}

	// Parse `LastUpdated` for each zone
	for i := range camera {
		camera[i].LastUpdated, _ = functions.ParseTimeData(camera[i].LastUpdated)
	}

	return camera, nil
}

// Get all camera
func GetAllCamera(ctx context.Context) ([]ResponseCamera, error) {
	var cam []ResponseCamera
	err := Db_GlobalVar.NewSelect().
		Model(&cam).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all cameras : %w", err)
	}

	for i := range cam {
		cam[i].LastUpdated, _ = functions.ParseTimeData(cam[i].LastUpdated)
	}
	return cam, nil
}

// Get all camera
func GetCameras(ctx context.Context) ([]ResponseCamera, error) {
	var cam []ResponseCamera
	err := Db_GlobalVar.NewSelect().
		Model(&cam).
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all cameras : %w", err)
	}

	return cam, nil
}

// Get camera by ID
func GetCameraByIDExtra(ctx context.Context, id int) (*Camera, error) {
	var cam Camera
	err := Db_GlobalVar.NewSelect().
		Model(&cam).
		Where("cam_id = ?", id).
		Where("is_deleted = ?", false).
		//Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting camera by id %d: %w", id, err)
	}

	cam.LastUpdated, _ = functions.ParseTimeData(cam.LastUpdated)
	return &cam, nil
}

func GetCameraByID(ctx context.Context, id int) (*ResponseCamera, error) {
	var cam ResponseCamera
	err := Db_GlobalVar.NewSelect().
		Model(&cam).
		Where("cam_id = ?", id).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting camera by id %d: %w", id, err)
	}
	cam.LastUpdated, _ = functions.ParseTimeData(cam.LastUpdated)
	return &cam, nil
}

func GetCamByID(ctx context.Context, id int) (*ResponseCamera, error) {
	var cam ResponseCamera
	err := Db_GlobalVar.NewSelect().
		Model(&cam).
		Where("cam_id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting camera by id %d: %w", id, err)
	}
	cam.LastUpdated, _ = functions.ParseTimeData(cam.LastUpdated)
	return &cam, nil
}
func GetCameraListEnabled(ctx context.Context) ([]ResponseCamera, error) {
	var camera []ResponseCamera
	err := Db_GlobalVar.NewSelect().
		Model(&camera).
		Where("is_enabled = ?", true).
		Scan(ctx, &camera)

	if err != nil {
		return nil, fmt.Errorf("error getting Enabled camera List: %w", err)
	}

	return camera, nil
}

func GetCameraListEnabledExtra(ctx context.Context) ([]Camera, error) {
	var camera []Camera
	err := Db_GlobalVar.NewSelect().
		Model(&camera).
		Where("is_enabled = ?", true).
		Scan(ctx, &camera)

	if err != nil {
		return nil, fmt.Errorf("error getting Enabled camera List: %w", err)
	}

	return camera, nil
}

func GetCameraEnabledByID(ctx context.Context, id int) (*Camera, error) {
	var cam Camera
	err := Db_GlobalVar.NewSelect().
		Model(&cam).
		Column("is_enabled").
		Where("cam_id = ?", id).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx, &cam)

	if err != nil {
		return nil, fmt.Errorf("error getting camera by id: %w", err)
	}

	return &cam, nil
}

// CreateCamera creates a new camera if the ID is not found, and checks if the database contains data.
func CreateCamera(ctx context.Context, newcam *Camera) error {
	log.Debug().Str("Camera Added AT:", functions.GetFormatedLocalTime())
	newcam.LastUpdated = functions.GetFormatedLocalTime()
	newcam.IsDeleted = false
	newcam.IsEnabled = true
	newcam.Direction = strings.ToLower(newcam.Direction)

	_, err := Db_GlobalVar.NewInsert().Model(newcam).Returning("cam_id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating a new camera: %w", err)
	}

	LoadCameralist()
	CamStartup()
	log.Debug().Msgf("New camera added with ID: %d", newcam.CamID)
	return nil
}

// Update a camera by ID
func UpdateCamera(ctx context.Context, cam_id int, updates *CameraNoBind) (int64, error) {
	log.Debug().Str("Camera Updated AT:", functions.GetFormatedLocalTime())
	updates.LastUpdated = functions.GetFormatedLocalTime()
	updates.Direction = strings.ToLower(updates.Direction)

	res, err := Db_GlobalVar.NewUpdate().
		Model(updates).
		Where("cam_id = ?", cam_id).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating camera with id %d: %w", cam_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated camera with ID: %d, rows affected: %d", cam_id, rowsAffected)
	CamStartup()
	LoadCameralist()
	return rowsAffected, nil
}

// Delete a camera by ID (soft delete: sets is_deleted to true)
func DeleteCamera(ctx context.Context, Cam_id int) (int64, error) {
	log.Debug().Str("Camera Added AT:", functions.GetFormatedLocalTime())

	res, err := Db_GlobalVar.NewUpdate().
		Model(&Camera{}).
		Where("cam_id = ?", Cam_id).
		//Where("is_enabled = ?", true).
		Set("is_deleted = ?", true).
		Set("is_enabled = ?", false).
		Set("last_update = ?", functions.GetFormatedLocalTime()).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error deleting Camera with id %d: %w", Cam_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Soft-deleted Camera with ID: %d, rows affected: %d", Cam_id, rowsAffected)
	LoadCameralist()
	CamStartup()
	return rowsAffected, nil
}
