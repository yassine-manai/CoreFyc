package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"fyc/functions"
)

type Zone struct {
	bun.BaseModel `json:"-" bun:"table:zone"`
	ID            int                    `bun:"id,pk,autoincrement" json:"-"`
	ZoneID        int                    `bun:"zone_id" json:"zone_id" binding:"required"`
	Name          map[string]interface{} `bun:"name,type:jsonb" json:"name" binding:"required" swaggertype:"object"`
	MaxCapacity   *int                   `bun:"max_capacity" json:"max_capacity" binding:"required"`
	FreeCapacity  *int                   `bun:"free_capacity" json:"free_capacity" binding:"required"`
	LastUpdated   string                 `bun:"last_update,type:timestamp" json:"-"`
	IsEnabled     bool                   `bun:"is_enabled,type:bool" json:"-"`
	IsDeleted     bool                   `bun:"is_deleted,type:bool" json:"-"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" swaggertype:"object"`
}

type ZoneNoBind struct {
	bun.BaseModel `json:"-" bun:"table:zone"`
	ID            int                    `bun:"id,pk,autoincrement" json:"-"`
	ZoneID        int                    `bun:"zone_id,pk" json:"-" `
	Name          map[string]interface{} `bun:"name,type:jsonb" json:"name" swaggertype:"object"`
	MaxCapacity   *int                   `bun:"max_capacity" json:"max_capacity" `
	FreeCapacity  *int                   `bun:"free_capacity" json:"free_capacity" `
	LastUpdated   string                 `bun:"last_update,type:timestamp" json:"-"`
	IsEnabled     *bool                  `bun:"is_enabled,type:bool" json:"is_enabled"`
	IsDeleted     *bool                  `bun:"is_deleted,type:bool" json:"is_deleted"`
	Extra         map[string]interface{} `bun:"extra" json:"extra" swaggertype:"object"`
}

type ZoneName struct {
	Ar string `json:"ar" `
	En string `json:"en" `
}

type ResponseZone struct {
	bun.BaseModel `json:"-" bun:"table:zone"`
	ID            int                    `bun:"id" json:"-"`
	ZoneID        *int                   `bun:"zone_id" json:"zone_id"`
	Name          ZoneName               `bun:"name" json:"name"`
	MaxCapacity   *int                   `bun:"max_capacity" json:"max_capacity"`
	FreeCapacity  *int                   `bun:"free_capacity" json:"free_capacity"`
	LastUpdated   string                 `bun:"last_update" json:"last_update"`
	IsEnabled     bool                   `bun:"is_enabled" json:"is_enabled"`
	IsDeleted     bool                   `bun:"is_deleted" json:"-"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" swaggertype:"object"`
}

type ResponseZoneExtra struct {
	bun.BaseModel `json:"-" bun:"table:zone"`
	ID            int      `bun:"id" json:"-"`
	ZoneID        *int     `bun:"zone_id" json:"zone_id"`
	Name          ZoneName `bun:"name" json:"name"`
	MaxCapacity   *int     `bun:"max_capacity" json:"max_capacity"`
	FreeCapacity  *int     `bun:"free_capacity" json:"free_capacity"`
	LastUpdated   string   `bun:"last_update" json:"last_update"`
	IsEnabled     bool     `bun:"is_enabled" json:"is_enabled"`
	IsDeleted     bool     `bun:"is_deleted" json:"-"`
}

func GetZoneData(ctx context.Context) ([]Zone, error) {
	var zone []Zone
	err := Db_GlobalVar.NewSelect().Model(&zone).Scan(ctx)
	if err != nil {

		return nil, fmt.Errorf("error getting all Zones with Extra: %w", err)
	}
	return zone, nil
}

func GetZones(ctx context.Context) ([]ResponseZone, error) {
	var zone []ResponseZone
	err := Db_GlobalVar.NewSelect().Model(&zone).Scan(ctx)
	if err != nil {

		return nil, fmt.Errorf("error getting all Zones with Extra: %w", err)
	}
	return zone, nil
}

// Get all Zones with extra data
func GetAllZoneExtra(ctx context.Context) ([]ResponseZone, error) {
	var zone []ResponseZone
	err := Db_GlobalVar.NewSelect().
		Model(&zone).
		ExcludeColumn().
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all Zones with Extra: %w", err)
	}

	for i := range zone {
		zone[i].LastUpdated, _ = functions.ParseTimeData(zone[i].LastUpdated)
	}

	return zone, nil
}

func GetAllZoneNoExtra(ctx context.Context) ([]ResponseZoneExtra, error) {
	var zone []ResponseZoneExtra
	err := Db_GlobalVar.NewSelect().
		Model(&zone).
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all Zones without Extra: %w", err)
	}

	for i := range zone {
		zone[i].LastUpdated, _ = functions.ParseTimeData(zone[i].LastUpdated)
	}

	return zone, nil
}

// Get all zone
func GetAllZone(ctx context.Context) ([]ResponseZone, error) {
	var EZ []ResponseZone
	err := Db_GlobalVar.NewSelect().
		Model(&EZ).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all zones : %w", err)
	}
	for i := range EZ {
		EZ[i].LastUpdated, _ = functions.ParseTimeData(EZ[i].LastUpdated)
	}
	return EZ, nil
}

// Get zone by id
func GetZoneByID(ctx context.Context, id int) (*Zone, error) {
	log.Debug().Int("ID ZONE", id)
	var zone Zone

	err := Db_GlobalVar.NewSelect().
		Model(&zone).
		Where("is_deleted = ?", false).
		Where("zone_id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting zone by id : %w", err)
	}

	zone.LastUpdated, _ = functions.ParseTimeData(zone.LastUpdated)

	return &zone, nil
}

func GetPresentZoneByID(ctx context.Context, id int) (*Zone, error) {
	log.Debug().Int("ID ZONE", id)
	var zone Zone

	err := Db_GlobalVar.NewSelect().
		Model(&zone).
		Where("is_deleted = ?", false).
		Where("zone_id != ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting zone by id : %w", err)
	}

	zone.LastUpdated, _ = functions.ParseTimeData(zone.LastUpdated)

	return &zone, nil
}
func GetZoneByIDExport(ctx context.Context, id int) (*ResponseZone, error) {
	log.Debug().Int("ID ZONE", id)
	var zone ResponseZone

	err := Db_GlobalVar.NewSelect().
		Model(&zone).
		Where("zone_id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting zone by id : %w", err)
	}

	zone.LastUpdated, _ = functions.ParseTimeData(zone.LastUpdated)

	return &zone, nil
}

func GetZoneByStatus(ctx context.Context, status string) ([]ResponseZone, error) {
	var zones []ResponseZone

	// Start building the query
	query := Db_GlobalVar.NewSelect().
		Model(&zones).
		Where("is_deleted = ?", false)

	switch status {
	case "enabled":
		query.Where("is_enabled = ?", true)
	case "disabled":
		query.Where("is_enabled = ?", false)
	case "":
		// No status filter applied
	default:
		return nil, fmt.Errorf("invalid status: %s, expected 'enabled', 'disabled', or empty", status)
	}

	// Execute the query
	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("error getting Zone -- status: %s, err: %w", status, err)
	}

	// Parse `LastUpdated` for each zone
	for i := range zones {
		zones[i].LastUpdated, _ = functions.ParseTimeData(zones[i].LastUpdated)
	}

	return zones, nil
}

func GetZoneByIDNoExtra(ctx context.Context, id int) (*ResponseZone, error) {
	var zone ResponseZone

	err := Db_GlobalVar.NewSelect().
		Model(&zone).
		Where("is_deleted = ?", false).
		Where("zone_id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting zone by id : %w", err)
	}

	zone.LastUpdated, _ = functions.ParseTimeData(zone.LastUpdated)

	return &zone, nil
}

// create a new zone
func CreateZone(ctx context.Context, zone *Zone) error {
	log.Debug().Str("Zone Creation Time:", functions.GetFormatedLocalTime()).Int("ZoneID", zone.ZoneID).Msg("Creating new zone")
	zone.LastUpdated = functions.GetFormatedLocalTime()
	zone.IsDeleted = false
	zone.IsEnabled = true

	// Initialize a new map to hold normalized names
	normalizedNames := make(map[string]interface{}, len(zone.Name))
	for lang, name := range zone.Name {
		if str, ok := name.(string); ok && str != "" {
			// Use the original name while lowering the key
			normalizedNames[strings.ToLower(lang)] = str
		}
	}
	zone.Name = normalizedNames

	_, err := Db_GlobalVar.NewInsert().Model(zone).Returning("zone_id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating zone: %w", err)
	}

	log.Debug().Msgf("New zone added with ID: %d", zone.ZoneID)
	LoadzoneList()
	LoadAllZonelist()

	return nil
}

// Update a zone by ID
func UpdateZone(ctx context.Context, zone_id int, updates ZoneNoBind) (int64, error) {
	// Log update time
	log.Debug().Str("Updated at", functions.GetFormatedLocalTime()).Int("Zone ID", zone_id).Msg("Starting zone update")
	updates.LastUpdated = functions.GetFormatedLocalTime()

	// Initialize a new map to hold normalized names
	normalizedNames := make(map[string]interface{}, len(updates.Name))
	for lang, name := range updates.Name {
		if str, ok := name.(string); ok && str != "" {
			// Use the original name while lowering the key
			normalizedNames[strings.ToLower(lang)] = str
		}
	}
	updates.Name = normalizedNames

	// log.Debug().Int("ZoneId", zone_id).Msgf("Data --- %v", updates)

	res, err := Db_GlobalVar.NewUpdate().
		Model(&updates).
		Where("zone_id = ?", zone_id).
		//Set("last_update = ?", GetFormatedLocalTime()). this is wrong
		OmitZero().
		Exec(ctx)

	log.Debug().Msgf(" Query: %s ", res)

	if err != nil {
		log.Error().Err(err).Msgf("Error updating zone with id %d", zone_id)
		return 0, fmt.Errorf("error updating zone with id %d: %w", zone_id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving rows affected")
		return 0, fmt.Errorf("error retrieving rows affected for zone with id %d: %w", zone_id, err)
	}
	// if updates.IsDeleted != nil {
	// 	log.Debug().Msgf("IsDeleted : %v", updates.IsDeleted)
	// 	Db_GlobalVar.NewUpdate().
	// 		Model(&Zone{}).
	// 		Set("is_deleted = ?", updates.IsDeleted).
	// 		// Where("is_deleted = ?", false).
	// 		Where("zone_id = ?", zone_id).
	// 		Exec(ctx)
	// }

	// if updates.IsEnabled != nil {
	// 	log.Debug().Msgf("is_enabled : %v", updates.IsEnabled)
	// 	Db_GlobalVar.NewUpdate().
	// 		Model(&Zone{}).
	// 		Set("is_enabled = ?", updates.IsEnabled).
	// 		// Where("is_deleted = ?", false).
	// 		Where("zone_id = ?", zone_id).
	// 		Exec(ctx)
	// }

	log.Debug().Msgf("Successfully updated zone with ID: %d, rows affected: %d", zone_id, rowsAffected)
	LoadzoneList()
	LoadAllZonelist()
	return rowsAffected, nil
}

// Delete a zone by ID
func DeleteZone(ctx context.Context, zone_id int) (int64, error) {
	log.Debug().Str("Deleted at:", functions.GetFormatedLocalTime()).Int("Zone ID:", zone_id)

	var zone *Zone
	res, err := Db_GlobalVar.NewUpdate().
		Model(zone).
		Where("zone_id = ?", zone_id).
		//Where("is_enabled = ?", true).
		Set("is_deleted = ?", true).
		Set("is_enabled = ?", false).
		Set("last_update = ?", functions.GetFormatedLocalTime()).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting Zone with id %d: %w", zone_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted zone with ID: %d, rows affected: %d", zone_id, rowsAffected)
	LoadAllZonelist()
	LoadzoneList()
	return rowsAffected, nil
}

func DeleteZoneError(ctx context.Context, zone_id int) error {
	var zone *Zone
	_, err := Db_GlobalVar.NewDelete().
		Model(zone).
		Where("zone_id = ?", zone_id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deleting Zone with id %d: %w", zone_id, err)
	}

	LoadAllZonelist()
	LoadzoneList()
	return nil
}

// Update a zone capacity by ID
func UpdateZoneCapacity(ctx context.Context, zoneID int, meth string) (int64, error) {
	var updateZone ZoneNoBind
	var adjustment int

	switch meth {
	case "dec":
		adjustment = -1
	case "inc":
		adjustment = 1
	default:
		log.Error().Str("method", meth).Msg("Invalid method for updating zone capacity")
		return 0, fmt.Errorf("invalid method: %s, use 'dec' or 'inc'", meth)
	}

	// Log the operation
	log.Debug().
		Str("Updated at", functions.GetFormatedLocalTime()).
		Int("Zone ID", zoneID).
		Str("Operation", meth).
		Msg("Updating zone capacity")

	// Perform the update
	res, err := Db_GlobalVar.NewUpdate().
		Model(&updateZone).
		Where("zone_id = ?", zoneID).
		Set("free_capacity = free_capacity + ?", adjustment).
		OmitZero().
		Exec(ctx)

	if err != nil {
		log.Error().Err(err).Msgf("Error updating zone with id %d", zoneID)
		return 0, fmt.Errorf("error updating zone with id %d: %w", zoneID, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving rows affected")
		return 0, fmt.Errorf("error retrieving rows affected for zone with id %d: %w", zoneID, err)
	}

	log.Debug().
		Msgf("Successfully updated zone with ID: %d, rows affected: %d", zoneID, rowsAffected)

	LoadzoneList()
	LoadAllZonelist()

	return rowsAffected, nil
}
