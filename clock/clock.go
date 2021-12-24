package clock

import (
	rsha "github.com/ZalgoNoise/meta/crypto/hash"

	"errors"
	"fmt"
)


func RecursiveSHA256(input string, iter int) ([]byte, error) {
	var hash []byte


	if iter <= 0 {
		return []byte{}, errors.New("number of iterations has to be greater than zero")
	}

	for i := 1 ; i <= iter ; i++ {
		if i == 1 {
			hash = rsha.Hash(input)
		} else {
			hash = rsha.Hash(hash)
		}	
	}

	return hash, nil
}


func RecursiveSHA256Logged(input string, iter int, breakpoint int) ([]byte, error) {
	var hash []byte

	if iter <= 0 {
		return []byte{}, errors.New("number of iterations has to be greater than zero")
	}

	for i := 1 ; i <= iter ; i++ {
		if i == 1 {
			hash = rsha.Hash(input)
		} else {
			hash = rsha.Hash(hash)

		}
		if breakpoint > 0 && i % breakpoint == 0 {
			fmt.Printf("#%v:\t\t%x\n", i, hash)
		}
		
	}
	return hash, nil
}

func RecursiveSHA256Inf(input interface{}, breakpoint int) {
	hash := rsha.Hash(input)
	var counter int = 1
	for {
		hash = rsha.Hash(hash)
		counter++
		if breakpoint > 0 && counter % breakpoint == 0 {
			fmt.Printf("#%v:\t\t%x\n", counter, hash)
		}
	}
}