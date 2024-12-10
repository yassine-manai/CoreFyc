package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/debug"
)

func DebugRoutes(r *gin.Engine) {

	r.GET("/fyc/debug", debug.Debuger_api)
}
