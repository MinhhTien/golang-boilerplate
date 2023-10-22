package controllers

import (
	"fmt"
	"log"
	"net/http"
	"todolist/api/middlewares"
	"todolist/api/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

// errList is a map to store errors
var errList = make(map[string]string)

// Initialize is a function that initializes the server
// It takes in the database driver, user, password, port, host and name
// It returns an error if there is one
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else {
		fmt.Println("Unknown Driver")
	}

	//database migration
	server.DB.Debug().AutoMigrate(
		&models.User{},
		&models.ToDo{},
	)

	// Init router
	server.Router = gin.Default()
	// User CORS middleware
	server.Router.Use(middlewares.CORSMiddleware())
	server.initializeRoutes()
}

// Run is a function that runs the server
// It takes in the address of the server
// It returns an error if there is one
func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
   }