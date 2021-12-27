package clock

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	// rhash "github.com/ZalgoNoise/meta/crypto/hash"
)

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

	c.request.seed = []byte(seed)
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
	hash := c.hasher.Hash(c.request.seed)
	target := []byte(c.request.hash)

	for {
		if iterations > 0 {
			hash = c.hasher.Hash(hash)
		}
		iterations++

		if matchHash(hash, target) {
			c.response = &HashClockResponse{
				Seed:       string(c.request.seed),
				Timeout:    c.request.timeout,
				Iterations: iterations,
				Hash:       string(hash),
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

	c.request.seed = []byte(seed)
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
		Seed:    string(c.request.seed),
		Timeout: c.request.timeout,
		Target:  c.request.hash,
	}
	target := []byte(c.request.hash)

	// recursive conversions are done with byte arrays
	// to preserve performance, instead of constantly
	// converting to string
	type timestamp struct {
		hash []byte
		id   int
	}

	ts := timestamp{
		hash: c.hasher.Hash(c.request.seed),
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
					ts.hash = c.hasher.Hash(ts.hash)
				}
				ts.id++

				if matchHash(ts.hash, target) {
					c.response.Iterations = ts.id
					c.response.Hash = string(ts.hash)
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
		Hash:       string(ts.hash),
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

	c.request.seed = []byte(seed)
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

	hash := c.hasher.Hash(c.request.seed)
	target := []byte(c.request.hash)

	// index starts at 2 since:
	// - index 0 is the seed
	// - index 1 is the first hash calculated (above)
	for i := 2; i <= c.request.iterations; i++ {
		hash = c.hasher.Hash(hash)
	}

	c.response = &HashClockResponse{
		Seed:       string(c.request.seed),
		Timeout:    c.request.timeout,
		Iterations: c.request.iterations,
		Hash:       string(hash),
		Target:     c.request.hash,
		Duration:   time.Since(timestamp),
	}

	if matchHash(hash, target) {
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
