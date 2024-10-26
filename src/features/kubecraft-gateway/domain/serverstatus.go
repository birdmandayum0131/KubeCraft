package domain

type ServerStatus string

const (
	Online   ServerStatus = "online"
	Offline  ServerStatus = "offline"
	Starting ServerStatus = "starting"
	Stopping ServerStatus = "stopping"
	Unknown  ServerStatus = "unknown"
)
