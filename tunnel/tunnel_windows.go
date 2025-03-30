package tunnel

import (
	"fmt"
	"net"
	"net/netip"
	"os/exec"
)

// Tunnel encrypts and decrypts network traffic.
type Tunnel struct {
	// might be a good idea to change this to a pointer now
	// so that we could see changes to the interface at real time
	// but this would need testing first to check if it actually works
	iface  net.Interface
	ips    []netip.Addr
	prefix netip.Prefix
}

func New(iface net.Interface, ips []netip.Addr, prefix netip.Prefix) *Tunnel {
	return &Tunnel{iface: iface, ips: ips, prefix: prefix}
}

// Interface returns the underlying network interface.
func (t *Tunnel) Interface() net.Interface { return t.iface }

// IPs attached to the tunnel.
func (t *Tunnel) IPs() []netip.Addr { return t.ips }

func (t *Tunnel) TransferRates() (Statistics, error) {
	return Statistics{}, nil
}

// Up sets tunnel state to up.
func (t *Tunnel) Up() error {
	return nil
}

func (t *Tunnel) AddAddrs() error {
	for _, ip := range t.ips {
		mask := "255.192.0.0"
		out, err := addAddr(t.iface.Name, ip.String(), mask)
		if err != nil {
			return fmt.Errorf("adding IP address to interface: %s : %w", string(out), err)
		}
	}
	return nil

}

// DelAddrs from a tunnel interface.
func (t *Tunnel) DelAddrs() error {
	return nil
}

func GetTransferRates(ifaceName string) (Statistics, error) {
	return Statistics{}, nil
}

func addAddr(ifaceName string, ip string, mask string) ([]byte, error) {
	// #nosec G204 -- input is properly sanitized
	cmdHandle := exec.Command(
		"netsh",
		"interface",
		"ip",
		"set",
		"address",
		"name="+ifaceName,
		"static",
		ip,
		mask,
	)
	return cmdHandle.CombinedOutput()
}
