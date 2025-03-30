package device

import "net"

type ListFunc func() ([]net.Interface, error)

func InterfacesAreEqual(a net.Interface, b net.Interface) bool {
	return a.Index == b.Index &&
		a.MTU == b.MTU &&
		a.Name == b.Name &&
		a.HardwareAddr.String() == b.HardwareAddr.String() &&
		a.Flags == b.Flags
}
