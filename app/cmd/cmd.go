package cmd

import (
	"fmt"
	"os"

	"github.com/ZalgoNoise/hashclock/flags"
	rsha "github.com/ZalgoNoise/meta/crypto/hash"

	// "hashclock/flags"
	"github.com/ZalgoNoise/hashclock/clock"
	// "hashclock/clock"
)

func Run() {
	cfg := flags.NewConfig()

	if cfg.ClockSeed == "" {
		fmt.Errorf("input seed string is undefined")
		os.Exit(1)
	}

	if cfg.Timeout > 0 {
		count, hash := clock.RecursiveSHA256Timed(cfg.ClockSeed, cfg.Timeout)

		fmt.Printf(
			"----\ntime: %v seconds; hashes: %v; seed: %s\n----\n%x\n----\n",
			cfg.Timeout,
			count,
			cfg.ClockSeed,
			hash,
		)
		os.Exit(0)
	}

	if cfg.Iterations == 0 {
		clock.RecursiveSHA256Inf(cfg.ClockSeed, cfg.Breakpoint)
	}
	if cfg.Iterations == 1 {
		hash := rsha.Hash(cfg.ClockSeed)
		fmt.Printf("%x\n", hash)

	}
	if cfg.Iterations >= 2 {
		if cfg.Breakpoint == 0 {
			hash, err := clock.RecursiveSHA256(cfg.ClockSeed, cfg.Iterations)
			if err != nil {
				fmt.Errorf("error while hashing seed: %v", err)
				os.Exit(1)
			}
			fmt.Printf("#%v:\t\t%x\n", cfg.Iterations, hash)
		}

		hash, err := clock.RecursiveSHA256Logged(cfg.ClockSeed, cfg.Iterations, cfg.Breakpoint)
		if err != nil {
			fmt.Errorf("error while hashing seed: %v", err)
			os.Exit(1)
		}
		fmt.Printf("----\n#%v:\t\t%x\n", cfg.Iterations, hash)

	}
}
