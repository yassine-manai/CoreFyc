package db

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Settings struct {
	bun.BaseModel      `json:"-" bun:"table:settings"`
	CarParkID          int                    `bun:"carpark_id,pk" json:"carpark_id" binding:"required"`
	CarParkName        map[string]interface{} `bun:"carpark_name,type:jsonb" binding:"required" json:"carpark_name" swaggertype:"object"`
	AppLogo            string                 `bun:"app_logo,type:bytea" binding:"required" json:"app_logo"`
	DefaultLang        string                 `bun:"default_lang" binding:"required" json:"default_lang" `
	TimeOutScreenKisok int                    `bun:"timeout_screenKiosk" binding:"required" json:"timeout_screenKiosk"`
	PkaImageSize       string                 `bun:"pka_image_size" binding:"required" json:"pka_image_size"`
	FycCleanCron       int                    `bun:"fyc_clean_cron" binding:"required" json:"fyc_clean_cron"`
	IsFycEnabled       bool                   `bun:"is_fyc_enabled,type:bool" json:"is_fyc_enabled"`
	CountingCleanCron  int                    `bun:"counting_clean_cron" binding:"required" json:"counting_clean_cron"`
	IsCountingEnabled  bool                   `bun:"is_counting_enabled,type:bool" json:"is_counting_enabled"`
	TC                 string                 `bun:"tc" json:"tc"`
}

type SettingsNoBind struct {
	bun.BaseModel      `json:"-" bun:"table:settings"`
	CarParkID          int                    `bun:"carpark_id" json:"carpark_id"`
	CarParkName        map[string]interface{} `bun:"carpark_name,type:jsonb" json:"carpark_name" swaggertype:"object"`
	AppLogo            string                 `bun:"app_logo,type:bytea" json:"app_logo"`
	DefaultLang        string                 `bun:"default_lang" json:"default_lang" `
	TimeOutScreenKisok int                    `bun:"timeout_screenKiosk" json:"timeout_screenKiosk"`
	PkaImageSize       string                 `bun:"pka_image_size" json:"pka_image_size"`
	FycCleanCron       int                    `bun:"fyc_clean_cron"  json:"fyc_clean_cron"`
	IsFycEnabled       *bool                  `bun:"is_fyc_enabled,type:bool" json:"is_fyc_enabled"`
	CountingCleanCron  int                    `bun:"counting_clean_cron" json:"counting_clean_cron"`
	IsCountingEnabled  *bool                  `bun:"is_counting_enabled,type:bool" json:"is_counting_enabled"`
	TC                 string                 `bun:"tc" json:"tc"`
}

type SettingsResponse struct {
	bun.BaseModel      `json:"-" bun:"table:settings"`
	CarParkID          int                    `bun:"carpark_id" json:"carpark_id"`
	CarParkName        map[string]interface{} `bun:"carpark_name,type:jsonb"  json:"carpark_name" swaggertype:"object"`
	AppLogo            string                 `bun:"app_logo,type:bytea" json:"app_logo"`
	DefaultLang        string                 `bun:"default_lang" binding:"required" json:"default_lang" `
	TimeOutScreenKisok *int                   `bun:"timeout_screenKiosk" binding:"required" json:"timeout_screenKiosk"`
}

type ResponseData struct {
	General     GeneralInfo `json:"general"`
	DefaultLang string      `json:"default_lang"`
	Cron        CronInfo    `json:"cron"`
	Kiosk       KioskInfo   `json:"Kiosk"`
}

type GeneralInfo struct {
	CarParkID    int    `json:"carpark_id"`
	CarParkName  string `json:"carpark_name"`
	PkaImageSize string `json:"pka_image_size"`
}

type CronInfo struct {
	FycCleanCron      int  `json:"fyc_clean_cron"`
	CountingCleanCron int  `json:"counting_clean_cron"`
	IsFycEnabled      bool `json:"is_fyc_enabled"`
	IsCountingEnabled bool `json:"is_counting_enabled"`
}

type KioskInfo struct {
	TimeOutScreenKiosk int    `json:"timeout_screenKiosk"`
	AppLogo            string `json:"app_logo"`
	TC                 string `json:"tc"`
}

func SettingsExists(ctx context.Context, carpark_id int) (bool, error) {
	var count int
	count, err := Db_GlobalVar.NewSelect().
		Model((*Settings)(nil)).
		Where("carpark_id = ?", carpark_id).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateSettings(ctx context.Context, settings *Settings) error {
	settings.DefaultLang = strings.ToLower(settings.DefaultLang)
	settings.IsCountingEnabled = true
	settings.IsFycEnabled = true

	_, err := Db_GlobalVar.NewInsert().Model(settings).Exec(ctx)
	return err
}

// GetSettings fetches a settings entry by CarParkID
func GetSettings(ctx context.Context, carParkID int) (*Settings, error) {
	var settings Settings
	err := Db_GlobalVar.NewSelect().Model(&settings).Where("carpark_id = ?", carParkID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func GetAllSettings(ctx context.Context) (*Settings, error) {
	var settings Settings
	settings.DefaultLang = strings.ToLower(settings.DefaultLang)

	err := Db_GlobalVar.NewSelect().Model(&settings).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func GetAllSettingsThirdParty(ctx context.Context) (*SettingsResponse, error) {
	var settingsTP SettingsResponse
	err := Db_GlobalVar.NewSelect().Model(&settingsTP).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &settingsTP, nil
}

// UpdateSettings updates a settings entry by CarParkID
func UpdateSettings(ctx context.Context, settings SettingsNoBind, cp_id int) error {
	log.Debug().Int("Settings", cp_id)
	settings.DefaultLang = strings.ToLower(settings.DefaultLang)

	res, err := Db_GlobalVar.NewUpdate().
		Model(&settings).
		Where("carpark_id = ?", cp_id).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("no rows updated")
	}
	return nil
}
