package flags

import "flag"

type PoHConfig struct {
	ClockSeed string
	Iterations int
	Breakpoint int
}

func NewConfig() *PoHConfig {
	inputClockSeed := flag.String("seed", "", "Input seed which will be hashed")
	inputIterations := flag.Int("iter", 1, "Number of iterations")
	inputBreakpoint := flag.Int("log", 1, "Log hashes every # of steps")

	flag.Parse()

	return &PoHConfig{
		ClockSeed: *inputClockSeed,
		Iterations: *inputIterations,
		Breakpoint: *inputBreakpoint,
	}
}