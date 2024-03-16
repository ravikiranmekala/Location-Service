package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ravikiranmekala/MapUp-Backend-assessment/database"
	"github.com/ravikiranmekala/MapUp-Backend-assessment/routes"
)

func main() {
	router := gin.New()
	database.ConnectDb()
	routes.LocationRoute(router)
	router.Run(":3000")
}
