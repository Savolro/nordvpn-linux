package internal

import (
	"log"
	"net"
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

// systemDFile returns a `os.systemDFile` object for
// systemDFile descriptor passed to this process via systemd fd-passing protocol.
//
// The order of the systemDFile descriptors is preserved in the returned slice.
// `unsetEnv` is typically set to `true` in order to avoid clashes in
// fd usage and to avoid leaking environment flags to child processes.
func systemDFile(unsetEnv bool) *os.File {
	defer func() {
		if unsetEnv {
			if err := os.Unsetenv(ListenPID); err != nil {
				log.Println(DeferPrefix, err)
			}
			if err := os.Unsetenv(ListenFDS); err != nil {
				log.Println(DeferPrefix, err)
			}
			if err := os.Unsetenv(ListenFDNames); err != nil {
				log.Println(DeferPrefix, err)
			}
		}
	}()

	pid, err := strconv.Atoi(os.Getenv(ListenPID))
	if err != nil || pid != os.Getpid() {
		return nil
	}

	nfds, err := strconv.Atoi(os.Getenv(ListenFDS))
	if err != nil || nfds != 1 {
		return nil
	}

	unix.CloseOnExec(listenFdsStart)
	name := os.Getenv(ListenFDNames)

	return os.NewFile(listenFdsStart, name)
}

// SystemDListener returns systemd defined, socket activated listener
func SystemDListener() (net.Listener, error) {
	file := systemDFile(true)
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(DeferPrefix, err)
		}
	}()
	return net.FileListener(file)
}
