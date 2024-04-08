package initializers

import "notes-api/models"

// sync with Supabase
func MigrateDB() {
	DB.AutoMigrate(&models.User{}) 
}
