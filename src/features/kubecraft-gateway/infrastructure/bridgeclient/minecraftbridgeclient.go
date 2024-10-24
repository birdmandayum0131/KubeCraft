package bridgeclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MineKubeMonitor implement services.ServerMonitor
type MinecraftBridgeClient struct {
	minecraft_bridge_url string
	ping_endpoint        string
	status_endpoint      string
}

func NewMinecraftBridgeClient(minecraft_bridge_url string) *MinecraftBridgeClient {
	if !strings.HasSuffix(minecraft_bridge_url, "/") {
		minecraft_bridge_url += "/"
	}

	return &MinecraftBridgeClient{
		minecraft_bridge_url: minecraft_bridge_url,
		ping_endpoint:        minecraft_bridge_url + "api/v1/minecraft/status",
		status_endpoint:      minecraft_bridge_url + "api/v1/minecraft/ping",
	}
}

// Use server replicas and mcstatus to determine server status
func (c *MinecraftBridgeClient) Ping() (int, error) {
	resp, err := http.Get(c.ping_endpoint)
	if err != nil {
		return -1, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		var pingResp pingResponse
		err = json.Unmarshal(body, &pingResp)
		if err != nil {
			return -1, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return pingResp.Latency, nil
	} else {
		var errResp errorResponse
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return -1, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return -1, &ConnectionError{s: errResp.Error}
	}
}
