package dns

// Setter is responsible for configuring DNS.
type Setter interface {
	Set(iface string, nameservers []string) error
	Unset(iface string) error
}

// Method is abstraction of DNS handling method
type Method interface {
	Set(iface string, nameservers []string) error
	Unset(iface string) error
	Name() string
}
