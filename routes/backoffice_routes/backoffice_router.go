package backoffice_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/backoffice"
)

func BackOfficeRouter(router *gin.RouterGroup) {

	// Debug routes
	router.GET("/backoffice/debug", backoffice.Debuger_BackOffice)

	// Dashboard routes
	router.GET("/backoffice/get_dashboard_data", backoffice.GetDashboardData)

	// Zones routes
	router.GET("/backoffice/get_zones", backoffice.GetZonesAPI)
	router.GET("/backoffice/get_zones_names", backoffice.GetZonesNames)
	router.POST("/backoffice/add_zone", backoffice.CreateZone)
	router.PUT("/backoffice/update_zone", backoffice.UpdateZoneDataAPI)
	router.DELETE("/backoffice/delete_zone", backoffice.DeleteZoneDataAPI)

	// Camera routes
	router.GET("/backoffice/getCameras", backoffice.GetCameraDataAPI)
	router.POST("/backoffice/addCamera", backoffice.AddCameraDataAPI)
	router.PUT("/backoffice/updateCamera", backoffice.UpdateCameraDataAPI)
	router.DELETE("/backoffice/deleteCameras", backoffice.DeleteCameraDataAPI)

	// Signes routes
	router.GET("/backoffice/getSign", backoffice.GetSignDataAPI)
	router.POST("/backoffice/addSign", backoffice.CreateSignDataAPI)
	router.PUT("/backoffice/updateSign", backoffice.UpdateSignDataAPI)
	router.DELETE("/backoffice/deleteSign", backoffice.DeleteSignDataAPI)

	// Client routes
	router.GET("/backoffice/get_clients", backoffice.GetClients)
	router.POST("/backoffice/addClient", backoffice.AddClientAPI)
	router.PUT("/backoffice/updateClient", backoffice.UpdateClientAPI)
	router.DELETE("/backoffice/deleteClient", backoffice.DeleteClientAPI)

	// Settings routes
	router.GET("/backoffice/getSettings", backoffice.GetSettingsDataAPI)
	router.PUT("/backoffice/updateSettings", backoffice.UpdateSettingsDataAPI)

	// Present car routes
	router.GET("/backoffice/get_present_car", backoffice.GetAllPresentTransactionsDataAPI)
	router.GET("/backoffice/get_present_car_id", backoffice.GetAllPresentTransactionsDataIDAPI)

}
