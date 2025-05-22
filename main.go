package main

import "fmt"
import "github.com/gin-gonic/gin"

func main() {
	fmt.Println("Corriendo en localhost:8080 !!!")
	
  	router := gin.Default()
  	router.GET("/", func(c *gin.Context) {
  	  	c.JSON(200, gin.H{
  	    	"message": "pong",
  	  	})
  	})
  	router.Run() // listen and serve on 0.0.0.0:8080
}