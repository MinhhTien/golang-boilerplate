package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"todolist/api/auth"
	"todolist/api/models"
	"todolist/api/utils/formaterror"

	"github.com/gin-gonic/gin"
)

// CreateToDo is a function that creates a new todo item
// It takes in a gin.Context and returns a JSON response
// It requires a valid token to be passed in the request header
// It requires a valid JSON body with the todo item details
func (server *Server) CreateToDo(c *gin.Context) {
	//clear previous error if any
	errList = map[string]string{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Unmarshal the request body into a todo item
	todo := models.ToDo{}
	err = json.Unmarshal(body, &todo)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Extract the token from the request header
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// check if the user exist:
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Set the author ID to the authenticated user
	todo.AuthorID = uid //the authenticated user is the one creating the todo
	// Prepare the todo item
	todo.Prepare()
	// Validate the todo item
	errorMessages := todo.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Save the todo item
	todoCreated, err := todo.SaveToDo(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	// Return the created todo item
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": todoCreated,
	})
}

// UpdateToDo is a function that updates a ToDo item in the database.
// It takes in a gin.Context object as an argument and returns a JSON response.
func (server *Server) UpdateToDo(c *gin.Context) {
	//clear previous error if any
	errList = map[string]string{}
	// Get the ToDo ID from the URL parameter
	todoID := c.Param("id")
	// Check if the todo id is valid
	pid, err := strconv.ParseUint(todoID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	// Extract the user ID from the auth token
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	//Check if the todo exist
	origToDo := models.ToDo{}
	err = server.DB.Debug().Model(models.ToDo{}).Where("id = ?", pid).Take(&origToDo).Error
	if err != nil {
		errList["No_todo"] = "No ToDo Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Check if the authenticated user is the owner of the ToDo
	if uid != origToDo.AuthorID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Read the data from the request body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Unmarshal the data into a ToDo object
	todo := models.ToDo{}
	err = json.Unmarshal(body, &todo)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Set the ToDo ID and Author ID
	todo.ID = origToDo.ID //this is important to tell the model the todo id to update, the other update field are set above
	todo.AuthorID = origToDo.AuthorID
	// Prepare the ToDo object
	todo.Prepare()
	// Validate the ToDo object
	errorMessages := todo.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Update the ToDo in the database
	todoUpdated, err := todo.UpdateAToDo(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	// Return the updated ToDo
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": todoUpdated,
	})
}

// DeleteToDo is a function that deletes a ToDo item from the database.
// It takes in a gin.Context object as an argument and returns a JSON response.
func (server *Server) DeleteToDo(c *gin.Context) {
	todoID := c.Param("id")
	// Check if the ToDo ID is valid
	pid, err := strconv.ParseUint(todoID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	fmt.Println("this is delete todo ")
	// Extract the user ID from the auth token
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Check if the todo exist
	todo := models.ToDo{}
	err = server.DB.Debug().Model(models.ToDo{}).Where("id = ?", pid).Take(&todo).Error
	if err != nil {
		errList["No_todo"] = "No ToDo Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Check if the authenticated user is the owner of the ToDo
	if uid != todo.AuthorID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// If all the conditions are met, delete the ToDo
	_, err = todo.DeleteAToDo(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	// Return a success message
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "todo deleted",
	})
}

// GetUserToDos is a function that takes in a gin.Context and returns a JSON response containing
// the ToDos associated with the userID provided in the request.
// Parameters: c *gin.Context: The gin.Context object containing the request information
// Returns: JSON response containing the ToDos associated with the userID provided in the request
func (server *Server) GetUserToDos(c *gin.Context) {
	//clear previous error if any
	errList = map[string]string{}
	userID := c.Param("id")
	// Check if the userID is valid
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	// Create a ToDo object
	todo := models.ToDo{}
	// Get the ToDos associated with the userID
	todos, err := todo.FindUserToDos(server.DB, uint32(uid))
	if err != nil {
		errList["No_todo"] = "No ToDo Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Return the ToDos in a JSON response
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": todos,
	})
}
