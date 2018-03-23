// Package profile provides configuration-less profiling enabled with a command-line flag.
package profile // import "github.com/motki/core/app/profile"

import (
	"flag"

	"github.com/pkg/profile"
)

var profileMode = flag.String("profile", "", "Enable profiling. Writes profiler data to current directory.\nPossible values: cpu, mem, mutex, block, trace.")

type profiler struct {
	mode string
	s    interface {
		Stop()
	}
}

func (p *profiler) Kind() string {
	return p.mode
}

func (p *profiler) Stop() {
	p.s.Stop()
}

func New() *profiler {
	var fn func(*profile.Profile)
	switch *profileMode {
	case "cpu":
		fn = profile.CPUProfile
	case "mem":
		fn = profile.MemProfile
	case "mutex":
		fn = profile.MutexProfile
	case "block":
		fn = profile.BlockProfile
	case "trace":
		fn = profile.TraceProfile
	}
	if fn != nil {
		return &profiler{
			mode: *profileMode,
			s:    profile.Start(fn, profile.ProfilePath("."), profile.NoShutdownHook)}
	}
	return nil
}
