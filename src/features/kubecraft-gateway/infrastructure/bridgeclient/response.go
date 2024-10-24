package bridgeclient

type pingResponse struct {
	Latency int `json:"latency"`
}

type errorResponse struct {
	Error string `json:"error"`
}
