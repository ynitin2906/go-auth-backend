package api

import (
	"fmt"
	"golang-auth/db"
	"golang-auth/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllNotesForUser retrieves all notes for a specific user
func GetAllNotesForLoginUser(c *fiber.Ctx, store *db.Store) error {
	// userId := c.Locals("userId").(primitive.ObjectID)

	userId, err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	notes, err := store.Notes.List(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.CreateErrorResponse("Error fetching notes", http.StatusInternalServerError, nil))
	}
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Notes retrieved successfully", fiber.StatusOK, notes))
}
func GetAllNotesForUserById(c *fiber.Ctx, store *db.Store) error {
	idParam := c.Params("id")
	userId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}

	notes, err := store.Notes.List(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.CreateErrorResponse("Error fetching notes", http.StatusInternalServerError, nil))
	}
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Notes retrieved successfully", fiber.StatusOK, notes))
}

func GetSingleNote(c *fiber.Ctx, store *db.Store) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}

	note, err := store.Notes.Get(c.Context(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			apiError := types.ErrResourceNotFound("Note")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error retrieving note")
		return c.Status(apiError.Code).JSON(apiError)
	}
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Note retrieved successfully", fiber.StatusOK, note))
}

func CreateNote(c *fiber.Ctx, store *db.Store) error {
	var note types.NotesRequest
	if err := c.BodyParser(&note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse("Invalid request body", http.StatusBadRequest, nil))
	}

	userId, err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	createNote := types.NotesCreate{
		Title:    note.Title,
		Category: note.Category,
		Note:     note.Note,
		UserID:   userId,
	}
	newNote, err := store.Notes.Create(c.Context(), &createNote)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.CreateErrorResponse("Error creating note", http.StatusInternalServerError, nil))
	}
	return c.Status(fiber.StatusCreated).JSON(types.CreateSuccessResponse("Note created successfully", fiber.StatusCreated, newNote))
}

func DeleteNote(c *fiber.Ctx, store *db.Store) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}

	// Check authorization
	_, err = CheckNoteAuthorization(c, store, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse(err.Error(), http.StatusBadRequest, nil))
	}

	// deletedNote, err := store.Notes.Delete(c.Context(), id)
	_, err = store.Notes.Delete(c.Context(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			apiError := types.ErrResourceNotFound("Note")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error deleting note")
		return c.Status(apiError.Code).JSON(apiError)
	}
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Note deleted successfully", fiber.StatusOK, nil))
}

func UpdateNote(c *fiber.Ctx, store *db.Store) error {
	// Parse note ID from request parameters
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse("Invalid ID", http.StatusBadRequest, nil))
	}

	// Parse the updated note from the request body
	var updatedNote types.NotesUpdate
	if err := c.BodyParser(&updatedNote); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse("Invalid request body", http.StatusBadRequest, nil))
	}

	// Check authorization and retrieve the existing note
	existingNote, err := CheckNoteAuthorization(c, store, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse(err.Error(), http.StatusBadRequest, nil))
	}

	// Merge existing note with updates (if fields are provided)
	modifiedNote := types.NotesUpdate{
		Title:    existingNote.Title,
		Category: existingNote.Category,
		Note:     existingNote.Note,
	}

	if updatedNote.Title != "" {
		modifiedNote.Title = updatedNote.Title
	}
	if updatedNote.Category != "" {
		modifiedNote.Category = updatedNote.Category
	}
	if updatedNote.Note != "" {
		modifiedNote.Note = updatedNote.Note
	}

	// Update the note in the database
	updatedNoteResult, err := store.Notes.Update(c.Context(), id, &modifiedNote)
	if err != nil {
		if err.Error() == "no note found" {
			apiError := types.ErrResourceNotFound("Note")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error updating note")
		return c.Status(apiError.Code).JSON(apiError)
	}

	// Return the updated note as the response
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Note updated successfully", fiber.StatusOK, updatedNoteResult))
}

// for checking purpose  that only logged in user or user who's role is
// admin can only delete not other user like if he using postman or swagger

func CheckNoteAuthorization(c *fiber.Ctx, store *db.Store, noteId primitive.ObjectID) (*types.Notes, error) {
	// Retrieve the logged-in user's ID from the request context
	// userId, ok := c.Locals("userId").(primitive.ObjectID)
	// if !ok {
	// 	apiError := types.ErrUnAuthorized()
	// 	return nil, c.Status(apiError.Code).JSON(apiError)
	// }

	userId, err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	// Get the note by ID
	note, err := store.Notes.Get(c.Context(), noteId)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("error retrieving note: %w", err)
	}

	// Check if the logged-in user is the owner of the note or an admin
	if c.Locals("role") != "admin" && note.UserID != userId {
		return nil, fmt.Errorf("unauthorized access to the note")
	}

	// If all checks pass, return the note without an error
	return note, nil

}
