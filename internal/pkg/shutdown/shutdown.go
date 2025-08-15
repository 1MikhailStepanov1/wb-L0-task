package shutdown

import "os"

type StopFn func()

func (s StopFn) Stop() {
	s()
}

type StopInterface interface {
	Stop()
}

type Stopper struct {
	stops []StopFn
}

func New() *Stopper {
	return &Stopper{}
}

func (s *Stopper) Register(toStop ...StopInterface) {
	for _, stop := range toStop {
		if stop != nil {
			s.stops = append(s.stops, stop.Stop)
		}
	}
}

func (s *Stopper) Wait(ch chan struct{}) {
	<-ch
	s.Stop()
}

func (s *Stopper) WaitSignal(ch chan os.Signal) {
	<-ch
	s.Stop()
}

func (s *Stopper) Stop() {
	for _, stop := range s.stops {
		stop.Stop()
	}
}
