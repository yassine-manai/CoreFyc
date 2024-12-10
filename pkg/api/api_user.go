package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fyc/pkg/db"
)

// GetUsersAPI godoc
//
//	@Summary		Get Users or a specific User by username
//	@Description	Get a list of users (all, enabled) or a specific user by username
//	@Tags			Users
//	@Produce		json
//	@Param			username	query		string					false	"Username"
//	@Param			type		query		string					false	"Type of user: all or enabled (default is 'all')"
//	@Success		200			{array}		db.User					"List of Users or a single User"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Failure		404			{object}	map[string]interface{}	"No user found"
//	@Router			/fyc/users [get]
func GetUsersAPI(c *gin.Context) {
	log.Debug().Msg("GetUsersAPI request")
	ctx := context.Background()
	username := c.Query("username")

	// Check if username is provided to fetch a specific user
	if username != "" {
		var user interface{}
		var err error

		user, err = db.GetUserByUsername(ctx, username)

		if err != nil {
			log.Err(err).Str("username", username).Msg("Error retrieving User by username")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "User not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("username", username).Msg("User fetched successfully")
		c.JSON(http.StatusOK, user)
		return
	}

	// Fetch list of users based on type
	var users []db.User
	var err error

	users, err = db.GetAllUsers(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Error retrieving users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving users",
			"code":    10,
		})
		return
	}

	if len(users) == 0 {
		log.Warn().Msg("No users found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No users found",
			"code":    9,
		})
		return
	}

	log.Info().Int("user_count", len(users)).Msg("Users fetched successfully")
	c.JSON(http.StatusOK, users)
}

// AddUserCred godoc
//
//	@Summary		Add a new User
//	@Description	Add a new User to the database
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			User	body		db.User	true	"User data"
//	@Success		201		{object}	db.User
//	@Router			/fyc/user [post]
func AddUserAPI(c *gin.Context) {
	var user db.User

	log.Info().Msg("Attempting to add new user")

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for user creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	if err := db.AddUser(ctx, &user); err != nil {
		log.Error().Err(err).Msg("Error creating user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Str("UserName", user.UserName).Msg("UserName created successfully")
	c.JSON(http.StatusCreated, user)
}

// UpdateClientCred godoc
//
//	@Summary		Update a client credential
//	@Description	Update an existing User by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			username	query		string	true	"Client ID"
//	@Param			clientCred	body		db.User	true	"Updated client credential data"
//	@Success		200			{object}	db.User
//	@Router			/fyc/user [put]
func UpdateUserAPI(c *gin.Context) {
	usernameStr := c.Query("username")

	log.Info().Str("username", usernameStr).Msg("Attempting to update User")

	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for User update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": "Invalid client credential data",
			"code":    12,
		})
		return
	}

	if user.UserName != usernameStr {
		log.Warn().Str("username param", usernameStr).Str("username", user.UserName).Msg("ID mismatch between path and body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The Username in the request body does not match the Username in the query parameter",
			"code":    13,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := db.UpdateUser(ctx, usernameStr, &user)
	if err != nil {
		log.Error().Err(err).Str("client_id", usernameStr).Msg("Error updating USER")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update User",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("username", usernameStr).Msg("No user found to update")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No  user found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("username", usernameStr).Msg("User updated successfully")
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a Use by username
//	@Tags			Users
//	@Param			username	query		string	true	"Username"
//	@Success		200			{string}	string	"User deleted successfully"
//	@Router			/fyc/user [delete]
func DeleteUserCredAPI(c *gin.Context) {
	userStr := c.Query("username")
	log.Info().Str("User", userStr).Msg("Attempting to delete user")
	ctx := context.Background()
	rowsAffected, err := db.DeleteUser(ctx, userStr)
	if err != nil {
		log.Error().Err(err).Str("user", userStr).Msg("Error deleting user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("username", userStr).Msg("No User found to delete")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No user found with the specified Username",
			"code":    9,
		})
		return
	}

	log.Info().Str("username", userStr).Msg("User deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": "User deleted successfully",
		"code":    8,
	})
}
