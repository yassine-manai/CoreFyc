package counting

import (
	"context"
	"fmt"
	"fyc/pkg/db"
	"fyc/pkg/valkey"

	"github.com/rs/zerolog/log"
)

func Sign_Data_Values(zone_id int, meth string, places_free string) string {
	ctx := context.Background()
	var valk valkey.ValkeyStrct

	valk.Valkey_Connect()

	if meth == "inc" {
		log.Debug().Msg(" - - - SIGN Increase Value - - - ")

		sign, err := db.GetSignByZoneId(ctx, zone_id)
		if err != nil || sign == nil {
			log.Error().Err(err).Int("sign_id", zone_id).Msg("Error retrieving sign by ID or sign not found")
			return ""
		}

		log.Info().Str("sign IP", sign.SignIP).Msg("Successfully retrieved sign data")

		// Fetching Capacity in Zone
		capacity, err := db.GetZoneByID(ctx, zone_id)
		if err != nil {
			log.Error().Err(err).Int("zone_id", zone_id).Msg("Error fetching capacity in the zone")
			return ""
		}

		log.Info().Int("zone_id", zone_id).Int("capacity_count", *capacity.FreeCapacity).Msg("Fetched capacity in the zone")

		if *capacity.FreeCapacity < *capacity.MaxCapacity && *capacity.FreeCapacity > 0 {
			log.Info().Int("zone_id", zone_id).Int("FreeCapacity", *capacity.FreeCapacity).Int("MAXCapacity", *capacity.MaxCapacity).Msg("Sign data values are valid")
			var SignHost = fmt.Sprintf("%s:%d", sign.SignIP, sign.SignPort)

			//valk.Valkey_Incr_Data(ctx, SignHost)
			valk.Valkey_Setter_Data(ctx, SignHost, places_free)
			valk.PublishMessage(ctx, SignHost)
		}

		return sign.SignIP
	}

	log.Debug().Msg(" - - - SIGN Decrease Value - - - ")
	sign, err := db.GetSignByZoneId(ctx, zone_id)
	if err != nil || sign == nil {
		log.Error().Err(err).Int("sign_id", zone_id).Msg("Error retrieving sign by ID or sign not found")
		return ""
	}

	log.Info().Str("sign IP", sign.SignIP).Msg("Successfully retrieved sign data")

	// Fetching Capacity in Zone
	capacity, err := db.GetZoneByID(ctx, zone_id)
	if err != nil {
		log.Error().Err(err).Int("zone_id", zone_id).Msg("Error fetching capacity in the zone")
		return ""
	}

	log.Info().Int("zone_id", zone_id).Int("capacity_count", *capacity.FreeCapacity).Msg("Fetched capacity in the zone")

	if *capacity.FreeCapacity < *capacity.MaxCapacity && *capacity.FreeCapacity > 0 {
		var SignHost = fmt.Sprintf("%s:%d", sign.SignIP, sign.SignPort)
		log.Info().Int("zone_id", zone_id).Str("Sign HOST", SignHost).Int("places_count", *capacity.FreeCapacity).Msg("Sign data values are valid")

		//valk.Valkey_Decr_Data(ctx, SignHost)
		valk.Valkey_Setter_Data(ctx, SignHost, places_free)
		valk.PublishMessage(ctx, SignHost)
	}

	return sign.SignIP
}
