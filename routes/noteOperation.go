package routes

import (
	"notes-api/handlers"
	"notes-api/middleware"

	"github.com/gofiber/fiber/v2"
)


// Notes CRUD routes
func SetNoteRoutes(app *fiber.App) {
	noteGroup := app.Group("/notes")

	noteGroup.Post("/", middleware.AuthHandler, handlers.CreateNoteHandler) 	//create a note
	noteGroup.Get("/", middleware.AuthHandler, handlers.GetAllNotes) 			//view all notes for a user
	noteGroup.Get("/:title", middleware.AuthHandler, handlers.FindNote) 		//find note by title (within a user scope)
	noteGroup.Put("/:id", middleware.AuthHandler, handlers.UpdateNote) 			//update notes by id (owned by user)
	noteGroup.Delete("/:id", middleware.AuthHandler, handlers.DeleteNote)		//delete notes by id (owned by user)

}
