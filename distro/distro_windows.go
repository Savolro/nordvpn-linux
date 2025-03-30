/*
Package distro provides information about the current Linux distribution.
*/
package distro

// ReleaseName of the currently running distribution.
func ReleaseName() (string, error) {
	return "Windows", nil
}

// ReleasePrettyName of the currently running distribution.
func ReleasePrettyName() (string, error) {
	return "Windows", nil
}

// KernelName of the currently running kernel.
func KernelName() string { return "Windows" }

// KernelFull name of the currently running kernel.
func KernelFull() string { return "Windows" }
