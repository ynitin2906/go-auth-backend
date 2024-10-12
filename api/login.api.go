package api

import (
	"golang-auth/db"
	"golang-auth/types"
	"golang-auth/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx, store *db.Store) error {
	var loginRequest types.Login

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	// Fetch the user by email
	user, err := store.User.FindByEmail(loginRequest.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}
	//compare has passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}
	//generate jwt token
	token, err := utils.GenerateJWT(user.Id.Hex(), user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate JWT token",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Login successfully",
		"token":   token,
	})
}
