package pprof

import (
	"fmt"
	"github.com/spf13/pflag"
	"time"
)

var (
	pprofType string
	startTime time.Time
	stopTime  time.Time
)

func AddProfilingFlags(flags *pflag.FlagSet) {
	flags.StringVar(&pprofType, "pprof", "time", "Type of prof capture.")
}

func InitProfiling() error {
	switch pprofType {
	case "none":
		return nil
	case "time":
		startTime = time.Now()
		fallthrough
	default:
	}

	return nil
}

func FlushProfiling() error {
	switch pprofType {
	case "none":
		return nil
	case "time":
		stopTime = time.Now()
		duration := stopTime.Sub(startTime)
		fmt.Printf("\nThis command use %f seconds.\n", duration.Seconds())
	default:
	}

	return nil
}
