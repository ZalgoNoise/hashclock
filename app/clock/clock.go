// Package clock contains the logic for the `HashClockService`,
// which defines several methods for (recursively) hashing a string.
//
// It is built generically to support multiple actions with a
// standard service:request/response structure
package clock

import (
	"time"

	rhash "github.com/ZalgoNoise/meta/crypto/hash"
)

// HashClockRequest struct defines the input configuration for
// a `HashClockService` request
type HashClockRequest struct {
	seed       []byte
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
	hasher   rhash.Hasher
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

	// initialize hasher
	c.hasher = &rhash.SHA256{}

	return c
}
