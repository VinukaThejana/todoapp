// Package database provides the models for the database tables
package database

import (
	"time"

	"gorm.io/gorm"
)

// Tables contains the list of tables to be created
var Tables = []struct {
	Name   string
	Schema any
}{
	{
		Name:   "users",
		Schema: User{},
	},
	{
		Name:   "todos",
		Schema: Todo{},
	},
	{
		Name:   "sessions",
		Schema: Session{},
	},
}

// User is a model for the user table
type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"type:varchar(255);unique_index"`
	Username string `gorm:"type:varchar(100);unique_index"`
	Password string `gorm:"not null"`
}

// Todo is a model for the todo table
type Todo struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Content     string `gorm:"not null"`
	Completed   bool   `gorm:"not null"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Session is a model for the session table
type Session struct {
	ID        string `gorm:"primarykey"`
	ExpiresAt int64  `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	LoginAt   time.Time
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
