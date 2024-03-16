package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ravikiranmekala/MapUp-Backend-assessment/controller"
)

func LocationRoute(router *gin.Engine) {
	router.GET("/locations/:category", controller.GetLocationController)
	router.POST("/locations", controller.CreateLocationController)
	router.POST("/search", controller.SearchLocationController)
	router.POST("/trip-cost/:location_id", controller.TripCostController)
}
