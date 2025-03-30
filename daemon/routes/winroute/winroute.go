package winroute

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/NordSecurity/nordvpn-linux/daemon/routes"
)

// Router uses netlink under the hood.
type Router struct {
	routes []routes.Route
	sync.Mutex
}

// Add appends route to a routes list via netlink if it does not exist yet.
func (r *Router) Add(route routes.Route) error {
	r.Lock()
	defer r.Unlock()
	if r.has(route) {
		return fmt.Errorf("route %+v already exists", route)
	}

	// check if such route does not exist in routing table
	if err := applyRoute("add", route); err != nil {
		if errors.Is(err, os.ErrExist) {
			return routes.ErrRouteToOtherDestinationExists
		}
		return fmt.Errorf("adding route %+v to a routing table: %w", route, err)
	}

	r.routes = append(r.routes, route)
	return nil
}

// Flush deletes all previously added routes via netlink.
func (r *Router) Flush() error {
	r.Lock()
	defer r.Unlock()
	var errs []error
	for _, route := range r.routes {
		if err := applyRoute("delete", route); err != nil {
			errs = append(errs, fmt.Errorf("deleting route %+v: %w", route, err))
			continue
		}
	}
	r.routes = nil
	return errors.Join(errs...)
}

// has returns true if router contains a given route in its memory.
func (r *Router) has(route routes.Route) bool {
	return slices.ContainsFunc(r.routes, route.IsEqual)
}

func applyRoute(addParam string, route routes.Route) error {
	ipNet := prefixToIPNet(route.Subnet)
	routeList := []routes.Route{}
	gateway := route.Gateway
	routeToIface := false
	if !gateway.IsValid() && route.Device.Name != "" {
		addrs, err := route.Device.Addrs()
		if err != nil {
			return fmt.Errorf("retrieving device")
		}
		if len(addrs) != 0 {
			var addr netip.Addr
			for _, iAddr := range addrs {
				addr, err = netip.ParseAddr(strings.Split(iAddr.String(), "/")[0])
				if err != nil {
					continue
				}
				if addr.Is4() {
					break
				}
			}
			route.Gateway = addr
			routeToIface = true
		}
	}
	routeList = append(routeList, route)
	for _, route := range routeList {
		// TODO: Remove this and investigate why routes are not there if added immediately after successful connection.
		time.Sleep(time.Second * 3)
		args := []string{
			addParam,
			ipNet.IP.String(),
			"mask", net.IP(ipNet.Mask).String(),
			route.Gateway.String(),
		}
		if addParam == "add" {
			args = append(args, "metric", "1")
		}
		if routeToIface && addParam == "delete" {
			args = []string{
				addParam,
				ipNet.IP.String(),
				"mask", net.IP(ipNet.Mask).String(),
				"if", strconv.Itoa(route.Device.Index),
			}
		}
		if out, err := exec.Command("route", args...).CombinedOutput(); err != nil {
			return fmt.Errorf("executing ip route command: %w: %s", err, string(out))
		}
	}
	return nil
}

func prefixToIPNet(prefix netip.Prefix) *net.IPNet {
	addr := prefix.Addr()
	bits := net.IPv4len * 8
	if addr.Is6() {
		bits = net.IPv6len * 8
	}
	return &net.IPNet{
		IP:   addr.AsSlice(),
		Mask: net.CIDRMask(prefix.Bits(), bits),
	}
}
