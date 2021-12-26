package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ZalgoNoise/hashclock/flags"

	// "hashclock/flags"
	"github.com/ZalgoNoise/hashclock/clock"
	// "hashclock/clock"
)

func printResponse(res *clock.HashClockResponse, toJSON bool) {
	if toJSON {
		printJSON(res)
	}
	printText(res)
}

func printJSON(res *clock.HashClockResponse) {
	type output struct {
		Seed       string `json:"seed,omitempty"`
		Iterations int    `json:"iterations,omitempty"`
		Timeout    int    `json:"timeout,omitempty"`
		Hash       string `json:"hash,omitempty"`
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
	out, err := json.Marshal(o)

	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	os.Exit(0)
}

func printText(res *clock.HashClockResponse) {
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
	out += s + res.Seed + nl
	out += pad + nl + res.Hash + nl + pad + nl

	fmt.Println(out)
	os.Exit(0)
}

func Run() {
	cfg := flags.NewConfig()
	cService := clock.NewService()

	if cfg.ClockSeed == "" {
		fmt.Errorf("input seed string is undefined")
		os.Exit(1)
	}

	if cfg.Timeout > 0 {
		res, err := cService.RecHashTime(cfg.ClockSeed, cfg.Timeout)
		if err != nil {
			fmt.Errorf(err.Error())
			os.Exit(1)
		}

		printResponse(res, cfg.SetJSON)
	}

	if cfg.Iterations == 0 {
		err := cService.RecHashLoop(cfg.ClockSeed, cfg.Breakpoint)
		if err != nil {
			fmt.Errorf(err.Error())
			os.Exit(1)
		}

	}
	if cfg.Iterations == 1 {
		res, err := cService.Hash(cfg.ClockSeed)
		if err != nil {
			fmt.Errorf(err.Error())
			os.Exit(1)
		}

		printResponse(res, cfg.SetJSON)
	}
	if cfg.Iterations >= 2 {
		if cfg.Breakpoint == 0 {
			res, err := cService.RecHash(cfg.ClockSeed, cfg.Iterations)
			if err != nil {
				fmt.Errorf(err.Error())
				os.Exit(1)
			}

			printResponse(res, cfg.SetJSON)
		}

		res, err := cService.RecHashPrint(cfg.ClockSeed, cfg.Iterations, cfg.Breakpoint)
		if err != nil {
			fmt.Errorf(err.Error())
			os.Exit(1)
		}

		printResponse(res, cfg.SetJSON)
	}
}
