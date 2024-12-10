package cron

import (
	"context"
	"database/sql"
	"fmt"
	"fyc/pkg/db"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

func CronCounting() {
	ctx := context.Background()
	defaultCronHour := 00
	cronHour := defaultCronHour
	CronEnabled := false

	TimeSettings, err := db.GetAllSettings(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			cronHour = defaultCronHour

			log.Warn().Msgf("No settings found in the database -- CRON COUNTING DISABLED -- Error %s", err)

		} else {
			log.Warn().Msgf("Error retrieving Settings from the database -- Error %s", err)
			cronHour = defaultCronHour
		}
	} else {
		cronHour = TimeSettings.CountingCleanCron
		CronEnabled = TimeSettings.IsCountingEnabled
	}

	if CronEnabled {
		cronExpression := fmt.Sprintf("0 %d * * *", cronHour)
		c := cron.New()

		_, err = c.AddFunc(cronExpression, CronJobCounting)
		if err != nil {
			log.Err(err).Msg("Error adding cron counting job")
			return
		}

		c.Start()
		log.Info().Msgf("------------------------------ # Cron Job Counting Scheduled to Run Daily at %v:00 # ------------------------------", cronHour)

		select {}
	} else {
		log.Warn().Msg("------------------------------ # CRON COUNTING DISABLED # ------------------------------ ")
	}

}

func CronJobCounting() {
	//ctx := context.Background()

	log.Debug().Msg("------------------------------ # Cron Counting Job STARTED # ------------------------------ ")
	/* if err := db.Db_GlobalVar.ResetModel(ctx, &db.PresentCar{}); err != nil {
		log.Error().Str("Error", err.Error()).Msg("Failed to reset PresentCar table")
		return
	}

	log.Info().Msg("Cron Successfully worked") */
	log.Debug().Msg("------------------------------ # Cron Counting Job FINISHED # ------------------------------ ")
}
