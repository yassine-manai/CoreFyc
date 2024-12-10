package third_party_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/third_party"
)

func ThirdPartyToken(r *gin.Engine) {
	r.POST("/token", third_party.GetToken)
	//r.POST("/fyc/v1/Auth/token", third_party.TokenHandler)
}
