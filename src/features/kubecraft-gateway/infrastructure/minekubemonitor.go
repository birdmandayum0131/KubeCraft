package infrastructure

import (
	"fmt"
	"kubecraft-gateway/domain"
	"kubecraft-gateway/infrastructure/bridgeclient"
)

type DeploymentWatcher interface {
	GetServerReplicas(deployment string, namespace string) (int32, int32, error)
}

// MineKubeMonitor implement services.ServerMonitor
type MineKubeMonitor struct {
	Config        MinecraftKubeConfig
	DeployWatcher DeploymentWatcher
	BridgeClient  bridgeclient.MinecraftBridgeClient
}

// Use server replicas and mcstatus to determine server status
func (c *MineKubeMonitor) GetServerStatus() (domain.ServerStatus, error) {
	targetReplicas, currentReplicas, err := c.DeployWatcher.GetServerReplicas(c.Config.DeploymentName, c.Config.Namespace)
	if err != nil {
		return domain.Unknown, fmt.Errorf("failed to get minecraft server deployment replicas: %w", err)
	}

	// * Ping minecraft server to check if server is ready
	if targetReplicas > 0 {
		_, err = c.BridgeClient.Ping()
		if err != nil {
			_, ok := err.(*bridgeclient.ConnectionError)
			if ok {
				// * Server is starting but not ready
				return domain.Starting, nil
			} else {
				return domain.Unknown, fmt.Errorf("failed to ping minecraft server: %w", err)
			}
		} else {
			// * Server is ready
			return domain.Online, nil
		}
	} else {
		if currentReplicas == 0 {
			// * Server is offline
			return domain.Offline, nil
		} else {
			// * Server is stopping
			return domain.Stopping, nil
		}
	}
}
