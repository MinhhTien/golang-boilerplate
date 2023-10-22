package middlewares

import (
	"net/http"
	"todolist/api/auth"

	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware is a middleware function that validates the token sent by the client
// and returns an error if the token is not valid
func TokenAuthMiddleware() gin.HandlerFunc {
	errList := make(map[string]string)
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			errList["unauthorized"] = "Unauthorized"
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  errList,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CORSMiddleware is a middleware function that enables us to interact with the React Frontend
// by setting the necessary headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		// If the request method is OPTIONS, abort the request with status code 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
