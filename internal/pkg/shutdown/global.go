package shutdown

import "os"

var globalShutdown = New() //nolint: gochecknoglobals

func Global() *Stopper {
	return globalShutdown
}

func Register(toStop ...StopInterface) {
	globalShutdown.Register(toStop...)
}

func RegisterFn(toStop ...func()) {
	for _, stop := range toStop {
		globalShutdown.Register(StopFn(stop))
	}
}

func Wait(ch chan struct{}) {
	globalShutdown.Wait(ch)
}

func WaitSignal(ch chan os.Signal) {
	globalShutdown.WaitSignal(ch)
}

func Stop() {
	globalShutdown.Stop()
}
