package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kevbradwick/tflapi/app"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(app.PanicHandler())
	station := router.Group("/stations")
	{
		station.GET("/search", app.SearchHandler)
		station.GET("/station/:id", app.GetStationHandler)
	}
	router.Run("0.0.0.0:8000")
}
