package handlers

import (
	"fmt"
	schemas "minecraftapi/interfaces/schemas/api"
	"minecraftapi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServerAPIHandler struct {
	ServerInteractor services.ServerInteractor
}

func (handler *ServerAPIHandler) GetServerStatusHandler(c *gin.Context) {
	status, err := handler.ServerInteractor.GetServerStatus()

	if err != nil {
		msg := fmt.Sprintf("Internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	response := schemas.ServerStatusResponse{
		ServerStatus: string(status),
	}

	c.IndentedJSON(http.StatusOK, response)
}


func (handler *ServerAPIHandler) StartServerHandler(c *gin.Context) {
	err := handler.ServerInteractor.StartServer()

	if err != nil {
		msg := fmt.Sprintf("Internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Server started"})
}

func (handler *ServerAPIHandler) StopServerHandler(c *gin.Context) {
	err := handler.ServerInteractor.StopServer()

	if err != nil {
		msg := fmt.Sprintf("Internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Server stopped"})
}
