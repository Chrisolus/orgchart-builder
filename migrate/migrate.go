package main

import (
	"org_chart/initializers"
	"org_chart/models"

	"github.com/spf13/viper"
)

func init() {
	initializers.InitViper("./")
	initializers.ConnectPsqlDB(viper.GetString("db.psql_dsn"))
}

func main() {
	initializers.DB.AutoMigrate(&models.Employee{})
	initializers.DB.AutoMigrate(&models.Role{})
	initializers.DB.AutoMigrate(&models.User{})
}
