package api

import (
	"golang-auth/db"
	"golang-auth/types"
	"golang-auth/utils"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllUsers retrieves all users from the database.
func GetAllUsers(c *fiber.Ctx, store *db.Store) error {
	users, err := store.User.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.CreateErrorResponse("Error fetching users", http.StatusInternalServerError, nil))
	}

	// For each user, fetch their notes
	for _, user := range users {
		notes, err := store.Notes.List(c.Context(), user.Id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.CreateErrorResponse("Error fetching notes for user", http.StatusInternalServerError, nil))
		}
		user.Notes = notes // Append notes to the user object
	}

	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Users retrieved successfully", fiber.StatusOK, users))
}

// GetSingleUser retrieves a user by ID from the database.

func GetLoggedInUser(c *fiber.Ctx, store *db.Store) error {

	idParam := c.Locals("userId").(string)
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}
	return CommonUserGet(c, store, id)
}

func UpdateLoggedInUser(c *fiber.Ctx, store *db.Store) error {

	idParam := c.Locals("userId").(string)
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}
	return CommmonUserUpdate(c, store, id)
}

// GetSingleUser retrieves a user by ID from the database.
func GetSingleUser(c *fiber.Ctx, store *db.Store) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}

	return CommonUserGet(c, store, id)
}

func CreateUser(c *fiber.Ctx, store *db.Store) error {
	var user types.UserRequest

	if err := c.BodyParser(&user); err != nil {
		apiError := types.ErrBadRequest("Error parsing request body")
		return c.Status(apiError.Code).JSON(apiError)
	}
	if user.Name == "" || user.Email == "" || user.Password == "" {
		apiError := types.ErrBadRequest("Name, email, and password are required fields.")
		return c.Status(apiError.Code).JSON(apiError)
	}

	_, err := store.User.FindByEmail(user.Email)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already in use"})
	} else if err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check email"})
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		apiError := types.NewError(fiber.StatusInternalServerError, "Error hashing password")
		return c.Status(apiError.Code).JSON(apiError)
	}
	createUser := types.UserCreate{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	// Create user in the database
	newUser, err := store.User.Create(c.Context(), &createUser)
	if err != nil {
		apiError := types.NewError(fiber.StatusInternalServerError, "Error creating user")
		return c.Status(apiError.Code).JSON(apiError)
	}

	// Generate JWT token for the new user
	token, err := utils.GenerateJWT(newUser.Id.Hex(), newUser.Email, "user")
	if err != nil {
		apiError := types.NewError(fiber.StatusInternalServerError, "Error generating JWT token")
		return c.Status(apiError.Code).JSON(apiError)
	}

	// Return the success response with the token
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Signup successful",
		"token":   token,
	})
}

// DeleteUser deletes a user by ID from the database.
func DeleteUser(c *fiber.Ctx, store *db.Store) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}

	deletedUser, err := store.User.Delete(c.Context(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			apiError := types.ErrResourceNotFound("User")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error deleting user")
		return c.Status(apiError.Code).JSON(apiError)
	}
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("User deleted successfully", fiber.StatusOK, deletedUser))
}

func UpdateUser(c *fiber.Ctx, store *db.Store) error {

	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}

	return CommmonUserUpdate(c, store, id)
}

func CommonUserGet(c *fiber.Ctx, store *db.Store, id primitive.ObjectID) error {
	user, err := store.User.Get(c.Context(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			apiError := types.ErrResourceNotFound("User")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error retrieving user")
		return c.Status(apiError.Code).JSON(apiError)
	}

	notes, err := store.Notes.List(c.Context(), id)
	tasks, err := store.Tasks.List(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.CreateErrorResponse("Error fetching notes", http.StatusInternalServerError, nil))
	}
	user.Notes = notes
	user.Tasks = tasks
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("User retrieved successfully", fiber.StatusOK, user))
}

func CommmonUserUpdate(c *fiber.Ctx, store *db.Store, id primitive.ObjectID) error {
	var updatedUser types.UserUpdate
	if err := c.BodyParser(&updatedUser); err != nil {
		apiError := types.ErrBadRequest("Error parsing request body")
		return c.Status(apiError.Code).JSON(apiError)
	}

	existingUser, err := store.User.Get(c.Context(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			apiError := types.ErrResourceNotFound("User")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error retrieving user")
		return c.Status(apiError.Code).JSON(apiError)
	}
	modifiedUser := types.UserUpdate{
		Name:           existingUser.Name,
		Email:          existingUser.Email,
		ProfilePicture: existingUser.ProfilePicture,
		SocialMedia:    existingUser.SocialMedia,
	}
	if updatedUser.Name != "" {
		modifiedUser.Name = updatedUser.Name
	}
	if updatedUser.Email != "" {
		modifiedUser.Email = updatedUser.Email
	}
	if updatedUser.ProfilePicture != "" {
		modifiedUser.ProfilePicture = updatedUser.ProfilePicture
	}
	if (updatedUser.SocialMedia != types.SocialMedia{}) {
		modifiedUser.SocialMedia = updatedUser.SocialMedia
	}

	updatedUserResult, err := store.User.Update(c.Context(), id, &modifiedUser)
	if err != nil {
		if err.Error() == "no user found" {
			apiError := types.ErrResourceNotFound("User")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error updating user")
		return c.Status(apiError.Code).JSON(apiError)
	}

	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("User updated successfully", fiber.StatusOK, updatedUserResult))
}

func GetAllAvatar(c *fiber.Ctx, store *db.Store) error {

	// Define the directory path
	dirPath := "uploads/avatar"

	// Read the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		apiError := types.NewError(fiber.StatusInternalServerError, "Unable to read directory")
		return c.Status(apiError.Code).JSON(apiError)
	}

	// Prepare a list to store file names
	var fileList []string

	// Loop through all files and add to the list
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	// Return the list of files as a JSON response
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("User updated successfully", fiber.StatusOK, fileList))
}
