package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ImageZone struct {
	bun.BaseModel `json:"-" bun:"table:zone_images"`
	ID            int                    `bun:"id,autoincrement,pk" json:"id"`
	ZoneID        *int                   `bun:"zone_id,pk" json:"zone_id" binding:"required"`
	Language      string                 `bun:"language,pk" json:"language" binding:"required"`
	ImageSm       string                 `bun:"image_s,type:bytea" json:"image_s" binding:"required"`
	ImageLg       string                 `bun:"image_l,type:bytea" json:"image_l" binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" swaggertype:"object"`
}

type ResponseImageZone struct {
	bun.BaseModel `json:"-" bun:"table:zone_images"`
	ID            int    `bun:"id" json:"id"`
	ZoneID        *int   `bun:"zone_id" json:"zone_id"`
	Language      string `bun:"language" json:"language"`
	ImageSm       string `bun:"image_s,type:bytea" json:"image_s"`
	ImageLg       string `bun:"image_l,type:bytea" json:"image_l"`
}

type ResponseImageLg struct {
	bun.BaseModel `json:"-" bun:"table:zone_images"`
	ID            int    `bun:"id" json:"id"`
	ZoneID        int    `bun:"zone_id" json:"zone_id"`
	Language      string `bun:"language" json:"language"`
	ImageLg       string `bun:"image_l,type:bytea" json:"image_l"`
}

type ResponseImageSm struct {
	bun.BaseModel `json:"-" bun:"table:zone_images"`
	ID            int    `bun:"id" json:"id"`
	ZoneID        int    `bun:"zone_id" json:"zone_id"`
	Language      string `bun:"language" json:"language"`
	ImageSm       string `bun:"image_s,type:bytea" json:"image_s"`
}

type ImageZoneNoBind struct {
	bun.BaseModel `json:"-" bun:"table:zone_images"`
	ID            int                    `bun:"id,autoincrement,pk" json:"id"`
	ZoneID        int                    `bun:"zone_id,pk" json:"zone_id" `
	Language      string                 `bun:"language,pk" json:"language"`
	ImageSm       string                 `bun:"image_s,type:bytea" json:"image_s"`
	ImageLg       string                 `bun:"image_l,type:bytea" json:"image_l" binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" swaggertype:"object"`
}

// Get all Zones with extra data
func GetAllZoneImageExtra(ctx context.Context) ([]ImageZone, error) {
	var zoneImage []ImageZone
	err := Db_GlobalVar.NewSelect().Model(&zoneImage).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Image Zones with Extra Data: %w", err)
	}
	return zoneImage, nil
}

// Get all zone
func GetAllZoneImage(ctx context.Context) ([]ResponseImageZone, error) {
	var Rzi []ResponseImageZone
	err := Db_GlobalVar.NewSelect().Model(&Rzi).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Zones Images : %w", err)
	}
	return Rzi, nil
}

// Get all zoneImage Small
func GetAllZoneImageSm(ctx context.Context) ([]ResponseImageSm, error) {
	var Zlg []ResponseImageSm

	err := Db_GlobalVar.NewSelect().Model(&Zlg).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Zones Images (small): %w", err)
	}

	return Zlg, nil
}

func GetAllZoneImageLg(ctx context.Context) ([]ResponseImageLg, error) {
	var Rzlg []ResponseImageLg

	err := Db_GlobalVar.NewSelect().Model(&Rzlg).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Zones Images (large): %w", err)
	}

	return Rzlg, nil
}

// Get zone by id
func GetZoneImageByID(ctx context.Context, id int) (*ImageZone, error) {
	zoneImg := new(ImageZone)
	err := Db_GlobalVar.NewSelect().Model(zoneImg).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Zone Image by id : %w", err)
	}
	return zoneImg, nil
}

func GetZoneImageByZONEIDLang(ctx context.Context, id int, language string) (*ImageZone, error) {
	var zoneImg ImageZone
	err := Db_GlobalVar.NewSelect().
		Model(&zoneImg).
		Where("zone_id = ? AND language = ?", id, language).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Zone Image by Language and id :%d, lang: %s err %w", id, language, err)
	}
	return &zoneImg, nil
}

func GetZoneImageByZONEIDLangs(ctx context.Context, id int, language string) ([]ImageZone, error) {
	var zoneImg []ImageZone
	err := Db_GlobalVar.NewSelect().
		Model(&zoneImg).
		Where("zone_id = ? AND language = ?", id, language).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Zone Image by Language and id :%d, lang: %s err %w", id, language, err)
	}
	return zoneImg, nil
}

func GetZoneImgByIDLang(ctx context.Context, id int, language string) (*ImageZone, error) {
	var zoneImg ImageZone
	err := Db_GlobalVar.NewSelect().
		Model(&zoneImg).
		Where("id = ?", id).
		Where("language = ?", language).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Zone Image by Language and id : %w", err)
	}
	return &zoneImg, nil
}

func GetAllImagesbyZoneID(ctx context.Context, zone_id int) ([]ImageZone, error) {
	var zoneImg []ImageZone
	err := Db_GlobalVar.NewSelect().
		Model(&zoneImg).
		Where("zone_id = ?", zone_id).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Zone Image by Language and id : %w", err)
	}
	return zoneImg, nil
}

// Gt zone by id
func GetZoneImageByZoneID(ctx context.Context, zone_id int) (*ImageZone, error) {
	var zoneImage ImageZone
	err := Db_GlobalVar.NewSelect().Model(&zoneImage).Where("zone_id = ?", zone_id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Zone Image by id : %w", err)
	}
	return &zoneImage, nil
}

// create a new zone
func CreateZoneImage(ctx context.Context, zoneImg *ImageZone) error {
	zoneImg.Language = strings.ToLower(zoneImg.Language)

	// Insert and get the auto-generated ID from the database
	_, err := Db_GlobalVar.NewInsert().Model(zoneImg).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating a Zone Image : %w", err)
	}
	log.Debug().Msgf("New zone image added with ID: %d", zoneImg.ID)

	return nil
}

// Update a zone img by ID
func UpdateZoneImage(ctx context.Context, zone_id int, updates *ImageZoneNoBind) (int64, error) {
	updates.Language = strings.ToLower(updates.Language)

	var rowsAffected int64
	log.Debug().Int("id", updates.ID).Int("Zoneid", updates.ZoneID).Str("lang", updates.Language)
	// res, err := Db_GlobalVar.NewUpdate().
	res, err := Db_GlobalVar.NewUpdate().
		Model(updates).
		//On("CONFLICT (id) DO UPDATE").
		//Set("id = EXCLUDED.id").
		Where("zone_id = ?", updates.ZoneID).
		Where("language = ?", updates.Language).
		//ExcludeColumn("id").
		OmitZero().
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating Zone Image with id %d: %w", zone_id, err)
	}

	rowsAffected, _ = res.RowsAffected()
	log.Debug().Msgf("Updated Zone Image with ID: %d, rows affected: %d ", zone_id, rowsAffected)
	if rowsAffected == 0 {
		log.Warn().Int("zone_id", zone_id).Str("lang", updates.Language).Msg("No row affected to update, insert it")
		r, err := Db_GlobalVar.NewInsert().
			Model(updates).
			// Where("zone_id = ?", zone_id).
			// ExcludeColumn("id").
			// OmitZero().
			Exec(ctx)
		if err != nil {
			return 0, fmt.Errorf("error updating Zone Image with id %d: %w", zone_id, err)
		}
		rowsAffected, _ = r.RowsAffected()
		log.Warn().Int("zone_id", zone_id).Int64("Affected", rowsAffected).Msg("ROW affected by insert it")
	}

	return rowsAffected, nil
}

// Delete a zone img by ID
func DeleteZoneImage(ctx context.Context, id int) (int64, error) {
	res, err := Db_GlobalVar.NewDelete().Model(&ImageZone{}).Where("zone_id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting Zone Image with id %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted Zone Image with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
