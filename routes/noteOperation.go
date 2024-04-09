package routes

import (
	"notes-api/handlers"
	"notes-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// Notes CRUD routes
func SetNoteRoutes(app *fiber.App) {
	noteGroup := app.Group("/notes")

	noteGroup.Post("/", limiter.New(middleware.NotesRateConfig), middleware.AuthHandler, handlers.CreateNoteHandler) 	//create a note
	noteGroup.Get("/", limiter.New(middleware.NotesRateConfig), middleware.AuthHandler, handlers.GetAllNotes) 			//view all notes for a user
	noteGroup.Get("/:title", limiter.New(middleware.NotesRateConfig), middleware.AuthHandler, handlers.FindNote) 		//find note by title (within a user scope)
	noteGroup.Put("/:id", limiter.New(middleware.NotesRateConfig), middleware.AuthHandler, handlers.UpdateNote) 			//update notes by id (owned by user)
	noteGroup.Delete("/:id", limiter.New(middleware.NotesRateConfig), middleware.AuthHandler, handlers.DeleteNote)		//delete notes by id (owned by user)

}
