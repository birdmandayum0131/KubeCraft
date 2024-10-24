package bridgeclient

type ConnectionError struct {
	s string
}

func (e *ConnectionError) Error() string {
	return e.s
}
