package main

import (
	"golang-auth/api"
	"golang-auth/db"
	"golang-auth/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the application routes
func SetupRoutes(app *fiber.App, store *db.Store) {

	setupAuthRoutes(app, store)
	app.Use(middleware.AuthMiddleware)

	setupLoggedInUserRoutes(app, store)
	setupNoteRoutes(app, store)
	setupTasksRoutes(app, store)

	// app.Use(middleware.AdminMiddleware)
	setupAdminRoutes(app, store)

}

// Setup Authentication routes
func setupAuthRoutes(app *fiber.App, store *db.Store) {
	app.Static("/avatar", "./uploads/avatar")

	app.Post("/login", func(c *fiber.Ctx) error {
		return api.Login(c, store)
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		return api.CreateUser(c, store)
	})
}
func setupAdminRoutes(app *fiber.App, store *db.Store) {
	app.Use(middleware.AdminMiddleware)

	app.Get("/users/all", func(c *fiber.Ctx) error {
		return api.GetAllUsers(c, store)
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		return api.DeleteUser(c, store)
	})
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return api.GetSingleUser(c, store)
	})

	app.Patch("/users/:id", func(c *fiber.Ctx) error {
		return api.UpdateUser(c, store)
	})
}

func setupLoggedInUserRoutes(app *fiber.App, store *db.Store) {

	app.Get("/loggedinuser", func(c *fiber.Ctx) error {
		return api.GetLoggedInUser(c, store)
	})

	app.Patch("/loggedinuser", func(c *fiber.Ctx) error {
		return api.UpdateLoggedInUser(c, store)
	})
	app.Get("/allavatar", func(c *fiber.Ctx) error {
		return api.GetAllAvatar(c, store)
	})
}

func setupNoteRoutes(app *fiber.App, store *db.Store) {

	app.Get("/notes/user", func(c *fiber.Ctx) error {
		return api.GetAllNotesForLoginUser(c, store)
	})
	app.Get("/notes/user/:id", func(c *fiber.Ctx) error {
		return api.GetAllNotesForUserById(c, store)
	})

	app.Get("/notes/:id", func(c *fiber.Ctx) error {
		return api.GetSingleNote(c, store)
	})

	app.Post("/notes", func(c *fiber.Ctx) error {
		return api.CreateNote(c, store)
	})

	app.Patch("/notes/:id", func(c *fiber.Ctx) error {
		return api.UpdateNote(c, store)
	})

	app.Delete("/notes/:id", func(c *fiber.Ctx) error {
		return api.DeleteNote(c, store)
	})

}
func setupTasksRoutes(app *fiber.App, store *db.Store) {

	app.Get("/tasks/user", func(c *fiber.Ctx) error {
		return api.GetAllTasksForLoginUser(c, store)
	})
	app.Get("/tasks/user/:id", func(c *fiber.Ctx) error {
		return api.GetAllTasksForUserById(c, store)
	})

	app.Get("/task/:id", func(c *fiber.Ctx) error {
		return api.GetSingleTask(c, store)
	})

	app.Post("/tasks", func(c *fiber.Ctx) error {
		return api.CreateTask(c, store)
	})

	app.Patch("/tasks/:id", func(c *fiber.Ctx) error {
		return api.UpdateTask(c, store)
	})

	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		return api.DeleteTask(c, store)
	})

}
