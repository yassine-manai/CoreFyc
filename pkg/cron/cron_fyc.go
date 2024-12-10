package cron

import (
	"context"
	"database/sql"
	"fmt"
	"fyc/pkg/db"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

func CronFyc() {
	ctx := context.Background()
	defaultCronHour := 00
	cronHour := defaultCronHour
	CronEnabled := false

	TimeSettings, err := db.GetAllSettings(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			// No settings found, use default value
			//log.Warn().Msg("No settings found in the database -- using default time %v:00", )
			cronHour = defaultCronHour

			log.Warn().Msgf("No settings found in the database -- CRON FYC DISABLED -- Error %s", err)

		} else {
			log.Warn().Msgf("Error retrieving Settings from the database -- Error %s", err)
			cronHour = defaultCronHour
		}
	} else {
		cronHour = TimeSettings.FycCleanCron
		CronEnabled = TimeSettings.IsFycEnabled
	}

	if CronEnabled {
		cronExpression := fmt.Sprintf("0 %d * * *", cronHour)
		c := cron.New()

		_, err = c.AddFunc(cronExpression, CronJobFyc)
		if err != nil {
			log.Err(err).Msg("Error adding cron fyc job")
			return
		}

		c.Start()
		log.Info().Msgf("------------------------------ # Cron Job FYC Scheduled to Run Daily at %v:00 # ------------------------------", cronHour)

		select {}
	} else {
		log.Warn().Msg("------------------------------ # CRON FYC DISABLED # ------------------------------ ")
	}

}

func CronJobFyc() {
	ctx := context.Background()

	log.Debug().Msg("------------------------------ # Cron Job STARTED # ------------------------------ ")
	if err := db.Db_GlobalVar.ResetModel(ctx, &db.PresentCar{}); err != nil {
		log.Error().Str("Error", err.Error()).Msg("Failed to reset PresentCar table")
		return
	}

	log.Info().Msg("Cron FYC Successfully worked")
	log.Debug().Msg("------------------------------ # Cron FYC Job FINISHED # ------------------------------ ")
}
