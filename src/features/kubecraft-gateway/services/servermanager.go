package services

import (
	"fmt"
	"kubecraft-gateway/domain"
)

type ServerMonitor interface {
	GetServerStatus() (domain.ServerStatus, error)
}

type ServerManager interface {
	StartServer() error
	StopServer() error
}

type ServerInteractor struct {
	ServerMonitor ServerMonitor
	ServerManager ServerManager
}

// GetServerStatus attempts to get the status of the server and returns an error if any occurs.
//
// The status of the server is either "online" or "offline".
//
// This function does not block.
func (s *ServerInteractor) GetServerStatus() (domain.ServerStatus, error) {
	status, err := s.ServerMonitor.GetServerStatus()

	if err != nil {
		return "", fmt.Errorf("failed to get server status: %w", err)
	}

	return status, nil
}

// StartServer attempts to start the server and returns an error if any occurs.
//
// This function won't block. It will return an error if the server is already running.
func (s *ServerInteractor) StartServer() error {
	status, err := s.ServerMonitor.GetServerStatus()

	if err != nil {
		return fmt.Errorf("failed to start server, error occurs when fetching server status: %w", err)
	}

	if status == domain.Online {
		return fmt.Errorf("failed to start server, server is already running")
	}

	err = s.ServerManager.StartServer()

	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// StopServer attempts to stop the server and returns an error if any occurs.
//
// This function won't block. It will return an error if the server is already stopped.
func (s *ServerInteractor) StopServer() error {
	status, err := s.ServerMonitor.GetServerStatus()

	if err != nil {
		return fmt.Errorf("failed to stop server, error occurs when fetching server status: %w", err)
	}

	if status == domain.Offline {
		return fmt.Errorf("failed to stop server, server is already stopped")
	}

	err = s.ServerManager.StopServer()

	if err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}
