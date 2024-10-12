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
func GetAllTasksForLoginUser(c *fiber.Ctx, store *db.Store) error {
	userId, err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	tasks, err := store.Tasks.List(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.CreateErrorResponse("Error fetching tasks", http.StatusInternalServerError, nil))
	}
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Tasks retrieved successfully", fiber.StatusOK, tasks))
}

func GetAllTasksForUserById(c *fiber.Ctx, store *db.Store) error {
	idParam := c.Params("id")
	userId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}

	tasks, err := store.Tasks.List(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.CreateErrorResponse("Error fetching tasks", http.StatusInternalServerError, nil))
	}
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Tasks retrieved successfully", fiber.StatusOK, tasks))
}

func GetSingleTask(c *fiber.Ctx, store *db.Store) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}

	task, err := store.Tasks.Get(c.Context(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			apiError := types.ErrResourceNotFound("Task")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error retrieving task")
		return c.Status(apiError.Code).JSON(apiError)
	}
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Task retrieved successfully", fiber.StatusOK, task))
}

func CreateTask(c *fiber.Ctx, store *db.Store) error {
	// Parse the request body into TasksRequest
	var task types.TasksRequest
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse("Invalid request body", http.StatusBadRequest, nil))
	}

	// Retrieve the user ID from the context (logged-in user)
	userId, err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Create initial status history from the task status and logged-in user
	statusEntry := &types.Status{
		Status: task.Status,
		UserId: userId.Hex(),
	}

	// Prepare the task creation struct
	createTask := types.TasksCreate{
		Title:         task.Title,
		Category:      task.Category,
		Task:          task.Task,
		UserID:        userId,
		StatusHistory: []*types.Status{statusEntry}, // Add status entry to the history
	}

	// Call the DB function to create the task
	newTask, err := store.Tasks.Create(c.Context(), &createTask)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.CreateErrorResponse("Error creating task", http.StatusInternalServerError, nil))
	}

	// Return success response with the created task
	return c.Status(fiber.StatusCreated).JSON(types.CreateSuccessResponse("Task created successfully", fiber.StatusCreated, newTask))
}
func UpdateTask(c *fiber.Ctx, store *db.Store) error {
	// Parse task ID from the request parameters
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse("Invalid ID", http.StatusBadRequest, nil))
	}

	// Parse the updated task from the request body
	var updatedTask types.TasksRequest
	if err := c.BodyParser(&updatedTask); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse("Invalid request body", http.StatusBadRequest, nil))
	}

	// Check authorization and retrieve the existing task
	existingTask, err := CheckTaskAuthorization(c, store, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse(err.Error(), http.StatusBadRequest, nil))
	}

	// Merge existing task with updates
	modifiedTask := types.TasksUpdate{
		Title:         existingTask.Title,
		Category:      existingTask.Category,
		Task:          existingTask.Task,
		StatusHistory: existingTask.StatusHistory, // Preserve existing status history
	}

	// Update task details if provided
	if updatedTask.Title != "" {
		modifiedTask.Title = updatedTask.Title
	}
	if updatedTask.Category != "" {
		modifiedTask.Category = updatedTask.Category
	}
	if updatedTask.Task != "" {
		modifiedTask.Task = updatedTask.Task
	}

	// Append new status to the status history
	userId, err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse("Invalid user ID format", http.StatusBadRequest, nil))
	}
	newStatus := &types.Status{
		Status: updatedTask.Status,
		UserId: userId.Hex(),
	}
	modifiedTask.StatusHistory = append(modifiedTask.StatusHistory, newStatus)

	// Update the task in the database
	updatedTaskResult, err := store.Tasks.Update(c.Context(), id, &modifiedTask)
	if err != nil {
		if err.Error() == "no task found" {
			apiError := types.ErrResourceNotFound("Task")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error updating task")
		return c.Status(apiError.Code).JSON(apiError)
	}

	// Return the updated task as the response
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Task updated successfully", fiber.StatusOK, updatedTaskResult))
}

func DeleteTask(c *fiber.Ctx, store *db.Store) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiError := types.ErrInvalidID()
		return c.Status(apiError.Code).JSON(apiError)
	}

	// Check authorization
	_, err = CheckTaskAuthorization(c, store, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.CreateErrorResponse(err.Error(), http.StatusBadRequest, nil))
	}

	_, err = store.Tasks.Delete(c.Context(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			apiError := types.ErrResourceNotFound("Task")
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := types.NewError(fiber.StatusInternalServerError, "Error deleting task")
		return c.Status(apiError.Code).JSON(apiError)
	}
	return c.Status(fiber.StatusOK).JSON(types.CreateSuccessResponse("Task deleted successfully", fiber.StatusOK, nil))
}

// for checking purpose  that only logged in user or user who's role is
// admin can only delete not other user like if he using postman or swagger

func CheckTaskAuthorization(c *fiber.Ctx, store *db.Store, noteId primitive.ObjectID) (*types.Tasks, error) {
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
	task, err := store.Tasks.Get(c.Context(), noteId)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("error retrieving note: %w", err)
	}

	// Check if the logged-in user is the owner of the task or an admin
	if c.Locals("role") != "admin" && task.UserID != userId {
		return nil, fmt.Errorf("unauthorized access to the task")
	}

	// If all checks pass, return the task without an error
	return task, nil

}
