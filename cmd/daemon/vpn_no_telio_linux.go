//go:build !telio

package main

import "github.com/NordSecurity/nordvpn-linux/daemon/vpn"

func getNordlynxVPN(envIsDev bool,
	eventsDbPath string,
	fwmark uint32,
	cfg vpn.LibConfigGetter,
	appVersion string,
	eventsPublisher *vpn.Events) (*nordlynx.KernelSpace, error) {
	return nordlynx.NewKernelSpace(fwmark, eventsPublisher), nil
}
