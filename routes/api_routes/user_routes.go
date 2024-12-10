package api_routes

import (
	"github.com/gin-gonic/gin"

	"fyc/pkg/api"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/fyc/users", api.GetUsersAPI)
	r.POST("/fyc/user", api.AddUserAPI)
	r.PUT("/fyc/user", api.UpdateUserAPI)
	r.DELETE("/fyc/user", api.DeleteUserCredAPI)

	//r.GET("/fyc/userEnabled", api.GetUserEnabledAPI)
	//r.GET("/fyc/userDeleted", api.GetUserDeletedAPI)
	//r.PUT("/fyc/userState", api.ChangeUserStateAPI)

}
