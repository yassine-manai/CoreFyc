package debug

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fyc/pkg/db"
)

// DebugAPI godoc
//
//	@Summary	Debug API
//	@Tags		Debug
//	@Produce	json
//	@Router		/fyc/debug [get]
func Debuger_api(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"ZoneList":   db.Zonelist,
		"CameraList": db.CameraList,
		"CamList":    db.CamList,
		"SignsList":  db.SignList,
		"ClientList": db.ClientDataList,
		"Clients":    db.ClientListAPI,
	})

}
