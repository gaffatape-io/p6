package test

import (
	"flag"
)

var Flags = flag.NewFlagSet("test-flags", flag.PanicOnError)

func init() {
	Flags.Set("logtostderr", "1")
	Flags.Set("stderrthreshold", "INFO")
}
