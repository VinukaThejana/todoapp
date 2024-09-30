package database

import (
	"fmt"

	"github.com/VinukaThejana/go-utils/logger"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Init initializes the database connection
func Init(e *env.Env) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(e.DatabaseURL), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		logger.Errorf(fmt.Errorf("Failed to connect to database: %v", err))
	}

	return db
}
