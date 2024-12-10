package db

import (
	"context"
	"fmt"
	"fyc/functions"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type PresentCarHistory struct {
	bun.BaseModel   `json:"-" bun:"table:present_car_history"`
	ID              *int                   `bun:"id,autoincrement" json:"id"`
	TransactionDate string                 `bun:"transaction_date,type:timestamp" json:"transaction_date" binding:"required"`
	CameraID        int                    `bun:"camera_id" json:"camera_id" binding:"required"`
	LPN             string                 `bun:"lpn" json:"lpn" binding:"required"`
	CurrZoneID      *int                   `bun:"current_zone_id" json:"current_zone_id" binding:"required"`
	LastZoneID      *int                   `bun:"last_zone_id" json:"last_zone_id" binding:"required"`
	Direction       string                 `bun:"direction" json:"direction" binding:"required"`
	Confidence      *int                   `bun:"confidence" json:"confidence" binding:"required"`
	CarDetailsID    *int                   `bun:"car_details_id" json:"car_details_id" binding:"required"`
	Extra           map[string]interface{} `bun:"extra,type:jsonb" json:"extra" swaggertype:"object"`
}

type ResponsePCH struct {
	bun.BaseModel   `json:"-" bun:"table:present_car_history"`
	ID              *int   `bun:"id" json:"id"`
	CarDetailsID    *int   `bun:"car_details_id" json:"car_details_id"`
	CameraID        int    `bun:"camera_id" json:"camera_id"`
	Confidence      *int   `bun:"confidence" json:"confidence"`
	CurrZoneID      *int   `bun:"current_zone_id" json:"current_zone_id"`
	LastZoneID      *int   `bun:"last_zone_id" json:"last_zone_id"`
	Direction       string `bun:"direction" json:"direction"`
	LPN             string `bun:"lpn" json:"lpn"`
	TransactionDate string `bun:"transaction_date" json:"transaction_date"`
}

// Get all present cars
func GetAllPresentHistoryExtra(ctx context.Context) ([]PresentCarHistory, error) {
	var cars []PresentCarHistory
	err := Db_GlobalVar.NewSelect().Model(&cars).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all present cars with Extra: %w", err)
	}
	for i := range cars {
		cars[i].TransactionDate, _ = functions.ParseTimeData(cars[i].TransactionDate)
	}
	return cars, nil
}

// Get all present cars
func GetAllPresentCarsHistory(ctx context.Context) ([]ResponsePCH, error) {
	var Pcars []ResponsePCH
	err := Db_GlobalVar.NewSelect().Model(&Pcars).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all present cars in history: %w", err)
	}
	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}
	return Pcars, nil
}

// Get present car by LPN
func GetPresentCarByLPNHistory(ctx context.Context, lpn string) (*PresentCarHistory, error) {
	var car PresentCarHistory
	err := Db_GlobalVar.NewSelect().Model(&car).Where("lpn = ?", lpn).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting present car in history  by LPN: %w", err)
	}
	car.TransactionDate, _ = functions.ParseTimeData(car.TransactionDate)
	return &car, nil
}

// Create a new present car
func CreatePresentCarHistory(ctx context.Context, car *PresentCarHistory) error {
	// Insert and get the auto-generated ID from the database
	_, err := Db_GlobalVar.NewInsert().Model(car).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding present car to History: %w", err)
	}
	log.Info().Msgf("New present car added to History with ID: %d", *car.ID)

	return nil
}

// Get all present cars with extra
func GetAllHistoryCarsCurrTimeExtra(ctx context.Context) ([]PresentCarHistory, error) {
	var Pcars []PresentCarHistory

	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("transaction_date::date = CURRENT_DATE").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all present cars: %w", err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}

	return Pcars, nil
}

// Get all present cars with no extra
func GetAllHistoryCarsCurrTimeNoExtra(ctx context.Context) ([]ResponsePCH, error) {
	var Pcars []ResponsePCH

	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("transaction_date::date = CURRENT_DATE").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all present cars: %w", err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}

	return Pcars, nil
}

func GetAllPresentHistoryCarsBiDateID(ctx context.Context, startDate, endDate, ID string) ([]PresentCarHistory, error) {
	var Pcars []PresentCarHistory
	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("transaction_date::date BETWEEN ? AND ?", startDate, endDate).
		Where("id = ?", ID).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars between %s and %s with spec_id %s: %w", startDate, endDate, ID, err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}
	return Pcars, nil
}

func GetAllHistoryresentCarsBiDateExtraID(ctx context.Context, startDate, endDate, ID string) ([]ResponsePCH, error) {
	var Pcars []ResponsePCH
	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("transaction_date::date BETWEEN ? AND ?", startDate, endDate).
		Where("id = ?", ID).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars between %s and %s with spec_id %s: %w", startDate, endDate, ID, err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}
	return Pcars, nil
}

// Update a present car by ID and return rows affected
func UpdatePresentCarHistory(ctx context.Context, id int, updates *PresentCarHistory) (int64, error) {
	res, err := Db_GlobalVar.NewUpdate().Model(updates).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating present car in history: %w", err)
	}

	rowsAffected, _ := res.RowsAffected() // Get the number of rows affected
	log.Debug().Msgf("Updated present car in history with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}

// update by LPN
func UpdatePresentCarByLpnHistory(ctx context.Context, lpn string, updates *PresentCarHistory) (int64, error) {
	log.Debug().Str("lpn", lpn).Msgf("Update Present Car by LPN:%v", updates)
	res, err := Db_GlobalVar.NewUpdate().Model(updates).Where("lpn = ?", lpn).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating present car to history: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated present car in history with LPN: %s, rows affected: %d", lpn, rowsAffected)

	return rowsAffected, nil
}

// Delete a present car by ID and return rows affected
func DeletePresentCarHistory(ctx context.Context, id int) (int64, error) {
	res, err := Db_GlobalVar.NewDelete().Model(&PresentCarHistory{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted present car with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
