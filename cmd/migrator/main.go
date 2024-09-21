package main

import (
	"fmt"
	"os"

	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/database"
	"github.com/VinukaThejana/todoapp/internal/enums"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var e = &env.Env{}
var db *gorm.DB

func init() {
	e.Load()
	db = database.Init(e)

	if e.Environ == string(enums.Dev) {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
		})
	}
}

func info(msg string) {
	log.Info().Msg(msg)
}

func main() {
	for _, table := range database.Tables {
		if !db.Migrator().HasTable(table.Schema) {
			db.AutoMigrate(table.Schema)
			info(fmt.Sprintf("creating the %s table", table.Name))
		} else {
			info(fmt.Sprintf("%s table already exists, skipping ... ", table.Name))
		}
	}
}
