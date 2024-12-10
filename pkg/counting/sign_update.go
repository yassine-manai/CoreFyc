package counting

import (
	"context"
	"fyc/pkg/db"

	"github.com/rs/zerolog/log"
)

func Increase_Zone_Capacity(ctx context.Context, CurrZone int, lpn string) int {

	log.Info().Msg("Sign Increase Operation ----------- ")

	zone, err := db.GetZoneByID(ctx, CurrZone)
	if err != nil {
		log.Error().Str("Error: ", err.Error()).Int("Zone ID", CurrZone).Msg("Error retrieving zone")
	}

	if *zone.FreeCapacity < *zone.MaxCapacity {
		rows, _ := db.UpdateZoneCapacity(ctx, CurrZone, "inc")
		if rows == 0 {
			log.Error().Str("Licence Plate", lpn).Int("ROWS AFF", int(rows)).Msg("No rows affected")
		}
	}

	zoneData, err := db.GetZoneByID(ctx, CurrZone)
	if err != nil {
		log.Error().Str("Error: ", err.Error()).Int("Zone ID", CurrZone).Msg("Error retrieving zone")
	}

	return *zoneData.FreeCapacity

}

func Decrease_Zone_Capacity(ctx context.Context, CurrZone int, lpn string) int {

	log.Info().Msg("Sign Decrease Operation ----------- ")

	zoneData, err := db.GetZoneByID(ctx, CurrZone)
	if err != nil {
		log.Error().Str("Error: ", err.Error()).Int("Zone ID", CurrZone).Msg("Error retrieving zone")
	}

	if *zoneData.FreeCapacity > 0 && *zoneData.FreeCapacity <= *zoneData.MaxCapacity {
		rows, _ := db.UpdateZoneCapacity(ctx, CurrZone, "dec")
		if rows == 0 {
			log.Error().Str("Licence Plate", lpn).Int("ROWS AFF", int(rows)).Msg("No rows affected")
		}
	} else {
		log.Warn().Str("Licence Plate", lpn).Int("Free Capacity", *zoneData.FreeCapacity).Msg("No free capacity in current zone")
	}

	zoneData, err = db.GetZoneByID(ctx, CurrZone)
	if err != nil {
		log.Error().Str("Error: ", err.Error()).Int("Zone ID", CurrZone).Msg("Error retrieving zone")
	}

	return *zoneData.FreeCapacity
}
