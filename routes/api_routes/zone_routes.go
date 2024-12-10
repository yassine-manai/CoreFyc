package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func ZoneRoutes(r *gin.Engine) {
	r.GET("/fyc/zones", api.GetZonesAPI)
	r.POST("/fyc/zones", api.CreateZoneAPI)
	r.PUT("/fyc/zones", api.UpdateZoneIdAPI)
	r.DELETE("/fyc/zones", api.DeleteZoneAPI)
	//r.GET("/fyc/zoneName", api.GetZoneNameAPI)
	//r.PUT("/fyc/zoneState", api.ChangeZoneStateAPI)
}
