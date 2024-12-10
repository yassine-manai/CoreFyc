package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func ZoneImageRoutes(r *gin.Engine) {
	r.GET("/fyc/zonesImages", api.GetAllImageZonesAPI)
	//r.GET("/fyc/zonesImage", api.GetZoneImageByIDAPI)
	r.POST("/fyc/zonesImage", api.CreateZoneImageAPI)
	r.PUT("/fyc/zonesImage/:id", api.UpdateZoneImageByIdAPI)
	r.DELETE("/fyc/zonesImage/:id", api.DeleteZoneImageAPI)
}
