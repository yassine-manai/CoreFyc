package migrations

import (
	"fyc/pkg/db"

	"github.com/rs/zerolog/log"
)

/*
	 func Validation_DB() {

		err := ValidateTableSchema(db.UserAudit{})
		if err != nil {
			log.Err(err).Msg("Schema validation failed")
		} else {
			//log.Info().Msg("Schema validation passed!")
			log.Debug().Msg("VALIDATING SCHEMAS PASSED ")

		}


}
*/

func Validation_Shema() error {
	log.Debug().Msg(" ------------------------------- #  VALIDATING SCHEMAS STARTED # ------------------------------ \n")

	table := GetTableName(db.Camera{})
	GetStructColomns(db.Camera{})
	GetDBColomns(table)

	/* models := []interface{}{
		db.UserAudit{},
		db.User{},
		db.ApiKey{},

		db.Settings{},
		db.PresentCar{},
		db.PresentCarHistory{},
		db.ErrorMessage{},
		db.Zone{},
		db.ImageZone{},

		db.Camera{},
		db.CarDetail{},
		db.Sign{},
	}

	for md, model := range models {
		err := ValidateTableSchema(model)
		if err != nil {
			log.Err(err).Int("Model ", md).Msg("Schema validation failed")
		} else {
			log.Info().Int("Model ", md).Msg("Schema validation passed!")

		}
	} */
	log.Debug().Msg(" -------------------------------- # VALIDATING SCHEMAS ENDED # ------------------------------ \n")

	return nil
}
