package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	cnf "fyc/config"
	"fyc/middleware"
	"fyc/routes/api_routes"
	"fyc/routes/backoffice_routes"
	"fyc/routes/third_party_routes"
)

func SetupRouter() *gin.Engine {
	ModeGin := cnf.Configvar.Server.GinReleaseMode
	if ModeGin == "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	//router.Use(config.CustomErrorHandler())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	log.Debug().Msg("--------------------------  START ROUTING  ----------------------")
	authorizedBackOffice := router.Group("/")
	authorizedBackOffice.Use(middleware.TokenMiddlewareBackOffice())

	authorzedThirdParty := router.Group("/")
	authorzedThirdParty.Use(middleware.TokenMiddlewareThirdParty())

	//rout := router.Group("/api")

	HikVisionRoutes(router) // CAMERA ROUTES ------------------------------------

	third_party_routes.ThirdPartyToken(router)               // Token Generator FOR THIRD PARTY ------------------
	third_party_routes.ThirdPartyRoutes(authorzedThirdParty) // THIRD PARTY ROUTES -------------------------------

	backoffice_routes.BackOfficeToken(router)                // Token Generator FOR BACKOFFICE ------------------
	backoffice_routes.BackOfficeRouter(authorizedBackOffice) // BACKOFFICE ROUTES --------------------------------
	backoffice_routes.ExportBackoffice(authorizedBackOffice) // BACKOFFICE EXPORT ROUTES ---------------------------

	PkaRoutes(router) // PKA ROUTES --------------------------------

	api_routes.DebugRoutes(router) // DEBUG ROUTES -------------------------------------
	api_routes.ZoneImageRoutes(router)
	api_routes.CameraRoutes(router)
	api_routes.ZoneRoutes(router)
	api_routes.PresentCarRoutes(router)
	api_routes.CarDetailRoutes(router)
	api_routes.ClientCredsRoutes(router)
	api_routes.SignRoutes(router)
	api_routes.UserAuditRoutes(router)
	api_routes.UserRoutes(router)
	api_routes.SettingsRoutes(router)
	api_routes.HistoryRoutes(router)
	api_routes.ErrorRoutes(router)

	//log.Debug().Msg("--------------------------  END ROUTING  ---------------------- ")

	// Swagger Endpoint
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/docs/index.html")
	})
	router.GET("/docs/*.any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
