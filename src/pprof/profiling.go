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
	flags.StringVar(&pprofType, "pprof", "time", "统计")
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
		fmt.Printf("\n命令耗时 %f 秒。\n", duration.Seconds())
	default:
	}

	return nil
}
