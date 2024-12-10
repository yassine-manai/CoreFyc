package db

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"fyc/functions"
)

type PresentCar struct {
	bun.BaseModel   `json:"-" bun:"table:presentcar"`
	ID              *int                   `bun:"id,pk,autoincrement" json:"id"`
	TransactionDate string                 `bun:"transaction_date,type:timestamp" json:"transaction_date" binding:"required"`
	CameraID        int                    `bun:"camera_id,pk" json:"camera_id" binding:"required"`
	LPN             string                 `bun:"lpn,pk" json:"lpn" binding:"required"`
	CurrZoneID      *int                   `bun:"current_zone_id" json:"current_zone_id" binding:"required"`
	LastZoneID      *int                   `bun:"last_zone_id" json:"last_zone_id" binding:"required"`
	Direction       string                 `bun:"direction" json:"direction" binding:"required"`
	Confidence      *int                   `bun:"confidence" json:"confidence" binding:"required"`
	CarDetailsID    *int                   `bun:"car_details_id" json:"car_details_id" binding:"required"`
	Extra           map[string]interface{} `bun:"extra,type:jsonb" json:"extra" swaggertype:"object"`
}

type ResponsePC struct {
	bun.BaseModel   `json:"-" bun:"table:presentcar"`
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
func GetAllPresentExtra(ctx context.Context) ([]PresentCar, error) {
	var cars []PresentCar
	err := Db_GlobalVar.NewSelect().Model(&cars).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all present cars with Extra: %w", err)
	}
	for i := range cars {
		cars[i].TransactionDate, _ = functions.ParseTimeData(cars[i].TransactionDate)
	}
	return cars, nil
}

// Get all present cars with  extra
func GetAllPresentCarsBiDateExtra(ctx context.Context, startDate, endDate string) ([]PresentCar, error) {
	var Pcars []PresentCar
	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("transaction_date::date BETWEEN ? AND ?", startDate, endDate).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars between %s and %s: %w", startDate, endDate, err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}
	return Pcars, nil
}

func GetAllPCLpnBiDateExtra(ctx context.Context, startDate, endDate, lpn string) ([]PresentCar, error) {
	var Pcars []PresentCar
	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("lpn LIKE ?", "%"+lpn+"%").
		Where("transaction_date::date BETWEEN ? AND ?", startDate, endDate).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars between %s and %s: %w", startDate, endDate, err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}
	return Pcars, nil
}

func GetAllLPNZONE(ctx context.Context, startDate, endDate, lpn string, zone int) ([]PresentCar, error) {
	var Pcars []PresentCar
	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("lpn LIKE ? AND current_zone_id = ?", "%"+lpn+"%", zone).
		Where("transaction_date::date BETWEEN ? AND ?", startDate, endDate).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars between %s and %s: %w", startDate, endDate, err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}
	return Pcars, nil
}

func GetZLSE(ctx context.Context, startDate, endDate string, zone int) ([]PresentCar, error) {
	var Pcars []PresentCar
	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("current_zone_id = ?", zone).
		Where("transaction_date::date BETWEEN ? AND ?", startDate, endDate).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars between %s and %s: %w", startDate, endDate, err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}
	return Pcars, nil
}

func GetAllPresent(ctx context.Context, startDate, endDate string, zone int) ([]PresentCar, error) {
	var Pcars []PresentCar
	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("current_zone_id != ?", zone).
		Where("transaction_date::date BETWEEN ? AND ?", startDate, endDate).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars between %s and %s: %w", startDate, endDate, err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}
	return Pcars, nil
}

// Get all present cars with no extra
func GetAllPresentCarsBiDateNoExtra(ctx context.Context, startDate, endDate string) ([]PresentCar, error) {
	var Pcars []PresentCar

	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("transaction_date::date BETWEEN ? AND ?", startDate, endDate).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars between %s and %s: %w", startDate, endDate, err)
	}

	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}
	return Pcars, nil
}

// Get all present cars with extra
func GetAllPresentCarsCurrTimeExtra(ctx context.Context) ([]PresentCar, error) {
	var Pcars []PresentCar

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
func GetAllPresentCarsCurrTimeNoExtra(ctx context.Context) ([]ResponsePC, error) {
	var Pcars []ResponsePC

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

func GetAllPresentCarsBiDateID(ctx context.Context, startDate, endDate, ID string) ([]PresentCar, error) {
	var Pcars []PresentCar
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

func GetAllPresentCarsBiDateExtraID(ctx context.Context, startDate, endDate, ID string) ([]ResponsePC, error) {
	var Pcars []ResponsePC
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

func GetPresentCarsID(ctx context.Context, ID string) (*PresentCar, error) {
	var Pcars PresentCar
	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("id = ?", ID).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars with id %s: %w", ID, err)
	}

	Pcars.TransactionDate, _ = functions.ParseTimeData(Pcars.TransactionDate)

	return &Pcars, nil
}

func GetCarsLPN(ctx context.Context, lpn string) (*PresentCar, error) {
	var Pcars PresentCar
	err := Db_GlobalVar.NewSelect().
		Model(&Pcars).
		Where("lpn = ?", lpn).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars with lpn %s: %w", lpn, err)
	}

	Pcars.TransactionDate, _ = functions.ParseTimeData(Pcars.TransactionDate)

	return &Pcars, nil
}

// -------------------------------------------- # DONT CHANGE THIS # --------------------------------------------
// Get all present cars
func GetAllPresentCars(ctx context.Context) ([]ResponsePC, error) {
	var Pcars []ResponsePC
	err := Db_GlobalVar.NewSelect().Model(&Pcars).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all present cars: %w", err)
	}
	for i := range Pcars {
		Pcars[i].TransactionDate, _ = functions.ParseTimeData(Pcars[i].TransactionDate)
	}

	return Pcars, nil
}

// Get present car by LPN
func GetPresentCarByLPN(ctx context.Context, lpn string) (*PresentCar, error) {
	var car PresentCar
	err := Db_GlobalVar.NewSelect().Model(&car).Where("lpn = ?", lpn).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting present car by LPN: %w", err)
	}
	car.TransactionDate, _ = functions.ParseTimeData(car.TransactionDate)
	return &car, nil
}

func GetPresentCarByLPNs(ctx context.Context, lpn string) ([]PresentCar, error) {
	var cars []PresentCar
	err := Db_GlobalVar.NewSelect().Model(&cars).Where("lpn = ?", lpn).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting present cars by LPN: %w", err)
	}

	for i := range cars {
		cars[i].TransactionDate, _ = functions.ParseTimeData(cars[i].TransactionDate)
	}
	return cars, nil
}

func GetPresentCarByLPNFuzzy(ctx context.Context, lpn string) ([]PresentCar, error) {
	var cars []PresentCar

	err := Db_GlobalVar.NewSelect().
		Model(&cars).
		Where("lpn LIKE ?", lpn).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting present cars by LPN: %w", err)
	}

	// Parse TransactionDate for each car in the result
	for i := range cars {
		cars[i].TransactionDate, _ = functions.ParseTimeData(cars[i].TransactionDate)
	}
	return cars, nil
}

func GetPresentFound(ctx context.Context, lpn string) (bool, error) {
	var car PresentCar
	err := Db_GlobalVar.NewSelect().Model(&car).Where("lpn = ?", lpn).Scan(ctx)
	if err != nil {
		return false, fmt.Errorf("error getting present car by LPN: %w", err)
	}
	return true, nil
}

// Create a new present car
func CreatePresentCar(ctx context.Context, car *PresentCar) error {
	_, err := Db_GlobalVar.NewInsert().Model(car).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating present car: %w", err)
	}
	log.Info().Msgf("New present car added with ID: %d", *car.ID)

	return nil
}

// Update a present car by ID and return rows affected
func UpdatePresentCar(ctx context.Context, id int, updates *PresentCar) (int64, error) {
	res, err := Db_GlobalVar.NewUpdate().Model(updates).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected() // Get the number of rows affected
	log.Debug().Msgf("Updated present car with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}

// update by LPN
func UpdatePresentCarByLpn(ctx context.Context, lpn string, updates *PresentCar) (int64, error) {
	log.Debug().Str("lpn", lpn).Msgf("Update Present Car by LPN")
	//log.Debug().Interface("DATA", updates).Send()

	res, err := Db_GlobalVar.NewUpdate().Model(updates).Where("lpn = ?", lpn).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated present car with LPN: %s, rows affected: %d", lpn, rowsAffected)

	return rowsAffected, nil
}

// Delete a present car by ID and return rows affected
func DeletePresentCar(ctx context.Context, id int) (int64, error) {
	res, err := Db_GlobalVar.NewDelete().Model(&PresentCar{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted present car with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
