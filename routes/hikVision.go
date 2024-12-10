package routes

import (
	"fyc/pkg/hikvision"

	"github.com/gin-gonic/gin"
)

func HikVisionRoutes(r *gin.Engine) {
	r.POST("/cam", hikvision.HikvisionHandler)
}
