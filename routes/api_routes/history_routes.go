package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func HistoryRoutes(r *gin.Engine) {
	r.GET("/fyc/history", api.GetHistoryAPI)
	r.GET("/fyc/history/:lpn", api.GetHistoryByLPNAPI)
	r.POST("/fyc/history", api.CreateHistoryAPI)
	r.PUT("/fyc/history/:id", api.UpdateHistoryAPI)
	r.DELETE("/fyc/history/:id", api.DeleteHistoryAPI)
}
