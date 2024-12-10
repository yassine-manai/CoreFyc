package functions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ValidateAndConvertID(c *gin.Context, idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return 0, fmt.Errorf("invalid ID format: %w", err)
	}
	return id, nil
}
