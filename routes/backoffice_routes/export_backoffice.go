package backoffice_routes

import (
	"fyc/pkg/export"

	"github.com/gin-gonic/gin"
)

func ExportBackoffice(router *gin.RouterGroup) {

	// Export routes
	router.POST("/backoffice/export_zone", export.ExportZone)
	router.POST("/backoffice/export_client", export.ExportClient)
	router.POST("/backoffice/export_camera", export.ExportCamera)
	router.POST("/backoffice/export_sign", export.ExportSign)
	router.POST("/backoffice/export_cars", export.ExportCars)

}
