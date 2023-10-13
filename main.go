package main

import (
	"fmt"
	//"net/http" --> Will need this for HTML frontend
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World")
	router := gin.Default()
	router.Run("localhost:8080")
}
