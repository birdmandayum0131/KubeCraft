package main

import (
	"flag"
	"fmt"
	"kubecraft-gateway/infrastructure"
	"kubecraft-gateway/infrastructure/bridgeclient"
	"kubecraft-gateway/interfaces/rest/handlers"
	"kubecraft-gateway/interfaces/rest/routes"
	"kubecraft-gateway/services"
	"os"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeConfigPath = flag.String("kubeConfig", "./.kube/config", "kube config files path")

func main() {
	flag.Parse()

	mineKubeConfig, err := createMineKubeConfig()
	if err != nil {
		panic(err)
	}

	minecraftKubeController, err := createKubeController(mineKubeConfig)
	if err != nil {
		panic(err)
	}

	bridgeClient, err := createBridgeClient()
	if err != nil {
		panic(err)
	}

	minecraftKubeMonitor := createMineKubeMonitor(*mineKubeConfig, *bridgeClient, minecraftKubeController)

	restAPIHandler := createRestHandler(minecraftKubeMonitor, minecraftKubeController)

	app := setupRouter(nil, restAPIHandler)
	err = app.Run(":8000")
	if err != nil {
		panic(err)
	}
}

func createMineKubeConfig() (*infrastructure.MinecraftKubeConfig, error) {
	// create k8s client config from environment variables
	envSrvNs, envSrvNsExist := os.LookupEnv("MINECRAFT_SERVER_NAMESPACE")
	envSrvDeploy, envSrvDeployExist := os.LookupEnv("MINECRAFT_SERVER_DEPLOYMENT")
	if !envSrvNsExist || !envSrvDeployExist {
		return nil, fmt.Errorf("MINECRAFT_SERVER_NAMESPACE and MINECRAFT_SERVER_DEPLOYMENT are not set")
	} else {
		return &infrastructure.MinecraftKubeConfig{
			Namespace:      envSrvNs,
			DeploymentName: envSrvDeploy,
		}, nil
	}
}

func createKubeController(config *infrastructure.MinecraftKubeConfig) (*infrastructure.MinecraftKubeController, error) {
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
		Config:    *config,
		Clientset: clientset,
	}, nil
}

func createBridgeClient() (*bridgeclient.MinecraftBridgeClient, error) {
	env_bridge_url, env_bridge_url_exist := os.LookupEnv("MINECRAFT_BRIDGE_URL")
	if !env_bridge_url_exist {
		return nil, fmt.Errorf("MINECRAFT_SERVER_NAMESPACE and MINECRAFT_SERVER_DEPLOYMENT are not set")
	} else {
		return bridgeclient.NewMinecraftBridgeClient(env_bridge_url), nil
	}
}

func createMineKubeMonitor(config infrastructure.MinecraftKubeConfig, bridgeClient bridgeclient.MinecraftBridgeClient, watcher infrastructure.DeploymentWatcher) *infrastructure.MineKubeMonitor {
	return &infrastructure.MineKubeMonitor{
		Config:        config,
		DeployWatcher: watcher,
		BridgeClient:  bridgeClient,
	}
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
