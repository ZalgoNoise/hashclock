package clock

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	rsha "github.com/ZalgoNoise/meta/crypto/hash"
)

// Hash method takes in a string to hash, returning an execution of the
// `newHashResponse` method
func (c *HashClockService) Hash(seed string) (*HashClockResponse, error) {
	// empty string exception
	if seed == "" {
		return &HashClockResponse{}, errors.New("seed cannot be empty")
	}

	c.request.seed = seed
	c.request.iterations = 1
	c.request.breakpoint = 1
	c.request.timeout = 0

	return c.newHashResponse()
}

// newHashResponse method will parse the `HashClockService.request` object
// and build its `HashClockResponse.response` with the hash for the seed string
func (c *HashClockService) newHashResponse() (*HashClockResponse, error) {
	hash := hex.EncodeToString(rsha.Hash(c.request.seed))

	c.response = &HashClockResponse{
		Seed:       c.request.seed,
		Timeout:    c.request.timeout,
		Iterations: c.request.iterations,
		Hash:       hash,
	}

	return c.response, nil
}

// RecHash method takes in a string to hash and the number of desired iterations,
// returning an execution of the `newRecHashResponse` method
func (c *HashClockService) RecHash(seed string, iter int) (*HashClockResponse, error) {
	// empty string exception
	if seed == "" {
		return &HashClockResponse{}, errors.New("seed cannot be empty")
	}

	// zero iterations exception
	if iter <= 0 {
		return &HashClockResponse{}, errors.New("number of iterations has to be greater than zero")
	}

	c.request.seed = seed
	c.request.iterations = iter
	c.request.breakpoint = 0
	c.request.timeout = 0

	return c.newRecHashResponse()
}

// newRecHashResponse method will parse the `HashClockService.request` object
// and build its `HashClockResponse.response`; by hashing the seed for the number
// of times defined in the iterations value, and setting them in the response object
func (c *HashClockService) newRecHashResponse() (*HashClockResponse, error) {
	var hash []byte

	// recursive SHA256 hash
	for i := 1; i <= c.request.iterations; i++ {
		if i == 1 {
			hash = rsha.Hash(c.request.seed)
		} else {
			hash = rsha.Hash(hash)
		}
	}

	c.response = &HashClockResponse{
		Seed:       c.request.seed,
		Timeout:    c.request.timeout,
		Iterations: c.request.iterations,
		Hash:       hex.EncodeToString(hash),
	}

	return c.response, nil
}

// RecHashPrint method takes in a string to hash, the number of desired iterations,
// and a breakpoint value; returning an execution of the `newRecHashResponse` method
func (c *HashClockService) RecHashPrint(seed string, iter int, breakpoint int) (*HashClockResponse, error) {
	// empty string exception
	if seed == "" {
		return &HashClockResponse{}, errors.New("seed cannot be empty")
	}

	// zero iterations exception
	if iter <= 0 {
		return &HashClockResponse{}, errors.New("number of iterations has to be greater than zero")
	}

	// negative breakpoint exception
	if breakpoint < 0 {
		return &HashClockResponse{}, errors.New("logging frequency cannot be negative")
	}

	// zero breakpoint exception
	if breakpoint == 0 {
		return &HashClockResponse{}, errors.New("invalid function call -- HashClockService.RecHashPrint() needs to be called with a breakpoint > 0. method RecHash() should be used instead")
	}

	c.request.seed = seed
	c.request.iterations = iter
	c.request.breakpoint = breakpoint
	c.request.timeout = 0

	return c.newRecHashPrintResponse()
}

// newRecHashPrintResponse method will parse the `HashClockService.request` object
// and build its `HashClockResponse.response`; by hashing the seed for the number
// of times defined in the iterations value, and setting them in the response object.
//
// During execution, if the counter modulo breakpoint is zero (counter % breakpoint == 0),
// the hash is printed to std-out.
func (c *HashClockService) newRecHashPrintResponse() (*HashClockResponse, error) {
	var hash []byte

	// recursive SHA256 hash
	for i := 1; i <= c.request.iterations; i++ {
		if i == 1 {
			hash = rsha.Hash(c.request.seed)
		} else {
			hash = rsha.Hash(hash)
		}

		// breakpoint logging
		if i%c.request.breakpoint == 0 {
			fmt.Printf("#%v:\t%x\n", i, hash)
		}
	}

	c.response = &HashClockResponse{
		Seed:       c.request.seed,
		Timeout:    c.request.timeout,
		Iterations: c.request.iterations,
		Hash:       hex.EncodeToString(hash),
	}

	return c.response, nil
}

// RecHashLoop method will recursively hash the seed string, infinitely
// (or until the program is halted) while printing out its hashes.
//
// During execution, if the counter modulo breakpoint is zero (counter % breakpoint == 0),
// the hash is printed to std-out.
//
// This means that a breakpoint of 1 prints every hash, while a breakpoint of 5 prints
// every 5th hash.
//
// Does not return a response object since it will be an infinite loop until the program
// is interrupted and/or killed; only an error in case the input values are invalid
func (c *HashClockService) RecHashLoop(seed string, breakpoint int) error {
	// empty string exception
	if seed == "" {
		return errors.New("seed cannot be empty")
	}

	// negative breakpoint exception
	if breakpoint <= 0 {
		return errors.New("logging frequency cannot be zero or below")
	}

	c.request.seed = seed
	c.request.iterations = 0
	c.request.breakpoint = breakpoint
	c.request.timeout = 0

	hash := rsha.Hash(seed)

	var counter int = 1

	for {
		hash = rsha.Hash(hash)
		counter++

		// breakpoint logging
		if counter%breakpoint == 0 {
			fmt.Printf("#%v:\t%x\n", counter, hash)
		}
	}
}

// RecHashTime method will take in a seed string and a timeout value (in seconds),
// returning an execution of the `newRecHashTimeResponse` method
func (c *HashClockService) RecHashTimeout(seed string, timeout int) (*HashClockResponse, error) {
	// empty string exception
	if seed == "" {
		return &HashClockResponse{}, errors.New("seed cannot be empty")
	}

	// empty timeout exception
	if timeout <= 0 {
		return &HashClockResponse{}, errors.New("timeout cannot be zero or below")
	}

	c.request.seed = seed
	c.request.iterations = 0
	c.request.breakpoint = 0
	c.request.timeout = timeout

	return c.newRecHashTimeoutResponse()
}

// newRecHashTimeResponse method will parse the `HashClockService.request` object
// and build its `HashClockResponse.response`; by continuously hashing the seed string
// in a goroutine limited by a timer (in seconds).
//
// Once the timer runs out, a created `done` channel interrupts the goroutine. The
// calculated hash and number of iterations are parsed into the `HashClockResponse.response`
// object
func (c *HashClockService) newRecHashTimeoutResponse() (*HashClockResponse, error) {
	r := &HashClockResponse{}
	r.Seed = c.request.seed
	r.Timeout = c.request.timeout

	// recursive conversions are done with byte arrays
	// to preserve performance, instead of constantly
	// converting to string
	type timestamp struct {
		hash []byte
		id   int
	}

	ts := timestamp{}
	done := make(chan bool)

	// recursively calculate hashes until timer is up
	go func() {

		ts.hash = rsha.Hash(c.request.seed)
		ts.id = 1

		for {
			select {
			case <-done:
				return
			default:
				ts.hash = rsha.Hash(ts.hash)
				ts.id++
			}
		}

	}()

	// kick off timer; then send done signal to goroutine
	time.Sleep(time.Second * time.Duration(c.request.timeout))
	done <- true

	// get calculated hash and number of iterations
	r.Hash = hex.EncodeToString(ts.hash)
	r.Iterations = ts.id

	return r, nil
}
