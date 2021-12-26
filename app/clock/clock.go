// Package clock contains the logic for the `HashClockService`,
// which defines several methods for (recursively) hashing a string.
//
// It is built generically to support multiple actions with a
// standard service:request/response structure
package clock

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	rsha "github.com/ZalgoNoise/meta/crypto/hash"
)

// HashClockRequest struct defines the input configuration for
// a `HashClockService` request
type HashClockRequest struct {
	seed       string
	iterations int
	breakpoint int
	timeout    int
	hash       string
}

// HashClockResponse struct defines the input configuration for
// a `HashClockService` response
type HashClockResponse struct {
	Seed       string `json:"seed,omitempty"`
	Timeout    int    `json:"timeout,omitempty"`
	Iterations int    `json:"iterations,omitempty"`
	Hash       string `json:"hash,omitempty"`
	Target     string `json:"target,omitempty"`
	Match      bool   `json:"match,omitempty"`
	Duration   time.Duration
}

// HashClockService struct is a placeholder for this service,
// containing a pointer to both the request and response objects,
// and being the container for all methods in this package
type HashClockService struct {
	request  *HashClockRequest
	response *HashClockResponse
}

// NewService function is a generic public function to spawn a
// pointer to a new HashClockService, with set default values
func NewService() *HashClockService {
	c := &HashClockService{}

	// initialize req
	c.request = &HashClockRequest{}

	// set default values
	c.request.iterations = 1
	c.request.breakpoint = 1
	c.request.timeout = 0

	return c
}

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

// Verify method will take in a seed string and a target hash,
// returning an execution of the `newVerifyResponse` method
func (c *HashClockService) Verify(seed string, hash string) (*HashClockResponse, error) {
	// empty string exception
	if seed == "" {
		return &HashClockResponse{}, errors.New("seed cannot be empty")
	}

	// empty hash exception
	if hash == "" {
		return &HashClockResponse{}, errors.New("hash cannot be empty")
	}

	// hash is not hex-encoded exception
	if _, err := hex.DecodeString(hash); err != nil {
		return &HashClockResponse{}, fmt.Errorf("hex encoder: invalid string -- %s", err)
	}

	// seed is hash exception
	if seed == hash {
		return &HashClockResponse{}, errors.New("seed cannot be the same as the hash (no verification involved)")
	}

	c.request.seed = seed
	c.request.iterations = 0
	c.request.breakpoint = 0
	c.request.timeout = 0
	c.request.hash = hash

	return c.newVerifyResponse()
}

// newVerifyResponse method will parse the `HashClockService.request` object
// and build its `HashClockResponse.response`; by recursively hashing the seed
// until it finds the target hash.
//
// This operation is infinitely recursive and will not be terminated unless
// halted by the user -- or, when the hash matches.
func (c *HashClockService) newVerifyResponse() (*HashClockResponse, error) {
	// timestamp is recorded when function is first called
	timestamp := time.Now()

	iterations := 0
	hash := rsha.Hash(c.request.seed)
	target := []byte(c.request.hash)
	enc := make([]byte, hex.EncodedLen(32))

	for {
		if iterations > 0 {
			hash = rsha.Hash(hash)
		}
		iterations++

		// hex-encode for comparison:
		hex.Encode(enc, hash)

		if matchHash(enc, target) {
			c.response = &HashClockResponse{
				Seed:       c.request.seed,
				Timeout:    c.request.timeout,
				Iterations: iterations,
				Hash:       string(enc),
				Target:     c.request.hash,
				Match:      true,
				Duration:   time.Since(timestamp),
			}

			return c.response, nil
		}

	}
}

// VerifyTimeout method will take in a seed string, a target hash,
// returning an execution of the `newVerifyTimeoutResponse` method
func (c *HashClockService) VerifyTimeout(seed, hash string, timeout int) (*HashClockResponse, error) {
	// empty string exception
	if seed == "" {
		return &HashClockResponse{}, errors.New("seed cannot be empty")
	}

	// empty hash exception
	if hash == "" {
		return &HashClockResponse{}, errors.New("hash cannot be empty")
	}

	// hash is not hex-encoded exception
	if _, err := hex.DecodeString(hash); err != nil {
		return &HashClockResponse{}, err
	}

	// seed is hash exception
	if seed == hash {
		return &HashClockResponse{}, errors.New("seed cannot be the same as the hash (no verification involved)")
	}

	// empty timeout exception
	if timeout <= 0 {
		return &HashClockResponse{}, errors.New("timeout cannot be zero or below")
	}

	c.request.seed = seed
	c.request.iterations = 0
	c.request.breakpoint = 0
	c.request.timeout = timeout
	c.request.hash = hash

	return c.newVerifyTimeoutResponse()
}

// newVerifyTimeoutResponse method will parse the `HashClockService.request` object
// and build its `HashClockResponse.response`; by recursively hashing the seed
// until it finds the target hash within a specific timeframe.
//
// This operation will stop with a match or when the timer is up. As the hashing is
// done in a goroutine, the ticker in the parent process will check for a match every
// 10ms. This does not affect performance when compared to 100ms, for instance.
func (c *HashClockService) newVerifyTimeoutResponse() (*HashClockResponse, error) {
	c.response = &HashClockResponse{
		Seed:    c.request.seed,
		Timeout: c.request.timeout,
		Target:  c.request.hash,
	}
	target := []byte(c.request.hash)
	enc := make([]byte, hex.EncodedLen(32))

	// recursive conversions are done with byte arrays
	// to preserve performance, instead of constantly
	// converting to string
	type timestamp struct {
		hash []byte
		id   int
	}

	ts := timestamp{
		hash: rsha.Hash(c.request.seed),
		id:   0,
	}
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				if ts.id > 0 {
					ts.hash = rsha.Hash(ts.hash)
				}
				ts.id++

				// hex-encode for comparison:
				hex.Encode(enc, ts.hash)

				if matchHash(enc, target) {
					c.response.Iterations = ts.id
					c.response.Hash = string(enc)
					c.response.Match = true

					return
				}
			}
		}
	}()

	// TODO: refactor; needs something more solid than this
	// kick off timer; then send done signal to goroutine.
	// Performs regular checks on the response object to avoid
	// waiting for the total length of the timeout
	//
	// Checking every 10ms seems not to be impactful (yet)
	for i := 0; i <= c.request.timeout*100; i++ {
		time.Sleep(time.Millisecond * 10)
		if c.response.Match {
			c.response.Duration = time.Millisecond * 10 * time.Duration(i)
			return c.response, nil
		}
	}
	done <- true

	c.response = &HashClockResponse{
		Iterations: ts.id,
		Hash:       string(enc),
		Match:      false,
	}

	return c.response, nil
}

// VerifyIndex method will take in a seed string, a target hash,
// returning an execution of the `newVerifyIndexResponse` method
func (c *HashClockService) VerifyIndex(seed string, hash string, iterations int) (*HashClockResponse, error) {

	// empty string exception
	if seed == "" {
		return &HashClockResponse{}, errors.New("seed cannot be empty")
	}

	// empty hash exception
	if hash == "" {
		return &HashClockResponse{}, errors.New("hash cannot be empty")
	}

	// hash is not hex-encoded exception
	if _, err := hex.DecodeString(hash); err != nil {
		return &HashClockResponse{}, err
	}

	// seed is hash exception
	if seed == hash {
		return &HashClockResponse{}, errors.New("seed cannot be the same as the hash (no verification involved)")
	}

	// iterations is zero or below exception
	if iterations <= 0 {
		return &HashClockResponse{}, errors.New("number of target iterations cannot be zero or below")
	}

	c.request.seed = seed
	c.request.iterations = iterations
	c.request.breakpoint = 0
	c.request.timeout = 0
	c.request.hash = hash

	return c.newVerifyIndexResponse()

}

// newVerifyIndexResponse method will parse the `HashClockService.request` object
// and build its `HashClockResponse.response`; by recursively hashing the seed
// a specific number of times.
//
// The resulting hash is matched to the target hash, and the results are returned.
func (c *HashClockService) newVerifyIndexResponse() (*HashClockResponse, error) {
	// timestamp is recorded when function is first called
	timestamp := time.Now()

	hash := rsha.Hash(c.request.seed)
	target := []byte(c.request.hash)
	enc := make([]byte, hex.EncodedLen(32))

	// index starts at 2 since:
	// - index 0 is the seed
	// - index 1 is the first hash calculated (above)
	for i := 2; i <= c.request.iterations; i++ {
		hash = rsha.Hash(hash)
	}

	// hex-encode for comparison:
	hex.Encode(enc, hash)

	c.response = &HashClockResponse{
		Seed:       c.request.seed,
		Timeout:    c.request.timeout,
		Iterations: c.request.iterations,
		Hash:       string(enc),
		Target:     c.request.hash,
		Duration:   time.Since(timestamp),
	}

	if matchHash(enc, target) {
		c.response.Match = true
		return c.response, nil
	}
	c.response.Match = false
	return c.response, nil
}

// matchHash function is a helper to read and compare each byte from both
// the input hash and the target hash
func matchHash(hash, target []byte) bool {
	for idx, t := range target {
		if t != hash[idx] {
			return false
		}
	}
	return true
}
