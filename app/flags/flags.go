package flags

import "flag"

type CLIConfig struct {
	ClockSeed  string
	Iterations int
	Breakpoint int
	Timeout    int
	SetJSON    bool
}

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
