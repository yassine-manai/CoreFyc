package third_party_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/third_party"
)

// @title			Main API
// @version		1.0
// @description	API documentation for main endpoints.
func ThirdPartyRoutes(r *gin.RouterGroup) {
	//r.POST("/token", third_party.GetToken)
	r.GET("/findmycar", third_party.FindMyCar)
	r.GET("/getpicture", third_party.GetPicture)
	r.GET("/getsettings", third_party.Getsettings)
	//r.POST("/fyc/v1/Auth/token", third_party.TokenHandler)
}
