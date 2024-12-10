package backoffice_routes

import (
	"fyc/pkg/backoffice"

	"github.com/gin-gonic/gin"
)

func BackOfficeToken(r *gin.Engine) {
	// LOGIN
	r.POST("/backoffice/login", backoffice.Login)
}
