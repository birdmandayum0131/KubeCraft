package routes

import (
	handlers "kubecraft-gateway/interfaces/rest/handlers"

	"github.com/gin-gonic/gin"
)

// Initialize route for minecraft api server
func InitMinecraftRoutes(route *gin.Engine, middlewares []gin.HandlerFunc, handler *handlers.ServerAPIHandler) {
	groupRoute := route.Group("/api/v1")
	groupRoute.Use(middlewares...)
	groupRoute.GET("/server/minecraft/status", handler.GetServerStatusHandler)
	groupRoute.POST("/server/minecraft/start", handler.StartServerHandler)
	groupRoute.POST("/server/minecraft/stop", handler.StopServerHandler)
}
