package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func UserAuditRoutes(r *gin.Engine) {
	r.GET("/fyc/UserAudit", api.GetUserAuditAPI)
	r.POST("/fyc/UserAudit", api.CreateUserAuditAPI)
	r.PUT("/fyc/UserAudit", api.UpdateUserAuditAPI)
	r.DELETE("/fyc/UserAudit", api.DeleteUserAuditAPI)
}
