package internal

import (
	"os"
	"os/signal"

	"golang.org/x/sys/windows"
)

func GetSignalChan() <-chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, windows.SIGTERM, windows.SIGHUP)
	return signals
}

// WaitSignal for app to shutdown
func WaitSignal() {
	signals := GetSignalChan()
	<-signals
}
