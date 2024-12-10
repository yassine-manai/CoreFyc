package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func PresentCarRoutes(r *gin.Engine) {
	r.GET("/fyc/presentcars", api.GetPresentCarsAPI)
	r.POST("/fyc/presentcars", api.CreatePresentCarAPI)
	r.PUT("/fyc/presentcars", api.UpdatePresentCarBylpnAPI)
	r.PUT("/fyc/presentcars/:lpn", api.UpdatePresentCarByIdAPI)
	r.DELETE("/fyc/presentcars/:id", api.DeletePresentCarAPI)

	//r.GET("/fyc/presentcars/:lpn", api.GetPresentCarByLPNAPI)

}
