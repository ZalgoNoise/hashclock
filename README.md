# hashclock

[![Bazel-CI](https://github.com/ZalgoNoise/hashclock/actions/workflows/bazel-ci.yaml/badge.svg)](https://github.com/ZalgoNoise/hashclock/actions/workflows/bazel-ci.yaml)
[![Dockerfile-CI](https://github.com/ZalgoNoise/hashclock/actions/workflows/docker-build.yaml/badge.svg)](https://github.com/ZalgoNoise/hashclock/actions/workflows/docker-build.yaml)
[![Go](https://github.com/ZalgoNoise/hashclock/actions/workflows/go.yaml/badge.svg)](https://github.com/ZalgoNoise/hashclock/actions/workflows/go.yaml)
[![Bazel-CD](https://github.com/ZalgoNoise/hashclock/actions/workflows/bazel-cd.yaml/badge.svg)](https://github.com/ZalgoNoise/hashclock/actions/workflows/bazel-cd.yaml)

_______


### Hashclock

_A recursive SHA256 hash function to print out cryptographic timestamps, or to benchmark CPU performance_

This repo contains a Golang binary and library to perform recursive hashing on raw data. The goal is to reduce any processing overhead as the library introduces a `Hasher` interface, compatible with the most popular hashing functions, which is described by a `Hash` method which must be recursive.

With these foundations, all the features and functionalities `hashclock` has (and will have) are built with such libraries. It allows for a simple command-line utility to calculate hashes recursively, for a certain amount of cycles, certain amount of time, or until it matches another input hash (for example).

The original concept for a hashclock emerged from studying different forms of transaction verification in cryptocurrencies, where the concept of [Proof of History](https://solana.com/solana-whitepaper.pdf) describes hashes as timestamps. In a nutshell, a (controlled; resource-limited) system continuously hashes a data object containing the number of rehashing cycles it has gone through, and if there is any input data (like a transaction), this newly joined data will force the stream of hashes to become different. Verification nodes will then verify if (and when) did the transaction occur, ensuring that from that point forward the resulting hashes will match:

Index | Operation | Output Hash
:----:|:---------:|:-----------:
1     | `sha256("genesis data")` | hash1
200   | `sha256(hash199)` | hash200
300   | `sha256(hash299)` | hash300
336   | `sha256(append(hash335, input_data_sha256))` | hash336

This repo does not attempt to improve or replicate this exact system, but merely as for research and analysis of recursive hashing functions, and their (extended) applications (apart from cryptocurrencies). 

_______

### Getting `hashclock`

You're able to get / build the `hashclock` binary in multiple ways, and for multiple targets:
- For all executables for most platforms, see the [Releases page](https://github.com/ZalgoNoise/hashclock/releases/latest)
- As Docker images, published in both `ghcr.io` and `docker.io` (see the [Releases page](https://github.com/ZalgoNoise/hashclock/releases/latest)):
  - `ghcr.io`: `docker pull ghcr.io/zalgonoise/hashclock`
  - `docker.io`: `docker pull zalgonoise/hashclock`
- Source code:
  - see the [Releases page](https://github.com/ZalgoNoise/hashclock/releases/latest) for stable versions
  - clone the repo for specific versions, or for a bleeding-edge version.

_______

### Building

> This section covers how to build a `hashclock` binary


#### Golang

You can build an executable with Golang's compiler, by changing directories to `hashclock/app`, and running `go build .`. Note that the `.gitignore` file is excluding the contents of the `hashclock/build` folder, which I use to export binaries to when using the native compiler. So, you can create an executable pointing to this folder with the command below:

```
cd app \
&& go build -o ../build/hashclock .

```

...Or simply build it in place:

```
cd app \
&& go build .
```

#### Bazel

There are plenty of Bazel targets in the `hashclock/app/BUILD` file for tooling, Go binary / library and also CI targets (`container_push` targets, etc). This section will focus on building `hashclock`.

To allow more _freedom_ in the root-level folder, `hashclock`'s logic has been moved to `hashclock/app`. In this folder you will also find the `WORKSPACE` file for the project, and as such you will need to first change-directories to this folder.

Within `hashclock/app`, you can build (or run, or test) `hashclock` with a simple `bazel build` command:

```
bazel build //:hashclock
```

________________

### Runtime / Command-line reference

The binary / executable is a command-line interface for `hashclock` which contains several features and applications. Taking note of the command-line flags:

```
Usage of hashclock:
  -alg string
        Hash function to use; lower-case or uppercase. One of: 'md5', 'sha1', 'sha224', 'sha256', 'sha384', 'sha512', 'sha512_224', 'sha512_256' (default "sha256")
  -hash string
        Input hash which will be verified, from hashing the seed
  -iter int
        Number of iterations (default 1)
  -json
        Returns the output in JSON format
  -log int
        Log hashes every # of steps (default 1)
  -seed string
        Input seed which will be hashed
  -time int
        Calculate hashes for # seconds
```

While there are plenty of flags and options, there are three main runtime modes to allow hashing / recursively hashing input data:
- for a (defined) number of iterations.
- for a certain amount of time (in seconds).
- until it matches another input hash string, with or without a timeout.

You can use the `-log {int}` flag to set an interval of when each calculated hash is printed. Printing all hashes (`-log 1`) will cause a lot of overhead and slower hashing rates. Printing no hashes (before the result) with `-log 0` is the most performant option. This flag is a modulo of the current index, so if you're printing 100 hashes with `-log 10`, it will print every 10th hash (`if idx % log_rate == 0`).

You can use the `-json` flag if you wish to further parse the resulting data in a JSON format.

Taking these modes as examples, please note below examples to these modes, when running the executable:

__Hash a string 1000000 times__

```
hashclock -seed "genesis_string" -iter 1000000 -log 0 -json | jq

{
  "seed": "genesis_string",
  "iterations": 1000000,
  "hash": "9fd84bb7d0d21dc6269590fc85aaf6ea564c16fd3b3c69a113c32c34fa60770c",
  "algorithm": "SHA256"
}
```

__Hash a string 10 times, printing all hashes__

```
hashclock -seed "genesis_string" -iter 10

#1:     32241cdef87d3717742dd16684f3ec711996233de3c9e8673c41939d68488d27
#2:     05d64f8e6ccd3810fbee93f6ee4120fd5a7061dddc2797dd1997a765a00e5006
#3:     d971baf34116ecb1bd23d9375baecae5d87a48ba3f09f76145ebed04986fb686
#4:     1adc3e2f51a8acceec2f4d9ba7d87d742ca545e927b8b2185b4f0c2e3238b238
#5:     d4961d23a08dccdfed9b6123a4aa951d1910fdecf15226cb97efab3814a17e36
#6:     2fe1298b76db3742fa8605c6c8032b3184e46da5babbc7969c473ee6e70d45e6
#7:     4ea68e5ad21d9021924d76e2029f52c6a9bf3d3d97815a6fabdfb8ccc32c8f5e
#8:     5db9827b7eeb6e100517b1b17f5d22cc2e8b41dbe10864491e24de51b8df1cb2
#9:     d8ad2a40b89b250daae9ea89ffcbb7896b74e7fdf56a873d1c6b28bfadbcc1a3
#10:    3a047a31ec5aa7ade3de9b013eefacc5d380f5133711d800fd8bcfcad670cbff

----
hashes: 10; seed: genesis_string; algo: SHA256; 
----
3a047a31ec5aa7ade3de9b013eefacc5d380f5133711d800fd8bcfcad670cbff
----
```

__Hash a string recursively for 10 seconds, most performant__

```
hashclock -seed "genesis_string" -log 0 -time 10 -json | jq
{
  "seed": "genesis_string",
  "iterations": 11458749,
  "timeout": 10,
  "hash": "d4b0aa00a3837fbbe18a2909b984b01f3c452427b59646c48008699dcefe3037",
  "algorithm": "SHA256"
}
```

__Recursively hash a string indefinitely, printing every 1000th hash (must close with Ctl+C)__

```
hashclock -seed "genesis_string" -iter 0 -log 1000

#1000:  6b63a2006f150e057e7a3ee69ecc6a607a2e6f3a3067ffb62d284aaf5a55d730
#2000:  60d856b4dbb41c864a41958d769917e29dabd170f2abbf2cfc2a8b346b7bb273
#3000:  f8f74404d329a7d855a8c32d0f6a10d24a8b57222b1b968314b2198bc99101f2
#4000:  db7f68cac04878aaf6ea646f5e1adbe78afba5503deea5a66b4e70f3d62b8513
#5000:  895e3570f7b9252d6fc6c6a63657913d11db619934947a635e2a2f1edb61980c
#6000:  a7c89f62c595af50281d87b77f07e12d9d5dbcee9e7f90f758ea7b98d42857f0
#7000:  540ad854e522f91f8890ee7ae83a1881675f60668dfb3a6403c29f7a03378d94
#8000:  51623c2b218dfe587195c0ea22d38825b8e73e759e73c1e3c8d024956a9cfd55
#9000:  f716660a6ddb3e3cd80be9b3a1be704b03cc474b1ce3f7a49be1007bcee08dea
(...)
#809000:        f0ad19c59fa470d7cd8dee9500a8fff4ee76fd20f029b0eb2d49e4eb1d1932fa
#810000:        4fb319123a127ce33ce8f9bf169e95c9e1cf5c7c91c726013c5fa5e9ff9fa5a0
```

__Recursively hash a string until it matches an input hash (without timeout)__

```
hashclock -seed "genesis_string" -hash 4fb319123a127ce33ce8f9bf169e95c9e1cf5c7c91c726013c5fa5e9ff9fa5a0 -json | jq

{
  "seed": "genesis_string",
  "iterations": 810000,
  "hash": "4fb319123a127ce33ce8f9bf169e95c9e1cf5c7c91c726013c5fa5e9ff9fa5a0",
  "target": "4fb319123a127ce33ce8f9bf169e95c9e1cf5c7c91c726013c5fa5e9ff9fa5a0",
  "match": true,
  "duration": "509.11935ms",
  "algorithm": "SHA256"
}
```

__Recursively hash a string until it matches an input hash (with timeout)__

```
hashclock -seed "genesis_string" -hash 4fb319123a127ce33ce8f9bf169e95c9e1cf5c7c91c726013c5fa5e9ff9fa5a0 -json -time 1 | jq

{
  "seed": "genesis_string",
  "iterations": 810000,
  "timeout": 1,
  "hash": "4fb319123a127ce33ce8f9bf169e95c9e1cf5c7c91c726013c5fa5e9ff9fa5a0",
  "target": "4fb319123a127ce33ce8f9bf169e95c9e1cf5c7c91c726013c5fa5e9ff9fa5a0",
  "match": true,
  "duration": "560ms",
  "algorithm": "SHA256"
}
```

________________________

#### Runtime with Bazel

Running the `hashclock` executable with `bazel` is very straight-forward, and all options / flags are of course available, for example:

```
blaze run //:hashclock -- -seed "genesis_string" -iter 10 -log 1

INFO: Invocation ID: {redacted}
INFO: Streaming build results to: https://app.buildbuddy.io/invocation/{redacted}
INFO: Analyzed target //:hashclock (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
Target //:hashclock up-to-date:
  bazel-bin/hashclock_/hashclock
INFO: Elapsed time: 0.823s, Critical Path: 0.01s
INFO: 1 process: 1 internal.
INFO: Running command line: bazel-bin/hashclock_/hashclock -seed genesis_string -iter 10 -log 1
INFO: Streaming build results to: https://app.buildbuddy.io/invocation/{redacted}
INFO: Build completed successfully, 1 total action
Waiting for build events upload: Build Event Service 1s

#1:     32241cdef87d3717742dd16684f3ec711996233de3c9e8673c41939d68488d27
#2:     05d64f8e6ccd3810fbee93f6ee4120fd5a7061dddc2797dd1997a765a00e5006
#3:     d971baf34116ecb1bd23d9375baecae5d87a48ba3f09f76145ebed04986fb686
#4:     1adc3e2f51a8acceec2f4d9ba7d87d742ca545e927b8b2185b4f0c2e3238b238
#5:     d4961d23a08dccdfed9b6123a4aa951d1910fdecf15226cb97efab3814a17e36
#6:     2fe1298b76db3742fa8605c6c8032b3184e46da5babbc7969c473ee6e70d45e6
#7:     4ea68e5ad21d9021924d76e2029f52c6a9bf3d3d97815a6fabdfb8ccc32c8f5e
#8:     5db9827b7eeb6e100517b1b17f5d22cc2e8b41dbe10864491e24de51b8df1cb2
#9:     d8ad2a40b89b250daae9ea89ffcbb7896b74e7fdf56a873d1c6b28bfadbcc1a3
#10:    3a047a31ec5aa7ade3de9b013eefacc5d380f5133711d800fd8bcfcad670cbff

----
hashes: 10; seed: genesis_string; algo: SHA256; 
----
3a047a31ec5aa7ade3de9b013eefacc5d380f5133711d800fd8bcfcad670cbff
----
```

________________

### Development and integration

Integrating `hashclock` or recursive hashing in your project may involve one of two routes:

- Getting simply the generic / abstract `Hasher` interface [from github.com/ZalgoNoise/meta/crypto/hash](https://github.com/ZalgoNoise/meta/tree/master/crypto/hash)
- Getting the `HashClockService` module [from github.com/ZalgoNoise/hashclock/app/clock](https://github.com/ZalgoNoise/hashclock/tree/master/app/clock)

#### About the `Hasher` interface

The [`Hasher` interface](https://github.com/ZalgoNoise/meta/blob/master/crypto/hash/hash.go) is a simple, generic and abstract interface to join all hashing algorithms:

```
package hash

// Hasher interface defines the behavior of a recursive hasher
// which takes in an empty interface (to support multiple formats)
// and returns a slice of bytes (the hashed seed)
type Hasher interface {
	Hash(data []byte) []byte
}
```

From this point, all hash functions implement the `Hash` method, which must work recursively (as described in the interface), in this case by passing a fixed-size slice of bytes to one without defined capacity (just like input data). Taking the `SHA256` struct as an example:

```
package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256 struct is a general placeholder for the SHA256
// algorithm, to implement the Hasher interface
type SHA256 struct{}

// Hash method satisfies the Hasher interface. It's a
// recursive hashing function to allow continuous hashing
func (hasher SHA256) Hash(data []byte) []byte {

	var hash []byte = make([]byte, hex.EncodedLen(32))
	var sum [32]byte = sha256.Sum256(data)

	hex.Encode(hash, sum[:])

	return hash

}
```

This means that, if your project simply needs a recursive hasher, you might as well just import the `zalgonoise/meta/crypto/hash` package instead of the entire `HashClockService` module. With this interface and these types you can create your own logic using it, with much more granularity.

____________

#### About the `HashClockService` module

This package (in `zalgonoise/hashclock/app/clock`) will contain logic to apply the `Hasher` interface as a generator / verifier of hash timestamps (based on the number of iterations). The module is very simple, containing three main types with a service's _request_ / _response_ approach:

```
package clock

import (
	"errors"
	"strings"
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
	algorithm  string
}

// HashClockResponse struct defines the input configuration for
// a `HashClockService` response
type HashClockResponse struct {
	Seed       string `json:"seed,omitempty"`
	Algorithm  string `json:"algorithm,omitempty"`
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
```

To spawn a HashClockService instance, it's best to use the public function `NewService()`, which return a pointer to a `HashClockService` with default settings, and an initialized request object (note that the response object is only generated as the response is being built):

```
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

	// initialize default hasher
	c.hasher = HasherMap[3]

	return c
}
```



#### Configuring the service

First of all, it's important to define the hashing algorithm that the service will use. This is done via the `SetHasher` method, that takes in a string (reference to a supported hashing algorithm) and sets the `HashClockService`'s hasher to the appropriate one, via two enum maps:

```
var HasherMap = map[int]rhash.Hasher{
	0: rhash.MD5{},
	1: rhash.SHA1{},
	2: rhash.SHA224{},
	3: rhash.SHA256{},
	4: rhash.SHA384{},
	5: rhash.SHA512{},
	6: rhash.SHA512_224{},
	7: rhash.SHA512_256{},
}

var HasherMapVals = map[int]string{
	0: "MD5",
	1: "SHA1",
	2: "SHA224",
	3: "SHA256",
	4: "SHA384",
	5: "SHA512",
	6: "SHA512_224",
	7: "SHA512_256",
}

// (...)

func (s *HashClockService) SetHasher(input string) error {
	for idx := 0; idx < len(HasherMapVals); idx++ {
		if input == HasherMapVals[idx] || input == strings.ToLower(HasherMapVals[idx]) {
			err := s.setHasher(idx)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("invalid hasher reference")
}

```

Quick example of the service being initialized in a new project, with a sha512 hasher:

```
package main 

import (
    "github.com/ZalgoNoise/hashclock/clock"
)

func main() {
    // set hasher string reference to sha512
    hasherRef := "sha512"

    // initialize the HashClockService module
    sClock := clock.NewService()

    // set service's hasher as sha512
    if err := sClock.SetHasher(hasherRef); err != nil {
        panic(err)
    }
}

```

#### Using the methods

All `HashClockService` methods can be used freely from this point forward, as the service is initialized and with a defined hasher. Here is a complete reference to all (current) methods in the `HashClockService`:

Method | Description | Header
:-----:|:-----------:|:-------:
`Hash` | This method will parse the `HashClockService.request` object and build its `HashClockResponse.response` with the hash for the seed string | `func (c *HashClockService) Hash(seed string) (*HashClockResponse, error) {}`
`RecHash` | This method takes in a string to hash and the number of desired iterations, returning an execution of the `newRecHashResponse` method | `func (c *HashClockService) RecHash(seed string, iter int) (*HashClockResponse, error) {}`
`RecHashPrint` | This method takes in a string to hash, the number of desired iterations, and a breakpoint value; returning an execution of the `newRecHashResponse` method | `func (c *HashClockService) RecHashPrint(seed string, iter int, breakpoint int) (*HashClockResponse, error) {}`
`RecHashLoop` | This method will recursively hash the seed string, infinitely (or until the program is halted) while printing out its hashes. | `func (c *HashClockService) RecHashLoop(seed string, breakpoint int) error {}`
`RecHashTimeout` | This method will take in a seed string and a timeout value (in seconds), returning an execution of the `newRecHashTimeResponse` method | `func (c *HashClockService) RecHashTimeout(seed string, timeout int) (*HashClockResponse, error) {}`
`Verify` | This method will take in a seed string and a target hash, returning an execution of the `newVerifyResponse` method | `func (c *HashClockService) Verify(seed string, hash string) (*HashClockResponse, error) {}`
`VerifyTimeout` | This method will take in a seed string, a target hash and a timeout value returning an execution of the `newVerifyTimeoutResponse` method | `func (c *HashClockService) VerifyTimeout(seed, hash string, timeout int) (*HashClockResponse, error) {}`
`VerifyIndex` | This method will take in a seed string, a target hash and target number of iterations returning an execution of the `newVerifyIndexResponse` method | `func (c *HashClockService) VerifyIndex(seed string, hash string, iterations int) (*HashClockResponse, error) {}`

______________

### CI/CD

Hashclock has a Continuous Integration system setup with Github Actions, with the following workflows set up:
- [Bazel-CI](https://github.com/ZalgoNoise/hashclock/actions/workflows/bazel-ci.yaml) will ensure that the libraries / executable are working with build / test / execution trials:
  - Runs `buildifier` to check and fix Bazel files
  - Runs `gazelle` to check and fix build files (Golang)
  - Runs `bazel build //...`
  - Runs `bazel test --test_output=all --test_summary=detailed --cache_test_results=no //...`
  - Runs `bazel run //:hashclock -- -seed hashclock -iter 10 -log 1`
- [Dockerfile-CI](https://github.com/ZalgoNoise/hashclock/actions/workflows/docker-build.yaml) will ensure that the existing `Dockerfile` in the project is working (since Bazel builds Docker images differently).
  - Runs `docker build -f ./Dockerfile .`
- [Go](https://github.com/ZalgoNoise/hashclock/actions/workflows/go.yaml) will ensure that, without Bazel, the project works as intended (using the Golang compiler):
  - Runs `go build -v ./...`
  - Runs `go test -v ./...`


For Continuous Deployment, this project takes leverage of Bazel with specific targets to both build and push Docker images to `ghcr.io` and `docker.io`. For the moment, building and distributing new `hashclock` binaries in a new release is being done manually. A `genrule` (or even a `.bzl` rule) should be developed in the future to automate that part of the process.
- [Bazel-CD](https://github.com/ZalgoNoise/hashclock/actions/workflows/bazel-cd.yaml) will ensure that new pushes to the `master` branch will create new `latest` versions of the Docker containers. For each container registry (Github and DockerHub), the following actions take place:
  - Launch `bazelisk` with mounted cache.
  - Build container with `bazel build //:latest`
  - Push new image to registry with `bazel run //:github-push` or `bazel run //:dockerhub-push`


__________

### Contributing

All contributions are welcome. Feel free to file a new PR. 