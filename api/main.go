package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.Use(LoggerMiddleware())
	// r.Use(RateLimitMiddleware())

	r.POST("/share-file", shareFileUpload)
	r.GET("/shares", getFileShare)
	r.GET("/download", downloadFileResumable)

	r.Run(":8080")

}
