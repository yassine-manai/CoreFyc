package backoffice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/pkg/db"
)

// DebugBackOfficeAPI godoc
//
//	@Summary	Debug BackOffice API
//	@Tags		BackOffice - Debug
//	@Security	BearerAuthBackOffice
//	@Produce	json
//	@Router		/backoffice/debug [get]
func Debuger_BackOffice(c *gin.Context) {

	log.Debug().Msg(" -- -- -- -- -- Debugging BackOffice API -- -- -- -- --")
	c.JSON(http.StatusOK, gin.H{
		"ZoneList":   db.Zonelist,
		"CameraList": db.CameraList,
		"CamList":    db.CamList,
		"SignsList":  db.SignList,
		"ClientList": db.ClientDataList,
		"Clients":    db.ClientListAPI,
	})

}
