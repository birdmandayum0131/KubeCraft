package bridgeclient

type pingResponse struct {
	Latency float64 `json:"latency"`
}

type errorResponse struct {
	Error string `json:"error"`
}
