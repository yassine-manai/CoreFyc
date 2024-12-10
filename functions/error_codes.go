package functions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseHandler struct{}

func NewResponseHandler() *ResponseHandler {
	return &ResponseHandler{}
}

func (r *ResponseHandler) RespondWithStatus(c *gin.Context, code int) {
	var message string
	var httpStatus int
	var success string

	switch code {
	case 200:
		success = "true"
		message = "OK"
		httpStatus = http.StatusOK
	case -1:
		success = "false"
		message = "Invalid credentials"
		httpStatus = http.StatusUnauthorized
	case -2:
		success = "false"
		message = "User is disabled"
		httpStatus = http.StatusForbidden
	case -3:
		success = "false"
		message = "Unauthorized, you need to connect first!"
		httpStatus = http.StatusUnauthorized
	case -4:
		success = "false"
		message = "No data found!"
		httpStatus = http.StatusNotFound
	case -5:
		success = "false"
		message = "Invalid request. 'id' parameter is required."
		httpStatus = http.StatusBadRequest
	case -6:
		success = "false"
		message = "Zone ID already exists!"
		httpStatus = http.StatusConflict
	case -7:
		success = "false"
		message = "Zone ID not found!"
		httpStatus = http.StatusNotFound
	case -8:
		success = "false"
		message = "Camera ID already exists!"
		httpStatus = http.StatusConflict
	case -9:
		success = "false"
		message = "Camera ID not found!"
		httpStatus = http.StatusNotFound
	case -10:
		success = "false"
		message = "Sign ID already exists!"
		httpStatus = http.StatusConflict
	case -11:
		success = "false"
		message = "Sign ID not found!"
		httpStatus = http.StatusNotFound
	case -12:
		success = "false"
		message = "Client ID already exists!"
		httpStatus = http.StatusConflict
	case -13:
		success = "false"
		message = "Client ID not found!"
		httpStatus = http.StatusNotFound
	case -500:
		success = "false"
		message = "An unexpected error occurred. Please try again later."
		httpStatus = http.StatusInternalServerError
	default:
		success = "false"
		message = "Unknown error"
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, gin.H{
		"success": success,
		"error":   message,
		"code":    code,
	})
}
