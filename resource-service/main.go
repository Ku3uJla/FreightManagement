package main

import (
	"user-service/internal/repository/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.InitDB()
	router := gin.Default()
}
