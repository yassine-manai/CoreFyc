package pka

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestSignAPI(c *gin.Context) {

	// Define the /testSign POST route
	var jsonData map[string]interface{}

	// Bind JSON body to jsonData map
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Respond with the JSON data received in the request body
	fmt.Printf("DATA ------------------ \n %v", jsonData)
	c.JSON(http.StatusOK, jsonData)

}
