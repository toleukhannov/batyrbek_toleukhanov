package main

import (
	"https://github.com/toleukhannov/batyrbek_toleukhanov/controllers"
	"https://github.com/toleukhannov/batyrbek_toleukhanov/database"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

}
