package api

import (
	"fmt"
	"log"
	"os"
	"todolist/api/controllers"

	"github.com/joho/godotenv"
)

var server = controllers.Server{}

// init is a function that is called automatically when the program starts
// it loads the environment variables from the .env file
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

// Run is the main function of the api package
// It is responsible for initializing the server and running it
func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting values")
	}
	// Initialize the server with the environment variables
	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	// This is for testing, when done, do well to comment
	// seed.Load(server.DB)
	// Get the API port from the environment variables
	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", apiPort)
	// Run the server
	server.Run(apiPort)
}
