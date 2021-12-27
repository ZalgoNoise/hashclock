package clock

import (
	"math/rand"
	"testing"

	rhash "github.com/ZalgoNoise/meta/crypto/hash"
)

type verificationCase struct {
	input      string
	iterations int
	timeout    int
	ok         string
	pass       bool
}

func TestVerify(t *testing.T) {

	var tests []verificationCase

	// run 500 random iteration blind tests
	for i := 0; i < 500; i++ {
		vCase := verificationCase{}
		if i == 0 || i%2 == 0 {
			vCase.input = testCases[0].seed
		} else {
			vCase.input = testCases[1].seed
		}
		vCase.iterations = rand.Intn(1000) + 1
		vCase.pass = true
		tests = append(tests, vCase)
	}

	// add breakage / failing tests
	tests = append(tests,
		verificationCase{
			input:      testBreakages[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testCases[0].hash,
			pass:       false,
		},
		verificationCase{
			input:      testCases[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testBreakages[0].hash,
			pass:       false,
		},
		verificationCase{
			input:      testCases[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testBreakages[5].hash,
			pass:       false,
		},
		verificationCase{
			input:      testCases[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testBreakages[6].hash,
			pass:       false,
		},
	)

	clock := NewService()

	for id, test := range tests {
		// get hash for this test
		calc, err := clock.RecHash(test.input, test.iterations)
		if test.pass && err != nil {
			t.Errorf(
				"#%v [HashClockService] RecHash(%s, %v) resulted in error: %s",
				id,
				test.input,
				test.iterations,
				err,
			)
		}

		// filter (purposefully) breakage / failing tests
		if test.pass {
			test.ok = calc.Hash

			result, err := clock.Verify(test.input, test.ok)
			if err != nil {
				t.Errorf(
					"#%v [HashClockService] Verify(%s, %s) resulted in error: %s",
					id,
					test.input,
					test.ok,
					err,
				)
			}

			if !result.Match {
				t.Errorf(
					"#%v [HashClockService] Verify(%s, %s) = %v ; expected true",
					id,
					test.input,
					test.ok,
					result.Match,
				)
			}
			if result.Hash != test.ok {
				t.Errorf(
					"#%v [HashClockService] Verify(%s, %s) = %v ; string comparison failed: %s != %s",
					id,
					test.input,
					test.ok,
					result.Match,
					result.Hash,
					test.ok,
				)
			}
			if result.Iterations != test.iterations {
				t.Errorf(
					"#%v [HashClockService] Verify(%s, %s) = %v ; iteration count comparison failed: %v != %v",
					id,
					test.input,
					test.ok,
					result.Match,
					result.Iterations,
					test.iterations,
				)
			}
			t.Logf(
				"#%v -- TESTED -- [HashClockService] Verify(%s, %s)",
				id,
				test.input,
				test.ok,
			)
		}
	}

}

func TestVerifyTimeout(t *testing.T) {

	var tests []verificationCase

	// run 500 random iteration blind tests
	for i := 0; i < 500; i++ {
		vCase := verificationCase{}
		if i == 0 || i%2 == 0 {
			vCase.input = testCases[0].seed
		} else {
			vCase.input = testCases[1].seed
		}
		vCase.iterations = rand.Intn(1000) + 1
		vCase.timeout = 3
		vCase.pass = true
		tests = append(tests, vCase)
	}

	// add breakage / failing tests
	tests = append(tests,
		verificationCase{
			input:      testBreakages[0].seed,
			iterations: testBreakages[0].iterations,
			timeout:    3,
			ok:         testCases[0].hash,
			pass:       false,
		},
		verificationCase{
			input:      testCases[0].seed,
			iterations: testBreakages[0].iterations,
			timeout:    3,
			ok:         testBreakages[0].hash,
			pass:       false,
		},
		verificationCase{
			input:      testCases[0].seed,
			iterations: testBreakages[0].iterations,
			timeout:    3,
			ok:         testBreakages[5].hash,
			pass:       false,
		},
		verificationCase{
			input:      testCases[0].seed,
			iterations: testBreakages[0].iterations,
			timeout:    3,
			ok:         testBreakages[6].hash,
			pass:       false,
		},
		verificationCase{
			input:      testBreakages[4].seed,
			iterations: testBreakages[4].iterations,
			timeout:    testBreakages[4].timeout,
			ok:         testBreakages[4].hash,
			pass:       false,
		},
	)

	clock := NewService()

	for id, test := range tests {
		// get hash for this test
		calc, err := clock.RecHash(test.input, test.iterations)
		if test.pass && err != nil {
			t.Errorf(
				"#%v [HashClockService] RecHash(%s, %v) resulted in error: %s",
				id,
				test.input,
				test.iterations,
				err,
			)
		}

		// filter (purposefully) breakage / failing tests
		if test.pass {
			test.ok = calc.Hash

			result, err := clock.VerifyTimeout(test.input, test.ok, test.timeout)
			if err != nil {
				t.Errorf(
					"#%v [HashClockService] VerifyTimeout(%s, %s, %v) resulted in error: %s",
					id,
					test.input,
					test.ok,
					test.timeout,
					err,
				)
			}

			if !result.Match {
				t.Errorf(
					"#%v [HashClockService] VerifyTimeout(%s, %s, %v) = %v ; expected true",
					id,
					test.input,
					test.ok,
					test.timeout,
					result.Match,
				)
			}
			if result.Hash != test.ok {
				t.Errorf(
					"#%v [HashClockService] VerifyTimeout(%s, %s, %v) = %v ; string comparison failed: %s != %s",
					id,
					test.input,
					test.ok,
					test.timeout,
					result.Match,
					result.Hash,
					test.ok,
				)
			}
			if result.Iterations != test.iterations {
				t.Errorf(
					"#%v [HashClockService] VerifyTimeout(%s, %s, %v) = %v ; iteration count comparison failed: %v != %v",
					id,
					test.input,
					test.ok,
					test.timeout,
					result.Match,
					result.Iterations,
					test.iterations,
				)
			}

			t.Logf(
				"#%v -- TESTED -- [HashClockService] VerifyTimeout(%s, %s, %v)",
				id,
				test.input,
				test.ok,
				test.timeout,
			)
		}
	}

}

func TestVerifyIndex(t *testing.T) {

	var tests []verificationCase

	// run 500 random iteration blind tests
	for i := 0; i < 500; i++ {
		vCase := verificationCase{}
		if i == 0 || i%2 == 0 {
			vCase.input = testCases[0].seed
		} else {
			vCase.input = testCases[1].seed
		}
		vCase.iterations = rand.Intn(1000) + 1
		vCase.pass = true
		tests = append(tests, vCase)
	}

	// add breakage / failing tests
	tests = append(tests,
		verificationCase{
			input:      testBreakages[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testCases[0].hash,
			pass:       false,
		},
		verificationCase{
			input:      testCases[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testBreakages[0].hash,
			pass:       false,
		},
		verificationCase{
			input:      testCases[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testBreakages[5].hash,
			pass:       false,
		},
		verificationCase{
			input:      testCases[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testBreakages[6].hash,
			pass:       false,
		},
		verificationCase{
			input:      testBreakages[1].seed,
			iterations: testBreakages[1].iterations,
			ok:         testCases[0].hash,
			pass:       false,
		},
	)

	clock := NewService()

	for id, test := range tests {
		// get hash for this test
		calc, err := clock.RecHash(test.input, test.iterations)
		if test.pass && err != nil {
			t.Errorf(
				"#%v [HashClockService] RecHash(%s, %v) resulted in error: %s",
				id,
				test.input,
				test.iterations,
				err,
			)
		}

		// filter (purposefully) breakage / failing tests
		if test.pass {
			test.ok = calc.Hash

			result, err := clock.VerifyIndex(test.input, test.ok, test.iterations)
			if err != nil {
				t.Errorf(
					"#%v [HashClockService] VerifyIndex(%s, %s, %v) resulted in error: %s",
					id,
					test.input,
					test.ok,
					test.iterations,
					err,
				)
			}

			if !result.Match {
				t.Errorf(
					"#%v [HashClockService] VerifyIndex(%s, %s, %v) = %v ; expected true",
					id,
					test.input,
					test.ok,
					test.iterations,
					result.Match,
				)
			}
			if result.Hash != test.ok {
				t.Errorf(
					"#%v [HashClockService] VerifyIndex(%s, %s, %v) = %v ; string comparison failed: %s != %s",
					id,
					test.input,
					test.ok,
					test.iterations,
					result.Match,
					result.Hash,
					test.ok,
				)
			}
			if result.Iterations != test.iterations {
				t.Errorf(
					"#%v [HashClockService] VerifyIndex(%s, %s, %v) = %v ; iteration count comparison failed: %v != %v",
					id,
					test.input,
					test.ok,
					test.iterations,
					result.Match,
					result.Iterations,
					test.iterations,
				)
			}

			t.Logf(
				"#%v -- TESTED -- [HashClockService] VerifyIndex(%s, %s, %v)",
				id,
				test.input,
				test.ok,
				test.iterations,
			)
		}
	}

}

func TestMatchHash(t *testing.T) {
	tests := []struct {
		input      string
		iterations int
		ok         string
	}{
		{
			input:      testCases[0].seed,
			iterations: testCases[0].iterations,
			ok:         testCases[0].hash,
		}, {
			input:      testCases[1].seed,
			iterations: testCases[1].iterations,
			ok:         testCases[1].hash,
		}, {
			input:      testCases[2].seed,
			iterations: testCases[2].iterations,
			ok:         testCases[2].hash,
		}, {
			input:      testCases[3].seed,
			iterations: testCases[3].iterations,
			ok:         testCases[3].hash,
		}, {
			input:      testCases[4].seed,
			iterations: testCases[4].iterations,
			ok:         testCases[4].hash,
		}, {
			input:      testCases[5].seed,
			iterations: testCases[5].iterations,
			ok:         testCases[5].hash,
		}, {
			input:      testCases[6].seed,
			iterations: testCases[6].iterations,
			ok:         testCases[6].hash,
		}, {
			input:      testCases[7].seed,
			iterations: testCases[7].iterations,
			ok:         testCases[7].hash,
		},
	}

	h := rhash.SHA256{}

	for id, test := range tests {
		hash := h.Hash([]byte(test.input))
		target := []byte(test.ok)

		// index 0 is the seed
		// index 1 is the first hash, above
		// rehashing index starts at 2
		for i := 2; i <= test.iterations; i++ {
			hash = h.Hash(hash)
		}

		if !matchHash(hash, target) {
			t.Errorf(
				"#%v [HashClockService] matchHash(%s, %s) = false ; expected true",
				id,
				string(hash),
				test.ok,
			)
		}

		t.Logf(
			"#%v -- TESTED -- [HashClockService] matchHash(%s, %s) = %v",
			id,
			string(hash),
			test.ok,
			matchHash(hash, target),
		)

	}
}
