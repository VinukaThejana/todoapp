// Package enums provides a way to define constants for the application.
package enums

// Env represents the environment of the application
type Env string

const (
	// Dev represents the local development environment
	Dev Env = "dev"
	// Stg represents the staging environment
	Stg Env = "stg"
	// Prd represents the production environment
	Prd Env = "prd"
)
