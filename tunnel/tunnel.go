package tunnel

import (
	"net"
	"net/netip"
)

// T describes tunnel behavior
// probably needs a better name, though
type T interface {
	Interface() net.Interface
	IPs() []netip.Addr
	TransferRates() (Statistics, error)
}

// Statistics defines what information can be collected about the tunnel
type Statistics struct {
	Tx uint64
	Rx uint64
}
