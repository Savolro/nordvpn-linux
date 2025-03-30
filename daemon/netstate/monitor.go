package monitor

// Reconnector interface to reconnect on network state changes
type Reconnector interface {
	Reconnect(stateIsUp bool)
}
