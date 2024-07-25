package server

import (
	"applied_system/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	graphCtrl := controllers.NewGraphController()
	router := gin.Default()
	router.POST("/graphs", graphCtrl.PostGraphHandler)
	router.GET("/graphs/:id/shortest-path", graphCtrl.GetShortestPathHandler)
	router.DELETE("/graphs/:id", graphCtrl.DeleteGraphHandler)

	return router
}
