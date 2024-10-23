package main

import (
	"flag"
	"fmt"
	"minecraftapi/infrastructure"
	"minecraftapi/interfaces/rest/handlers"
	"minecraftapi/interfaces/rest/routes"
	"minecraftapi/services"
	"os"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeConfigPath = flag.String("kubeConfig", "./.kube/config", "kube config files path")

func main() {
	flag.Parse()

	minecraftKubeController, err := createKubeController()
	if err != nil {
		panic(err)
	}

	restAPIHandler := createRestHandler(minecraftKubeController, minecraftKubeController)

	app := setupRouter(nil, restAPIHandler)
	err = app.Run(":8000")
	if err != nil {
		panic(err)
	}
}

func createKubeController() (*infrastructure.MinecraftKubeController, error) {
	// create k8s client config from environment variables
	envSrvNs, envSrvNsExist := os.LookupEnv("MINECRAFT_SERVER_NAMESPACE")
	envSrvDeploy, envSrvDeployExist := os.LookupEnv("MINECRAFT_SERVER_DEPLOYMENT")
	var kubeCtrlConfig infrastructure.MinecraftKubeCtrlConfig
	if !envSrvNsExist && !envSrvDeployExist {
		return nil, fmt.Errorf("MINECRAFT_SERVER_NAMESPACE and MINECRAFT_SERVER_DEPLOYMENT are not set")
	} else {
		kubeCtrlConfig = infrastructure.MinecraftKubeCtrlConfig{
			Namespace:      envSrvNs,
			DeploymentName: envSrvDeploy,
		}
	}

	// read cluster config from service account token
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		// load kube config if in development environment
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load kube config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	return &infrastructure.MinecraftKubeController{
		Config:    kubeCtrlConfig,
		Clientset: clientset,
	}, nil
}

func createRestHandler(serverMonitor services.ServerMonitor, serverManager services.ServerManager) *handlers.ServerAPIHandler {
	return &handlers.ServerAPIHandler{
		ServerInteractor: services.ServerInteractor{
			ServerMonitor: serverMonitor,
			ServerManager: serverManager,
		},
	}
}

func setupRouter(middlewares []gin.HandlerFunc, handler *handlers.ServerAPIHandler) *gin.Engine {
	router := gin.Default()
	routes.InitMinecraftRoutes(router, middlewares, handler)
	return router
}