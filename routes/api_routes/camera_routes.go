package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func CameraRoutes(r *gin.Engine) {
	r.GET("/fyc/cameras", api.GetCameraAPI)
	r.POST("/fyc/cameras", api.CreateCameraAPI)
	r.PUT("/fyc/cameras", api.UpdateCameraAPI)
	r.DELETE("/fyc/cameras", api.DeleteCameraAPI)

	//r.GET("/fyc/camerasEnabled", api.GetCameraEnabledAPI)
	//r.GET("/fyc/camerasDeleted", api.GetCameraDeletedAPI)
	//r.PUT("/fyc/cameraState", api.ChangeCameraStateAPI)

}
