package handlers

import (
	"fmt"
	"notes-api/initializers"
	"notes-api/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// validates and organizes request body
type CreateNoteDTO struct {
	Title	string 		`json:"title" bson:"title"`
	Content string 		`json:"content" bson:"content"`
}

type UpdateNoteDTO struct {
	Title	string 		`json:"title,omitempty" bson:"title,omitempty"`
	Content string 		`json:"content,omitempty" bson:"content,omitempty"`
}

type FullyFormedDTO struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	UserID   string `json:"userId" bson:"userId"`
	Author   string `json:"author" bson:"author"`
	Title    string `json:"title" bson:"title"`
	Content  string `json:"content" bson:"content"`
}


// create a new note
// POST
func CreateNoteHandler(c *fiber.Ctx) error {
	// collect req user deets
	user, ok := c.Locals("user").(models.User)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	author, userId := user.Username, fmt.Sprintf("%d", user.ID)

	// collect and validate req body
	dto := new(CreateNoteDTO)
	err := c.BodyParser(&dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Body"})
	}

	// build mongo query
	note := FullyFormedDTO{
		UserID:  userId,
		Author:  author,
		Title:   dto.Title,
		Content: dto.Content,
	}

	collection := initializers.NoteDB.Collection("dune_showcase")

	result, err := collection.InsertOne(c.Context(), note)
	// respond back
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":"Failed to create Note",
			"message":err.Error(),
		})
	}

	
	return c.Status(201).JSON(fiber.Map{
		"result":result,
	})

}


// get all notes for a user
// GET
func GetAllNotes(c *fiber.Ctx) error {
	// collect req user deets
	user, ok := c.Locals("user").(models.User)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// slice to hold notes
	notes := make([]models.Note, 0)

	userId := fmt.Sprintf("%d", user.ID)

	// build mongo query
	collection := initializers.NoteDB.Collection("dune_showcase")

	filter := bson.D{{Key: "userId", Value: userId}}

	cursor, err := collection.Find(c.Context(), filter)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	// save notes to slice
	for cursor.Next(c.Context()) {
		note := models.Note{}
		err := cursor.Decode(&note)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":err.Error(),
			})
		}
		notes = append(notes, note)
	}

	// respond back
	return c.Status(200).JSON(fiber.Map{
		"result":notes,
	})
}


// search a user's note by "title"
// GET
func FindNote(c *fiber.Ctx) error {
	// collect req user deets

	title := c.Params("title")

	user, ok := c.Locals("user").(models.User)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userId := fmt.Sprintf("%d", user.ID)

	// build query to mongo
	collection := initializers.NoteDB.Collection("dune_showcase")

	filter := bson.D{{Key: "userId", Value: userId}, {Key: "title", Value: title} }

	note := models.Note{}

	err := collection.FindOne(c.Context(), filter).Decode(&note)

	// respond back
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"result":note,
	})

}

// update a note based on _id and userId
// PUT
func UpdateNote(c *fiber.Ctx) error {
	// collect deets from request
	user, ok := c.Locals("user").(models.User)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userId := fmt.Sprintf("%d", user.ID)

	// validate and index req body
	dto := new(UpdateNoteDTO)
	err := c.BodyParser(&dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Body"})
	}
	id := c.Params("id")
	if id == ""{
		return c.Status(400).JSON(fiber.Map{
			"error":"id is required",
		})
	}
	noteId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":"invalid id",
		})
	}

	// build update query
	collection := initializers.NoteDB.Collection("dune_showcase")

	filter := bson.D{{Key: "_id", Value: noteId}, {Key: "userId", Value: userId}}

	result, err := collection.UpdateOne(c.Context(), filter, bson.M{"$set": dto})

	// respond back to user
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":"could not update note",
			"message":err.Error(),
		})
	}

	if result.ModifiedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"Message":"Note id not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"result":result,
	})
	
}


// delete a note based on _id and userID
// DELETE
func DeleteNote(c *fiber.Ctx) error {
	// collect deets from req
	user, ok := c.Locals("user").(models.User)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userId := fmt.Sprintf("%d", user.ID)

	id := c.Params("id")
	if id == ""{
		return c.Status(400).JSON(fiber.Map{
			"error":"id is required",
		})
	}

	noteId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":"invalid id",
		})
	}

	// build delete query
	collection := initializers.NoteDB.Collection("dune_showcase")

	filter := bson.D{{Key: "_id", Value: noteId}, {Key: "userId", Value: userId}}

	result, err := collection.DeleteOne(c.Context(), filter)

	// respond back
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":"could not delete book",
			"message":err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"Message":"Note id not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"result":result,
	})
}