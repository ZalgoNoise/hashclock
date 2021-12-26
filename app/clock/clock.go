package clock

import (
	"encoding/hex"
	"time"

	rsha "github.com/ZalgoNoise/meta/crypto/hash"

	"errors"
	"fmt"
)

type HashClockRequest struct {
	seed       string
	iterations int
	breakpoint int
	timeout    int
}

type HashClockResponse struct {
	Seed       string `json:"seed,omitempty"`
	Timeout    int    `json:"timeout,omitempty"`
	Iterations int    `json:"iterations,omitempty"`
	Hash       string `json:"hash,omitempty"`
}

type HashClockService struct {
	request  *HashClockRequest
	response *HashClockResponse
}

func NewService() *HashClockService {
	c := &HashClockService{}

	// initialize req & res
	c.request = &HashClockRequest{}
	c.response = &HashClockResponse{}

	// set default values
	c.request.iterations = 1
	c.request.breakpoint = 1
	c.request.timeout = 0

	return c
}

func (c *HashClockService) Hash(seed string) (*HashClockResponse, error) {
	// empty string exception
	if seed == "" {
		return &HashClockResponse{}, errors.New("seed cannot be empty")
	}

	c.request.seed = seed
	c.request.iterations = 1
	c.request.breakpoint = 1
	c.request.timeout = 0

	return newHashResponse(c.request)
}

func newHashResponse(req *HashClockRequest) (*HashClockResponse, error) {
	hash := hex.EncodeToString(rsha.Hash(req.seed))

	return &HashClockResponse{
		Seed:       req.seed,
		Timeout:    req.timeout,
		Iterations: req.iterations,
		Hash:       hash,
	}, nil
}

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

	return newRecHashResponse(c.request)
}

func newRecHashResponse(req *HashClockRequest) (*HashClockResponse, error) {
	var hash []byte

	// recursive SHA256 hash
	for i := 1; i <= req.iterations; i++ {
		if i == 1 {
			hash = rsha.Hash(req.seed)
		} else {
			hash = rsha.Hash(hash)
		}
	}

	return &HashClockResponse{
		Seed:       req.seed,
		Timeout:    req.timeout,
		Iterations: req.iterations,
		Hash:       hex.EncodeToString(hash),
	}, nil
}

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

	return newRecHashPrintResponse(c.request)
}

func newRecHashPrintResponse(req *HashClockRequest) (*HashClockResponse, error) {
	var hash []byte

	// recursive SHA256 hash
	for i := 1; i <= req.iterations; i++ {
		if i == 1 {
			hash = rsha.Hash(req.seed)
		} else {
			hash = rsha.Hash(hash)
		}

		// breakpoint logging
		if i%req.breakpoint == 0 {
			fmt.Printf("#%v:\t%x\n", i, hash)
		}
	}

	return &HashClockResponse{
		Seed:       req.seed,
		Timeout:    req.timeout,
		Iterations: req.iterations,
		Hash:       hex.EncodeToString(hash),
	}, nil
}

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

func (c *HashClockService) RecHashTime(seed string, timeout int) (*HashClockResponse, error) {
	// empty string exception
	if seed == "" {
		return &HashClockResponse{}, errors.New("seed cannot be empty")
	}

	if timeout <= 0 {
		return &HashClockResponse{}, errors.New("timeout cannot be zero or below")
	}

	c.request.seed = seed
	c.request.iterations = 0
	c.request.breakpoint = 0
	c.request.timeout = timeout

	return newRecHashTimeResponse(c.request)
}

func newRecHashTimeResponse(req *HashClockRequest) (*HashClockResponse, error) {
	r := &HashClockResponse{}
	r.Seed = req.seed
	r.Timeout = req.timeout

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

		ts.hash = rsha.Hash(req.seed)
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
	time.Sleep(time.Second * time.Duration(req.timeout))
	done <- true

	// get calculated hash and number of iterations
	r.Hash = hex.EncodeToString(ts.hash)
	r.Iterations = ts.id

	return r, nil
}
