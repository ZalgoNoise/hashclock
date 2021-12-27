package clock

import (
	"testing"
)

type testCase struct {
	seed       string
	iterations int
	hash       string
	breakpoint int
	timeout    int
}

var testCases []testCase = []testCase{
	{
		seed:       "Hello World!",
		iterations: 1,
		hash:       "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
	}, {
		seed:       "Hello, 世界!",
		iterations: 1,
		hash:       "7de2f06498b5b4d53b170000c311101b55046a3c889efd54351cb3697fcf57cc",
	}, {
		seed:       "Hello World!",
		iterations: 2,
		hash:       "4163fb4ab9e1e0a51709a51bc7e13ab6792907905960145c722d2c1479caac42",
	}, {
		seed:       "Hello, 世界!",
		iterations: 2,
		hash:       "8888cd035a543f3c28b335c2afe8a09fc49530f37380271a5b1ba1925b5f355a",
	}, {
		seed:       "Hello World!",
		iterations: 10,
		hash:       "1fada6a9b8084eff6347baeae0812aac1c14fcc0a97759a6bfa2d8a2ea087705",
	}, {
		seed:       "Hello, 世界!",
		iterations: 10,
		hash:       "eb98b6704c17b6add5642ea05e12e5ad2d3c94c39b9950dacfb9bea8882b9a98",
	}, {
		seed:       "Hello World!",
		iterations: 10000,
		hash:       "aee751f97c1fc4f9df69a52b02a60e40a1a05c872a879cfef077577419a7814b",
	}, {
		seed:       "Hello, 世界!",
		iterations: 10000,
		hash:       "a2ee566ab660737832f3081753057005f1d559258480f630b02e194b639ce22a",
	},
}

var testBreakages []testCase = []testCase{
	// empty seed / empty hash case
	{
		seed:       "",
		iterations: 1,
		hash:       "",
	},

	// zero iterations case
	{
		seed:       "Hello World!",
		iterations: 0,
		hash:       "",
	},

	// zero breakpoint case
	{
		seed:       "Hello World!",
		iterations: 1,
		hash:       "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
		breakpoint: 0,
	},

	// negative breakpoint case
	{
		seed:       "Hello World!",
		iterations: 1,
		hash:       "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
		breakpoint: -2,
	},

	// zero timeout case
	{
		seed:       "Hello World!",
		iterations: 1,
		hash:       "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
		timeout:    0,
	},

	// invalid hash string (hex encoding) case
	{
		seed:       "Hello World!",
		iterations: 1,
		hash:       "zzz3b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
		timeout:    0,
	},

	// string == hash case
	{
		seed:       "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
		iterations: 1,
		hash:       "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
		timeout:    0,
	},
}

func TestHash(t *testing.T) {

	tests := []struct {
		input string
		ok    string
		pass  bool
	}{
		{
			input: testCases[0].seed,
			ok:    testCases[0].hash,
			pass:  true,
		}, {
			input: testCases[1].seed,
			ok:    testCases[1].hash,
			pass:  true,
		}, {
			input: testBreakages[0].seed,
			ok:    testBreakages[0].hash,
			pass:  false,
		},
	}

	clock := NewService()

	for id, test := range tests {
		result, err := clock.Hash(test.input)
		if test.pass && err != nil {
			t.Errorf(
				"#%v [HashClockService] Hash(%s) resulted in error: %s ; expected %s",
				id,
				test.input,
				err,
				test.ok,
			)
		}

		if test.pass {
			if test.ok != result.Hash {
				t.Errorf(
					"#%v [HashClockService] Hash(%s) = %s ; expected %s",
					id,
					test.input,
					result.Hash,
					test.ok,
				)
			}

			t.Logf(
				"#%v -- TESTED -- [HashClockService] Hash(%s) = %s",
				id,
				test.input,
				test.ok,
			)
		}
	}
}

func TestRecHash(t *testing.T) {

	tests := []struct {
		input      string
		iterations int
		ok         string
		pass       bool
	}{
		{
			input:      testCases[2].seed,
			iterations: testCases[2].iterations,
			ok:         testCases[2].hash,
			pass:       true,
		},
		{
			input:      testCases[3].seed,
			iterations: testCases[3].iterations,
			ok:         testCases[3].hash,
			pass:       true,
		}, {
			input:      testCases[4].seed,
			iterations: testCases[4].iterations,
			ok:         testCases[4].hash,
			pass:       true,
		},
		{
			input:      testCases[5].seed,
			iterations: testCases[5].iterations,
			ok:         testCases[5].hash,
			pass:       true,
		}, {
			input:      testCases[6].seed,
			iterations: testCases[6].iterations,
			ok:         testCases[6].hash,
			pass:       true,
		},
		{
			input:      testCases[7].seed,
			iterations: testCases[7].iterations,
			ok:         testCases[7].hash,
			pass:       true,
		}, {
			input:      testBreakages[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testBreakages[0].hash,
			pass:       false,
		}, {
			input:      testBreakages[1].seed,
			iterations: testBreakages[1].iterations,
			ok:         testBreakages[1].hash,
			pass:       false,
		},
	}

	clock := NewService()

	for id, test := range tests {
		result, err := clock.RecHash(test.input, test.iterations)
		if test.pass && err != nil {
			t.Errorf(
				"#%v [HashClockService] RecHash(%s, %v) resulted in error: %s ; expected %s",
				id,
				test.input,
				test.iterations,
				err,
				test.ok,
			)
		}

		if test.pass {
			if test.ok != result.Hash {
				t.Errorf(
					"#%v [HashClockService] RecHash(%s, %v) = %s ; expected %s",
					id,
					test.input,
					test.iterations,
					result.Hash,
					test.ok,
				)
			}

			t.Logf(
				"#%v -- TESTED -- [HashClockService] RecHash(%s, %v) = %s",
				id,
				test.input,
				test.iterations,
				test.ok,
			)
		}
	}
}

func TestRecHashPrint(t *testing.T) {
	tests := []struct {
		input      string
		iterations int
		breakpoint int
		ok         string
		pass       bool
	}{
		{
			input:      testCases[2].seed,
			iterations: testCases[2].iterations,
			breakpoint: 1,
			ok:         testCases[2].hash,
			pass:       true,
		},
		{
			input:      testCases[3].seed,
			iterations: testCases[3].iterations,
			breakpoint: 1,
			ok:         testCases[3].hash,
			pass:       true,
		}, {
			input:      testCases[4].seed,
			iterations: testCases[4].iterations,
			breakpoint: 5,
			ok:         testCases[4].hash,
			pass:       true,
		},
		{
			input:      testCases[5].seed,
			iterations: testCases[5].iterations,
			breakpoint: 5,
			ok:         testCases[5].hash,
			pass:       true,
		}, {
			input:      testBreakages[0].seed,
			iterations: testBreakages[0].iterations,
			ok:         testBreakages[0].hash,
			pass:       false,
		}, {
			input:      testBreakages[1].seed,
			iterations: testBreakages[1].iterations,
			ok:         testBreakages[1].hash,
			pass:       false,
		}, {
			input:      testBreakages[2].seed,
			iterations: testBreakages[2].iterations,
			ok:         testBreakages[2].hash,
			pass:       false,
		}, {
			input:      testBreakages[3].seed,
			iterations: testBreakages[3].iterations,
			ok:         testBreakages[3].hash,
			pass:       false,
		},
	}

	clock := NewService()

	for id, test := range tests {
		result, err := clock.RecHashPrint(test.input, test.iterations, test.breakpoint)
		if test.pass && err != nil {
			t.Errorf(
				"#%v [HashClockService] RecHashPrint(%s, %v, %v) resulted in error: %s ; expected %s",
				id,
				test.input,
				test.iterations,
				test.breakpoint,
				err,
				test.ok,
			)
		}

		if test.pass {
			if test.ok != result.Hash {
				t.Errorf(
					"#%v [HashClockService] RecHashPrint(%s, %v, %v) = %s ; expected %s",
					id,
					test.input,
					test.iterations,
					test.breakpoint,
					result.Hash,
					test.ok,
				)
			}
			t.Logf(
				"#%v -- TESTED -- [HashClockService] RecHashPrint(%s, %v, %v) = %s",
				id,
				test.input,
				test.iterations,
				test.breakpoint,
				test.ok,
			)
		}
	}
}

// This is an infinite-loop function; all tests must be breakage / error tests
func TestRecHashLoop(t *testing.T) {
	tests := []struct {
		input      string
		breakpoint int
		pass       bool
	}{
		{
			input:      testBreakages[0].seed,
			breakpoint: 1,
			pass:       false,
		}, {
			input:      testBreakages[2].seed,
			breakpoint: testBreakages[2].breakpoint,
			pass:       false,
		}, {
			input:      testBreakages[3].seed,
			breakpoint: testBreakages[3].breakpoint,
			pass:       false,
		},
	}

	clock := NewService()

	for id, test := range tests {
		err := clock.RecHashLoop(test.input, test.breakpoint)
		if test.pass && err != nil {
			t.Errorf(
				"%v [HashClockService] RecHashLoop(%s, %v) resulted in an unexpected error: %s",
				id,
				test.input,
				test.breakpoint,
				err,
			)
		}
		t.Logf(
			"#%v -- TESTED -- [HashClockService] RecHashLoop(%s, %v) = %s",
			id,
			test.input,
			test.breakpoint,
			err,
		)
	}
}

func TestRecHashTimeout(t *testing.T) {

	tests := []struct {
		input   string
		timeout int
		pass    bool
	}{
		{
			input:   testCases[0].seed,
			timeout: 1,
			pass:    true,
		}, {
			input:   testCases[1].seed,
			timeout: 1,
			pass:    true,
		}, {
			input:   testBreakages[0].seed,
			timeout: 10,
			pass:    false,
		}, {
			input:   testBreakages[4].seed,
			timeout: testBreakages[4].timeout,
			pass:    false,
		},
	}

	clock := NewService()

	for id, test := range tests {
		result, err := clock.RecHashTimeout(test.input, test.timeout)
		if test.pass && err != nil {
			t.Errorf(
				"#%v [HashClockService] RecHashTimeout(%s, %v) resulted in an unexpected error: %s",
				id,
				test.input,
				test.timeout,
				err,
			)
		}

		if test.pass {
			// verify calculated hash
			verify, err := clock.VerifyIndex(test.input, result.Hash, result.Iterations)
			if err != nil {
				t.Errorf(
					"#%v [HashClockService] VerifyIndex(%s, %s, %v) resulted in an unexpected error: %s",
					id,
					test.input,
					result.Hash,
					result.Iterations,
					err,
				)
			}

			if result.Hash != verify.Hash {
				t.Errorf(
					"#%v [HashClockService] RecHashTimeout(%s, %v) calculated an invalid hash: %s with %v iterations. Verification expected hash %s",
					id,
					test.input,
					test.timeout,
					result.Hash,
					result.Iterations,
					verify.Hash,
				)
			}

			t.Logf(
				"#%v -- TESTED -- [HashClockService] RecHashTimeout(%s, %v) = %s with %v iterations",
				id,
				test.input,
				test.timeout,
				result.Hash,
				result.Iterations,
			)
		}
	}
}
