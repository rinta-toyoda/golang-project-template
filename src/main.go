package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	err := server.Run(":8081")
	if err != nil {
		return
	} // localhost:8080
}
