package infrastructure

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// MinecraftKubeController implement services.ServerManager and services.ServerMonitor
// through kubernetes api
type MinecraftKubeController struct {
	Config    MinecraftKubeConfig
	Clientset *kubernetes.Clientset
}

func (c *MinecraftKubeController) GetServerReplicas(deploymentName string, namespace string) (int32, int32, error) {
	deployClient := c.Clientset.AppsV1().Deployments(namespace)
	deployment, err := deployClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get minecraft server deployment: %w", err)
	}

	return *deployment.Spec.Replicas, deployment.Status.Replicas, nil
}

func (c *MinecraftKubeController) StartServer() error {
	deployClient := c.Clientset.AppsV1().Deployments(c.Config.Namespace)
	deployment, err := deployClient.Get(context.TODO(), c.Config.DeploymentName, metav1.GetOptions{})

	if err != nil {
		return fmt.Errorf("failed to get minecraft server deployment: %w", err)
	}

	// Start server by setting its replicas to 1
	deployment.Spec.Replicas = new(int32)
	*deployment.Spec.Replicas = 1
	_, err = deployClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update minecraft server deployment: %w", err)
	}
	return nil
}

func (c *MinecraftKubeController) StopServer() error {
	deployClient := c.Clientset.AppsV1().Deployments(c.Config.Namespace)
	deployment, err := deployClient.Get(context.TODO(), c.Config.DeploymentName, metav1.GetOptions{})

	if err != nil {
		return fmt.Errorf("failed to get minecraft server deployment: %w", err)
	}

	// Stop server by setting its replicas to 0
	deployment.Spec.Replicas = new(int32)
	*deployment.Spec.Replicas = 0
	_, err = deployClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update minecraft server deployment: %w", err)
	}
	return nil
}
