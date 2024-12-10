package routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/pka"
)

func PkaRoutes(r *gin.Engine) {
	r.GET("/v2/bays.json", pka.PkaSearchAPI)
	r.GET("/v2/maps/:imagename", pka.PkaImageAPI)
}
