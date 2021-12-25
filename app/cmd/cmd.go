package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ZalgoNoise/hashclock/flags"
	rsha "github.com/ZalgoNoise/meta/crypto/hash"

	// "hashclock/flags"
	"github.com/ZalgoNoise/hashclock/clock"
	// "hashclock/clock"
)

type output struct {
	Seed       string `json:"seed,omitempty"`
	Iterations int    `json:"iterations,omitempty"`
	Timeout    int    `json:"timeout,omitempty"`
	Hash       string `json:"hash,omitempty"`
}

func printJSON(seed string, hash string, iterations int, timeout int) {
	o := &output{}

	o.Seed = seed
	o.Hash = hash

	if iterations > 0 {
		o.Iterations = iterations
	}

	if timeout > 0 {
		o.Timeout = timeout
	}
	out, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	os.Exit(0)
}

func printText(seed string, hash string, iterations int, timeout int) {
	const (
		pad string = "----"
		tb  string = "time: "
		tas string = " second"
		tam string = " seconds"
		i   string = "hashes: "
		s   string = "seed: "
		sp  string = "; "
		nl  string = "\n"
	)

	out := nl + pad + nl

	if timeout > 0 {
		if timeout == 1 {
			out += tb + strconv.Itoa(timeout) + tas + sp
		} else {
			out += tb + strconv.Itoa(timeout) + tam + sp
		}
	}
	if iterations > 0 {
		out += i + strconv.Itoa(iterations) + sp
	}
	out += s + seed + nl
	out += pad + nl + hash + nl + pad + nl

	fmt.Println(out)
	os.Exit(0)
}

func Run() {
	cfg := flags.NewConfig()

	if cfg.ClockSeed == "" {
		fmt.Errorf("input seed string is undefined")
		os.Exit(1)
	}

	if cfg.Timeout > 0 {
		count, hash := clock.RecursiveSHA256Timed(cfg.ClockSeed, cfg.Timeout)

		if cfg.SetJSON {
			printJSON(cfg.ClockSeed, hex.EncodeToString(hash), count, cfg.Timeout)
		}
		printText(cfg.ClockSeed, hex.EncodeToString(hash), count, cfg.Timeout)

	}

	if cfg.Iterations == 0 {
		clock.RecursiveSHA256Inf(cfg.ClockSeed, cfg.Breakpoint)
	}
	if cfg.Iterations == 1 {
		hash := rsha.Hash(cfg.ClockSeed)

		if cfg.SetJSON {
			printJSON(cfg.ClockSeed, hex.EncodeToString(hash), cfg.Iterations, cfg.Timeout)
		}
		printText(cfg.ClockSeed, hex.EncodeToString(hash), cfg.Iterations, cfg.Timeout)

	}
	if cfg.Iterations >= 2 {
		if cfg.Breakpoint == 0 {
			hash, err := clock.RecursiveSHA256(cfg.ClockSeed, cfg.Iterations)
			if err != nil {
				fmt.Errorf("error while hashing seed: %v", err)
				os.Exit(1)
			}
			if cfg.SetJSON {
				printJSON(cfg.ClockSeed, hex.EncodeToString(hash), cfg.Iterations, cfg.Timeout)
			}
			printText(cfg.ClockSeed, hex.EncodeToString(hash), cfg.Iterations, cfg.Timeout)

		}

		hash, err := clock.RecursiveSHA256Logged(cfg.ClockSeed, cfg.Iterations, cfg.Breakpoint)
		if err != nil {
			fmt.Errorf("error while hashing seed: %v", err)
			os.Exit(1)
		}
		if cfg.SetJSON {
			printJSON(cfg.ClockSeed, hex.EncodeToString(hash), cfg.Iterations, cfg.Timeout)
		}
		printText(cfg.ClockSeed, hex.EncodeToString(hash), cfg.Iterations, cfg.Timeout)
	}
}
