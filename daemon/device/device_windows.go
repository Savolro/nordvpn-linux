package device

import (
	mapset "github.com/deckarep/golang-set/v2"
)

func InterfacesWithDefaultRoute(ignoreSet mapset.Set[string]) mapset.Set[string] {
	return mapset.NewSet[string]()
}
