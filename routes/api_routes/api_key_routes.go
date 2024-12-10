package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func ClientCredsRoutes(r *gin.Engine) {
	r.GET("/fyc/apikey", api.GetAllClientCredsApi)
	r.PUT("/fyc/apikey", api.UpdateClientCredAPI)
	r.POST("/fyc/apikey", api.AddClientCredAPI)
	r.DELETE("/fyc/apikey", api.DeleteClientCredAPI)

	//r.GET("/fyc/clientEnabled", api.GetClientEnabledAPI)
	//r.GET("/fyc/clientsDeleted", api.GetClientDeletedAPI)
	//r.PUT("/fyc/clientState", api.ChangeClientStateAPI)

}
