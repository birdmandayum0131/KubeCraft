package infrastructure

import (
	"fmt"
	"kubecraft-gateway/domain"
	"kubecraft-gateway/infrastructure/bridgeclient"
)

type KubeWatcher interface {
	GetServerReplicas(deployment string, namespace string) (*int32, error)
	GetServerPodsNumber(namespace string) (*int, error)
}

// MineKubeMonitor implement services.ServerMonitor
type MineKubeMonitor struct {
	Config       MinecraftKubeConfig
	KubeWatcher  KubeWatcher
	BridgeClient bridgeclient.MinecraftBridgeClient
}

// Use server replicas and mcstatus to determine server status
func (c *MineKubeMonitor) GetServerStatus() (domain.ServerStatus, error) {
	replicas, err := c.KubeWatcher.GetServerReplicas(c.Config.DeploymentName, c.Config.Namespace)
	if err != nil {
		return domain.Unknown, fmt.Errorf("failed to get minecraft server deployment replicas: %w", err)
	}

	pods, err := c.KubeWatcher.GetServerPodsNumber(c.Config.Namespace)
	if err != nil {
		return domain.Unknown, fmt.Errorf("failed to get minecraft server pods: %w", err)
	}

	if *replicas > 0 {
		// * Server is started or starting
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
		// * Server is stopped or stopping
		if *pods == 0 {
			// * Server is offline
			return domain.Offline, nil
		} else {
			// * Server is stopping
			return domain.Stopping, nil
		}
	}
}
