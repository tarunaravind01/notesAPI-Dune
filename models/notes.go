package models

// data model for notes
type Note struct {
	ID      string 		`json:"id" bson:"_id"`
	UserID  string 		`json:"userId" bson:"userId"`
	Author	string		`json:"author" bson:"author"`
	Title   string 		`json:"title" bson:"title"`
	Content string 		`json:"content" bson:"content"`
}


