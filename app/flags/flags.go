// Package flags registers the user's command-line input
// as an object, which is parsed and used in `hashclock/cmd`
package flags

import "flag"

// CLIConfig struct defines the set configuration for hashclock
// in an object which is parsed and used in `hashclock/cmd`
type CLIConfig struct {
	ClockSeed  string
	Iterations int
	Breakpoint int
	Timeout    int
	SetJSON    bool
}

// NewConfig function captures the set command-line flags and their
// values, and stores them in a `CLIConfig` object
func NewConfig() *CLIConfig {
	inputClockSeed := flag.String("seed", "", "Input seed which will be hashed")
	inputIterations := flag.Int("iter", 1, "Number of iterations")
	inputBreakpoint := flag.Int("log", 1, "Log hashes every # of steps")
	inputTimeout := flag.Int("time", 0, "Calculate hashes for # seconds")
	inputSetJSON := flag.Bool("json", false, "Returns the output in JSON format")

	flag.Parse()

	return &CLIConfig{
		ClockSeed:  *inputClockSeed,
		Iterations: *inputIterations,
		Breakpoint: *inputBreakpoint,
		Timeout:    *inputTimeout,
		SetJSON:    *inputSetJSON,
	}
}
