package controllers

import (
	"todolist/api/middlewares"
)

func (s *Server) initializeRoutes() {
	v1 := s.Router.Group("/api/v1")
	{
		v1.POST("/login", s.Login)
		v1.POST("/signup", s.CreateUser)

		v1.POST("/todos", middlewares.TokenAuthMiddleware(), s.CreateToDo)
		v1.PUT("/todos/:id", middlewares.TokenAuthMiddleware(), s.UpdateToDo)
		v1.DELETE("/todos/:id", middlewares.TokenAuthMiddleware(), s.DeleteToDo)
		v1.GET("/user_todos/:id", middlewares.TokenAuthMiddleware(), s.GetUserToDos)
	}
}
