package main

import (
	"fmt"
	"log"
	hdl "org_chart/handlers"
	"org_chart/hub"
	"org_chart/initializers"
	mdl "org_chart/middleware"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	initializers.InitViper("./config.yaml")
	go initializers.ConnectPsqlDB(viper.GetString("db.psql_dsn"))
	go initializers.InitScyllaConnection(viper.GetString("db.scylla_dsn"))
	go hub.InitHub()
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//USER
	router.POST("/api/user/register", hdl.RegisterUser)
	router.POST("/api/user/login", hdl.LoginUser)
	router.POST("/api/token/refresh", hdl.RefreshToken)
	router.POST("/api/token/isvalid", hdl.VerifyToken)
	router.GET("/api/users", mdl.ProtectedRoute(), hdl.FetchAllUsers)

	//ROLE
	roles := router.Group("/api/roles")
	roles.Use(mdl.ProtectedRoute())

	roles.POST("", hdl.PostRoles) // CREATE
	roles.GET("", hdl.GetRoles)
	roles.GET("/:id", hdl.GetRoleById)   // READ
	roles.PATCH("/:id", hdl.PatchRole)   // UPDATE
	roles.DELETE("/:id", hdl.DeleteRole) // DELETE

	//EMPLOYEE
	employees := router.Group("/api/employees")
	employees.Use(mdl.ProtectedRoute())

	employees.POST("", hdl.PostEmployee) // CREATE
	employees.GET("", hdl.GetEmployees)
	employees.GET("/summary", hdl.GetEmployeeSummary) // READ
	employees.PATCH("/:id", hdl.PatchEmployee)        // UPDATE
	employees.DELETE("/:id", hdl.DeleteEmployee)      // DELETE

	//WEBSOCKET
	router.GET("/ws", mdl.ProtectedRoute(), hdl.WebSocketHandler)
	router.GET("/api/messages", mdl.ProtectedRoute(), hdl.GetConversation)

	if err := router.Run(viper.GetString("gin.PORT")); err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("API Active")
	}
}
