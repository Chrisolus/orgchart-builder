package initializers

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPsqlDB(psql_dsn string) {

	var err error
	DB, err = gorm.Open(postgres.Open(psql_dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB Connection Failed: ", err.Error())
	} else {
		log.Println("Postgres Connected")
	}
}
