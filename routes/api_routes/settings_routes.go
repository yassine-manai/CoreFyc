package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func SettingsRoutes(r *gin.Engine) {
	r.GET("/fyc/settings", api.GetSettingsAPI)
	r.POST("/fyc/settings", api.AddSettingsAPI)
	r.PUT("/fyc/settings", api.UpdateSettingsAPI)
}
