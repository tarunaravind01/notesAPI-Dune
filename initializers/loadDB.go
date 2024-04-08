package initializers

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// postgres handler(users table)
var DB *gorm.DB

// mongo handler(notes table)
var NoteDB *mongo.Database

// postgres
func ConnectToUserDB() {
	var err error
	dsn := os.Getenv("POSTGRES_CONN")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to auth db")
	}
}

// mongo
func InitNotesDB() {
	uri := os.Getenv("MONGO_CONN")

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	NoteDB = client.Database("notes")

}

func CloseDB() error {
	return NoteDB.Client().Disconnect(context.Background())
}
