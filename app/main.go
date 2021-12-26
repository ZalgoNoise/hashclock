// Hashclock is a cryptographic timestamp generator; or a tool
// to benchmark CPU speed (with number of hashes per second)
//
// Input strings will be hashed with a SHA256 algorithm
// and the resulting slice of (32) bytes is hex-encoded
//
// Each hash value can be input into the same hashing function,
// producing a "hash of a hash".
//
// This is done in succession for a determined number of loops,
// certain amount of time, or indefinitely.
//
// The number of iterations combined with the hash and the seed
// can be used to verify whether the claim is true, by calculating
// the hash of the seed the same amount of times as the number of
// iterations in the claim.
//
package main

import (
	"github.com/ZalgoNoise/hashclock/cmd"
)

func main() {
	cmd.Run()
}
