package main


import (
	"net/http"
	dbhandler "vumble/back-end/db"
	"github.com/gin-gonic/gin"
)

db := dbhandler.newVumbleDB()

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// login route
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		// Check if the username exists
		// If it doesn't, return an error
		
		if(db.getUser(username) == nil) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		

		res, err := db.login(username, password)
		if(err != nil) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incorrect password"})
			return
		}

		// Set the cookie to the OAuth token
		c.SetCookie("token", res, 3600, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	
	

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in
	println("Server started")
	r.Run(":8080")
}
