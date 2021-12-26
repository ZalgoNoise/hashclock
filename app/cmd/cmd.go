// Package cmd will contain the logic for a CLI deployment
// of hashclock
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ZalgoNoise/hashclock/clock"
	"github.com/ZalgoNoise/hashclock/flags"
)

// printResponse function is a generic fork based on the input `toJSON` value
// to either run `printJSON` or `printText`. This is to avoid repetition in `Run`
func printResponse(res *clock.HashClockResponse, toJSON bool) {
	if toJSON {
		printJSON(res)
	}
	printText(res)
}

// printJSON function will parse the set values in `clock.HashClockResponse`
// and build a new JSON object only containing the set values
func printJSON(res *clock.HashClockResponse) {
	type output struct {
		Seed       string `json:"seed,omitempty"`
		Iterations int    `json:"iterations,omitempty"`
		Timeout    int    `json:"timeout,omitempty"`
		Hash       string `json:"hash,omitempty"`
		Target     string `json:"target,omitempty"`
		Match      bool   `json:"match,omitempty"`
		Duration   string `json:"duration,omitempty"`
	}

	o := &output{}

	o.Seed = res.Seed
	o.Hash = res.Hash

	if res.Iterations > 0 {
		o.Iterations = res.Iterations
	}

	if res.Timeout > 0 {
		o.Timeout = res.Timeout
	}

	if res.Target != "" {
		o.Target = res.Target
		o.Match = res.Match
		o.Duration = res.Duration.String()
	}

	out, err := json.Marshal(o)

	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	os.Exit(0)
}

// printText function will parse the set values in `clock.HashClockResponse`
// and build a std-out message only containing the set values
func printText(res *clock.HashClockResponse) {
	const (
		pad string = "----"
		tb  string = "timeout: "
		tas string = " second"
		tam string = " seconds"
		i   string = "hashes: "
		s   string = "seed: "
		t   string = "target: "
		m   string = "match: "
		d   string = "duration: "
		sp  string = "; "
		nl  string = "\n"
	)

	out := nl + pad + nl

	if res.Timeout > 0 {
		if res.Timeout == 1 {
			out += tb + strconv.Itoa(res.Timeout) + tas + sp
		} else {
			out += tb + strconv.Itoa(res.Timeout) + tam + sp
		}
	}
	if res.Iterations > 0 {
		out += i + strconv.Itoa(res.Iterations) + sp
	}
	out += s + res.Seed + sp

	if res.Target != "" {
		out += t + res.Target + sp + m + strconv.FormatBool(res.Match) + sp
	}

	if res.Duration > 0 {
		out += d + res.Duration.String() + sp
	}

	out += nl + pad + nl + res.Hash + nl + pad + nl

	fmt.Println(out)
	os.Exit(0)
}

// Run function is the entrypoint for a CLI deployment of hashclock
//
// The configuration is defined from the parsed flags, and a new
// `clock.HashClockService` is created. From this point, depending on the
// defined values, different methods are called.
//
// All `clock.HashClockService` methods (except for `RecHashLoop`) return
// a `clock.HashClockResponse` object, which is parsed in the
// `printResponse` function
func Run() {
	cfg := flags.NewConfig()
	cService := clock.NewService()

	if cfg.Seed == "" {
		fmt.Println("input seed string is undefined")
		os.Exit(1)
	}

	// hash is set
	// verify hashes from seed
	if cfg.Hash != "" {
		// timeout is also set
		// verify hash within # of seconds
		if cfg.Timeout > 0 {
			res, err := cService.VerifyTimeout(cfg.Seed, cfg.Hash, cfg.Timeout)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			printResponse(res, cfg.SetJSON)
		}

		// iterations is set
		// verify hashes with a certain index; default value
		// for iterations in CLI is 1, so detector starts at 2
		if cfg.Iterations > 1 {
			res, err := cService.VerifyIndex(cfg.Seed, cfg.Hash, cfg.Iterations)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			printResponse(res, cfg.SetJSON)
		}

		// simple verifier (runs indefinitely)
		res, err := cService.Verify(cfg.Seed, cfg.Hash)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		printResponse(res, cfg.SetJSON)
	}

	// timeout is set
	// calculate hashes for # seconds
	if cfg.Timeout > 0 {
		res, err := cService.RecHashTimeout(cfg.Seed, cfg.Timeout)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		printResponse(res, cfg.SetJSON)
	}

	// 0 iterations set
	// infinite recursive hashing
	if cfg.Iterations == 0 {
		err := cService.RecHashLoop(cfg.Seed, cfg.Breakpoint)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	}

	// 1 iteration set
	// clac only 1 hash of a seed string
	if cfg.Iterations == 1 {
		res, err := cService.Hash(cfg.Seed)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		printResponse(res, cfg.SetJSON)
	}

	// 2+ iterations set
	// recursive hashing
	if cfg.Iterations >= 2 {

		// breakpoint is 0
		// don't print calculated hashes
		if cfg.Breakpoint == 0 {
			res, err := cService.RecHash(cfg.Seed, cfg.Iterations)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			printResponse(res, cfg.SetJSON)
		}

		// breakpoint is 1+
		// log every X hashes
		res, err := cService.RecHashPrint(cfg.Seed, cfg.Iterations, cfg.Breakpoint)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		printResponse(res, cfg.SetJSON)
	}
}
