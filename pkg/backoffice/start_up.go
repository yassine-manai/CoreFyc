package backoffice

import (
	"context"
	"fyc/config"
	"fyc/pkg/db"

	"github.com/rs/zerolog/log"
)

func StartUpData() {
	// Add Admin USER to database
	addDefaultAdminUser()

	// Add Default Settings Data to database
	addDefaultSettingsData()

	log.Debug().Msg("-------------------------------- # LOADING DATA LIST # ------------------------------")
	db.LoadzoneList()
	log.Debug().Msgf("Zonelist Length - - - - %v", len(db.Zonelist))

	db.LoadAllZonelist()
	log.Debug().Msgf("Zonelist ++  Length * * * *  %v", len(db.AllZonelist))

	db.LoadCameralist()
	log.Debug().Msgf("CameraList Length %v", len(db.CameraList))

	db.LoadClientsApi()
	db.LoadClientDataList()
	db.LoadClientlist()
	log.Debug().Msgf("API Clients Length %v -- %v", len(db.ClientList), len(db.ClientList))

	db.CamStartup()
	log.Debug().Msgf("Camera Fetched Data List %v , Length %v", db.CamList, len(db.CamList))

	db.LoadSignlist()
	log.Debug().Msgf("Signs Data List %v , Length %v", db.SignList, len(db.SignList))

}

func addDefaultAdminUser() {
	ctx := context.Background()
	adminUserName := "admin"

	exists, err := db.UserExists(ctx, adminUserName)
	if err != nil {
		log.Error().Err(err).Msg("Error checking admin user existence")
		return
	}

	if !exists {
		defaultAdmin := &db.User{
			UserName:  config.Configvar.AdminUser.Username,
			Password:  config.Configvar.AdminUser.Password,
			FirstName: "System",
			LastName:  "Administrator",
			Role:      "Super-Admin",
		}

		if err := db.AddUser(ctx, defaultAdmin); err != nil {
			log.Error().Err(err).Msg("Failed to add default admin user")
		} else {
			log.Info().Str("Username", config.Configvar.AdminUser.Username).Msg("Default admin user added successfully")
		}
	} else {
		log.Info().Str("Username", config.Configvar.AdminUser.Username).Msg("Default admin user already exists")
	}
}

func addDefaultSettingsData() {
	ctx := context.Background()
	carpark_id := 7077

	exists, err := db.SettingsExists(ctx, carpark_id)
	if err != nil {
		log.Error().Err(err).Msg("Error checking Settings existence")
		return
	}

	if !exists {
		defaultSettings := db.Settings{
			CarParkID:    7077,
			CarParkName:  map[string]interface{}{"en": "Fyc Car Park", "ar": "موقف سيارات Fyc"},
			DefaultLang:  "en",
			PkaImageSize: "small",

			FycCleanCron:      00,
			CountingCleanCron: 00,

			IsFycEnabled:      true,
			IsCountingEnabled: true,

			TimeOutScreenKisok: 10000,
			AppLogo:            "data:image/png;base64,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",
			TC:                 "My Terms and Conditions",
		}

		if err := db.CreateSettings(ctx, &defaultSettings); err != nil {
			log.Error().Err(err).Msg("Failed to add default Settings to database")
		} else {
			log.Info().Int("CarPark ID", carpark_id).Msg("Default Settings added successfully")
		}
	} else {
		log.Info().Int("CarPark ID", carpark_id).Msg("Default Settings already exists")
	}

}
