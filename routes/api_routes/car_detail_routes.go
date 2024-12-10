package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func CarDetailRoutes(r *gin.Engine) {
	r.GET("/fyc/carDetails", api.GetCarDetailsAPI)
	r.POST("/fyc/carDetails", api.CreateCarDetailAPI)
	r.PUT("/fyc/carDetails", api.UpdateCarDetailByIdAPI)
	r.DELETE("/fyc/carDetails", api.DeleteCarDetailAPI)
}
