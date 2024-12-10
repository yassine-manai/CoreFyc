package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func ErrorRoutes(r *gin.Engine) {
	r.GET("/fyc/errors", api.GetAllErrorCode)
	r.POST("/fyc/errors", api.CreateErrorMessageAPI)
	r.PUT("/fyc/errors", api.UpdateErrorMessageAPI)
	r.DELETE("/fyc/errors", api.DeleteErrorMessageAPI)
}
