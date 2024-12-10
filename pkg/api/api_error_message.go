package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/pkg/db"
)

// GetAllErrorCode godoc
//
//	@Summary		Get all error messages or a specific one
//	@Description	Get a list of all error messages or a specific one by code and language
//	@Tags			Errors
//	@Produce		json
//	@Param			code	query		string	false	"Error code to fetch specific error message"
//	@Param			lang	query		string	false	"Language of the error message"
//	@Success		200		{object}	[]db.ErrorMessage
//	@Router			/fyc/errors [get]
func GetAllErrorCode(c *gin.Context) {
	ctx := context.Background()
	code := c.Query("code")
	lang := c.Query("lang")

	if code != "" {
		codeReq, err := strconv.Atoi(code)
		if err != nil {
			log.Error().Err(err).Msg("Invalid code format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid code format",
				"message": "Code must be a valid integer",
				"code":    400,
			})
			return
		}

		errorMessage, err := db.GetErrorMessageByFilter(ctx, codeReq, lang)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving error message by code")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error retrieving error message",
				"code":    10,
			})
			return
		}

		if errorMessage.Code == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No error message found for the provided code",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, errorMessage)
	} else {
		errorMessages, err := db.GetErrorMessage(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving error messages")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error retrieving error messages",
				"code":    10,
			})
			return
		}

		if len(errorMessages) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No error messages found",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, errorMessages)
	}
}

// CreateErrorMessageAPI godoc
//
//	@Summary		Create a new error message
//	@Description	Create a new error message
//	@Tags			Errors
//	@Accept			json
//	@Produce		json
//	@Param			errMsg	body		db.ErrorMessage	true	"Error message object"
//	@Success		201		{object}	db.ErrorMessage
//	@Router			/fyc/errors [post]
func CreateErrorMessageAPI(c *gin.Context) {
	var errMsg db.ErrorMessage
	if err := c.ShouldBindJSON(&errMsg); err != nil {
		log.Error().Err(err).Msg("Invalid input for error message")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	// Create context
	ctx := context.Background()

	// Check if the error code already exists in the database
	existingErrorMsg, err := db.GetErrorMessageByCode(ctx, errMsg.Code)
	if err == nil && existingErrorMsg.Code != 0 {
		log.Warn().Int("code", errMsg.Code).Msg("Error code already exists")
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Conflict",
			"message": "Error message with this code already exists",
			"code":    409,
		})
		return
	}

	// Convert the keys of Messages to lowercase
	lowercaseMessages := make(map[string]string)
	for lang, msg := range errMsg.Messages {
		lowercaseMessages[strings.ToLower(lang)] = msg
	}
	errMsg.Messages = lowercaseMessages

	// Create the error message in the database
	if err := db.CreateErrorMessage(ctx, &errMsg); err != nil {
		log.Error().Err(err).Msg("Failed to create error message")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Failed to create error message",
			"code":    500,
		})
		return
	}

	log.Info().Int("code", errMsg.Code).Msg("Error Message created successfully")
	c.JSON(http.StatusCreated, errMsg)
}

// UpdateErrorMessageAPI godoc
//
//	@Summary		Update an existing error message
//	@Description	Update an existing error message by code
//	@Tags			Errors
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string			true	"Error message code"
//	@Param			errMsg	body		db.ErrorMessage	true	"Updated error message object"
//	@Success		200		{object}	db.ErrorMessage
//	@Router			/fyc/errors [put]
func UpdateErrorMessageAPI(c *gin.Context) {
	codeStr := c.Query("code")
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid error code format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid code format",
			"message": "Code must be a valid integer",
			"code":    400,
		})
		return
	}

	var errMsg db.ErrorMessage
	if err := c.ShouldBindJSON(&errMsg); err != nil {
		log.Error().Err(err).Msg("Invalid input for error message")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	if errMsg.Code != code {
		log.Info().Msg("The Code in the request body does not match the query Code")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Code mismatch",
			"message": "The Code in the request body does not match the query code",
			"code":    13,
		})
		return
	}

	ctx := context.Background()

	// Check if the error code exists in the database
	existingErrorMsg, err := db.GetErrorMessageByCode(ctx, errMsg.Code)
	if err == nil && existingErrorMsg.Code == 0 {
		log.Warn().Int("code", errMsg.Code).Msg("Error code doesn't exist")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Error message with this code doesn't exist",
			"code":    404,
		})
		return
	}

	// Convert the keys of Messages to lowercase
	lowercaseMessages := make(map[string]string)
	for lang, msg := range errMsg.Messages {
		lowercaseMessages[strings.ToLower(lang)] = msg
	}
	errMsg.Messages = lowercaseMessages

	if err := db.UpdateErrorMessage(ctx, &errMsg); err != nil {
		log.Error().Err(err).Msg("Failed to update error message")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Failed to update error message",
			"code":    500,
		})
		return
	}

	log.Info().Int("code", errMsg.Code).Msg("Error message updated successfully")
	c.JSON(http.StatusOK, errMsg)
}

// DeleteErrorMessageAPI godoc
//
//	@Summary		Delete a specific language from an error message
//	@Description	Delete a specific language entry from the messages field of an error message by code
//	@Tags			Errors
//	@Param			code	query	string	true	"Error message code"
//	@Param			lang	query	string	true	"Language of the error message"
//	@Router			/fyc/errors [delete]
func DeleteErrorMessageAPI(c *gin.Context) {
	codeStr := c.Query("code")
	langQuery := c.Query("lang")

	// Validate the code parameter
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid error code format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid code format",
			"message": "Code must be a valid integer",
			"code":    400,
		})
		return
	}

	// Ensure language parameter is provided
	if langQuery == "" {
		log.Error().Msg("Missing language parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing language parameter",
			"message": "A valid language parameter must be provided",
			"code":    400,
		})
		return
	}

	// Call the DeleteErrorMessage function to remove the specific language entry
	ctx := context.Background()
	rowsAffected, err := db.DeleteErrorMessage(ctx, code, langQuery)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete error message")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete error message",
			"message": err.Error(),
			"code":    500,
		})
		return
	}

	// Handle case where no rows were affected (i.e., no such code or language found)
	if rowsAffected == 0 {
		log.Info().Int("code", code).Str("lang", langQuery).Msg("No error message found with the specified code and language")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No error message found for the provided code and language",
			"code":    404,
		})
		return
	}

	// Success response
	log.Info().Int("code", code).Str("lang", langQuery).Msg("Error message language deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success":      "Deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
