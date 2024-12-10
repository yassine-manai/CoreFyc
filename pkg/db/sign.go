package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"fyc/functions"
)

type Sign struct {
	bun.BaseModel `json:"-" bun:"table:sign"`
	ID            int                    `bun:"id,autoincrement,pk" json:"-"`
	SignID        int                    `bun:"sign_id,pk" binding:"required" json:"sign_id"`
	SignName      map[string]interface{} `bun:"sign_name,type:jsonb" binding:"required" json:"sign_name" swaggertype:"object"`
	SignUserName  string                 `bun:"sign_username" binding:"required" json:"sign_username"`
	SignPassword  string                 `bun:"sign_password" binding:"required" json:"sign_password"`
	SignType      string                 `bun:"sign_type" binding:"required" json:"sign_type"`
	SignIP        string                 `bun:"sign_ip" binding:"required" json:"sign_ip"`
	SignPort      int                    `bun:"sign_port" binding:"required" json:"sign_port"`
	ZoneID        int                    `bun:"zone_id" binding:"required" json:"zone_id"`
	IsEnabled     bool                   `bun:"is_enabled,type:bool" json:"-"`
	IsDeleted     bool                   `bun:"is_deleted,type:bool" json:"-"`
	LastUpdated   string                 `bun:"last_update,type:timestamp" json:"-"`
}

type SignNoBind struct {
	bun.BaseModel `json:"-" bun:"table:sign"`
	ID            int                    `bun:"id,autoincrement" json:"-"`
	SignID        int                    `bun:"sign_id,pk" json:"sign_id"`
	SignName      map[string]interface{} `bun:"sign_name" json:"sign_name" swaggertype:"object"`
	SignUserName  string                 `bun:"sign_username" json:"sign_username"`
	SignPassword  string                 `bun:"sign_password"  json:"sign_password"`
	SignType      string                 `bun:"sign_type"  json:"sign_type"`
	SignIP        string                 `bun:"sign_ip" json:"sign_ip"`
	SignPort      int                    `bun:"sign_port"  json:"sign_port"`
	ZoneID        int                    `bun:"zone_id"  json:"zone_id"`
	IsEnabled     *bool                  `bun:"is_enabled,type:bool" json:"is_enabled"`
	IsDeleted     *bool                  `bun:"is_deleted,type:bool" json:"-"`
	LastUpdated   string                 `bun:"last_update,type:timestamp" json:"-"`
}

type SignResp struct {
	bun.BaseModel `json:"-" bun:"table:sign"`
	ID            int                    `bun:"id,autoincrement" json:"-"`
	SignID        int                    `bun:"sign_id,pk" json:"sign_id"`
	SignName      map[string]interface{} `bun:"sign_name" json:"sign_name" swaggertype:"object"`
	SignUserName  string                 `bun:"sign_username" json:"sign_username"`
	SignPassword  string                 `bun:"sign_password"  json:"sign_password"`
	SignType      string                 `bun:"sign_type" json:"sign_type"`
	SignIP        string                 `bun:"sign_ip" json:"sign_ip"`
	SignPort      int                    `bun:"sign_port" json:"sign_port"`
	ZoneID        int                    `bun:"zone_id"  json:"zone_id"`
	IsEnabled     bool                   `bun:"is_enabled,type:bool" json:"is_enabled"`
	IsDeleted     bool                   `bun:"is_deleted,type:bool" json:"-"`
	LastUpdated   string                 `bun:"last_update,type:timestamp" json:"last_update"`
}

func CreateSign(ctx context.Context, sign *Sign) error {
	//log.Debug().Str("Sign Added AT:", functions.GetFormatedLocalTime())
	sign.LastUpdated = functions.GetFormatedLocalTime()
	sign.IsDeleted = false
	sign.IsEnabled = true

	// Initialize a new map to hold normalized names
	normalizedNames := make(map[string]interface{}, len(sign.SignName))
	for lang, name := range sign.SignName {
		if str, ok := name.(string); ok && str != "" {
			// Use the original name while lowering the key
			normalizedNames[strings.ToLower(lang)] = str
		}
	}
	sign.SignName = normalizedNames
	log.Debug().Int("Sign ID", sign.SignID).Str("Time", functions.GetFormatedLocalTime()).Msg("Adding sign ")
	_, err := Db_GlobalVar.NewInsert().Model(sign).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding sign: %w", err)
	}
	LoadSignlist()
	return nil
}

func GetSignById(ctx context.Context, signID int) (*SignResp, error) {
	var sign SignResp

	err := Db_GlobalVar.NewSelect().Model(&sign).
		Where("sign_id = ?", signID).
		//Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("sign with SignID %d not found", signID)
		}
		return nil, fmt.Errorf("error retrieving sign with SignID %d: %w", signID, err)
	}

	sign.LastUpdated, _ = functions.ParseTimeData(sign.LastUpdated)
	return &sign, nil
}

func GetSignByZoneId(ctx context.Context, zoneID int) (*SignResp, error) {
	var sign SignResp

	err := Db_GlobalVar.NewSelect().Model(&sign).
		Where("zone_id = ?", zoneID).
		//Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("sign with Zone ID %d not found", zoneID)
		}
		return nil, fmt.Errorf("error retrieving sign with Zone ID %d: %w", zoneID, err)
	}
	sign.LastUpdated, _ = functions.ParseTimeData(sign.LastUpdated)
	return &sign, nil
}

func GetSigns(ctx context.Context) ([]SignResp, error) {
	var signs []SignResp
	err := Db_GlobalVar.NewSelect().Model(&signs).
		Where("is_deleted = ?", false).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all signs: %w", err)
	}

	return signs, nil
}

func GetSignByStatus(ctx context.Context, status string) ([]SignResp, error) {
	var signStat []SignResp

	query := Db_GlobalVar.NewSelect().
		Model(&signStat).
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
	for i := range signStat {
		signStat[i].LastUpdated, _ = functions.ParseTimeData(signStat[i].LastUpdated)
	}

	return signStat, nil
}

func GetZoneByName(ctx context.Context, name string) (*Zone, error) {
	var zone Zone
	query := Db_GlobalVar.NewSelect().
		Model(&zone).
		Where("is_deleted = ?", false).
		Where("name->>'en' = ? OR name->>'ar' = ?", name, name)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == "no rows found" {
			return nil, fmt.Errorf("zone not found")
		}
		return nil, fmt.Errorf("error fetching zone: %w", err)
	}

	return &zone, nil
}

func GetSignByStatusZone(ctx context.Context, status string, zone_id int) ([]SignResp, error) {
	var signStat []SignResp

	query := Db_GlobalVar.NewSelect().
		Model(&signStat).
		Where("zone_id = ?", zone_id).
		Where("is_deleted = ?", false)

	switch status {
	case "enabled":
		query.Where("is_enabled = ?", true)
	case "disabled":
		query.Where("is_enabled = ?", false)
	case "all":
	default:
		return nil, fmt.Errorf("invalid status: %s, expected 'enabled', 'disabled', or 'all'", status)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("error getting signs by status and zone: %w", err)
	}

	for i := range signStat {
		signStat[i].LastUpdated, _ = functions.ParseTimeData(signStat[i].LastUpdated)
	}

	return signStat, nil
}

func GetAllSigns(ctx context.Context) ([]SignResp, error) {
	var signs []SignResp
	err := Db_GlobalVar.NewSelect().Model(&signs).
		Where("is_deleted = ?", false).
		//Where("is_enabled = ?", true).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all signs: %w", err)
	}

	for i := range signs {
		signs[i].LastUpdated, _ = functions.ParseTimeData(signs[i].LastUpdated)
	}
	return signs, nil
}

func UpdateSign(ctx context.Context, signID int, updatedSign SignNoBind) (int64, error) {
	updatedSign.LastUpdated = functions.GetFormatedLocalTime()

	// convert sign name key to lowercase
	normalizedNames := make(map[string]interface{}, len(updatedSign.SignName))
	for lang, name := range updatedSign.SignName {
		if str, ok := name.(string); ok && str != "" {
			normalizedNames[strings.ToLower(lang)] = str
		}
	}
	updatedSign.SignName = normalizedNames

	log.Debug().Int("Sign ID", signID).Str("Time", functions.GetFormatedLocalTime()).Msg("Updating sign ")
	result, err := Db_GlobalVar.NewUpdate().
		Model(&updatedSign).
		//Where("is_deleted = ?", false).
		//Where("is_enabled = ?", true).
		Where("sign_id = ?", signID).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating sign with SignID %d: %w", signID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	LoadSignlist()
	return rowsAffected, nil
}

func DeleteSign(ctx context.Context, signID int) (int64, error) {
	log.Debug().Msgf("Deleting Sign with SignID: %d", signID)

	log.Debug().Str("Sign Deleted AT:", functions.GetFormatedLocalTime())

	result, err := Db_GlobalVar.NewUpdate().
		Model(&Sign{}).
		Where("sign_id = ?", signID).
		//Where("is_enabled = ?", true).
		Set("is_deleted = ?", true).
		Set("is_enabled = ?", false).
		Set("last_update = ?", functions.GetFormatedLocalTime()).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error deleting sign with SignID %d: %w", signID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	LoadSignlist()
	return rowsAffected, nil
}
