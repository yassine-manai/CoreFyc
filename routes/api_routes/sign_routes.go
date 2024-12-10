package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func SignRoutes(r *gin.Engine) {
	r.GET("/fyc/sign", api.GetSignAPI)
	r.POST("/fyc/sign", api.CreateSignAPI)
	r.PUT("/fyc/sign", api.UpdateSignAPI)
	r.DELETE("/fyc/sign", api.DeleteSignAPI)

	//r.GET("/fyc/signEnabled", api.GetSignEnabledAPI)
	//r.GET("/fyc/signDeleted", api.GetSignDeletedAPI)
	//r.PUT("/fyc/signState", api.ChangeSigntateAPI)
}
