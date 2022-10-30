package entities

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseEnvRuntimesWithValidRuntimes(t *testing.T) {
	testCases := []struct {
		test             string
		runtimes         []string
		expectedRuntimes EnvRuntimes
	}{
		{
			test:     "with mixed runtimes",
			runtimes: []string{"docker@latest", "php@8.1", "node", "rust@1.61.0", "python@", "go@1.19.0-alpha"},
			expectedRuntimes: EnvRuntimes{
				"docker": "latest",
				"php":    "8.1",
				"node":   "latest",
				"rust":   "1.61.0",
				"python": "latest",
				"go":     "1.19.0-alpha",
			},
		},

		{
			test:     "without passed version",
			runtimes: []string{"docker", "java", "clang", "ruby", "go", "python", "php", "node"},
			expectedRuntimes: EnvRuntimes{
				"docker": "latest",
				"java":   "latest",
				"clang":  "latest",
				"ruby":   "latest",
				"go":     "latest",
				"python": "latest",
				"php":    "latest",
				"node":   "latest",
			},
		},

		{
			test:     "with passed latest version",
			runtimes: []string{"docker@latest", "java@latest", "clang@latest", "go@latest"},
			expectedRuntimes: EnvRuntimes{
				"docker": "latest",
				"java":   "latest",
				"clang":  "latest",
				"go":     "latest",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			parsedRuntimes, err := ParseEnvRuntimes(
				tc.runtimes,
			)

			if err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}

			if !reflect.DeepEqual(parsedRuntimes, tc.expectedRuntimes) {
				t.Fatalf(
					"expected runtimes to equal '%+v', got '%+v'",
					tc.expectedRuntimes,
					parsedRuntimes,
				)
			}
		})
	}
}

func TestParseEnvRuntimesWithInvalidRuntimes(t *testing.T) {
	testCases := []struct {
		test                   string
		runtimes               []string
		expectedInvalidRuntime string
	}{
		{
			test:                   "with valid and invalid runtime",
			runtimes:               []string{"docker@latest", "php@8.1", "bibi"},
			expectedInvalidRuntime: "bibi",
		},

		{
			test:                   "with only invalid runtime",
			runtimes:               []string{"bibi", "sisi"},
			expectedInvalidRuntime: "bibi",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := ParseEnvRuntimes(
				tc.runtimes,
			)

			if err == nil {
				t.Fatalf("expected error, got nothing")
			}

			if !errors.As(err, &ErrEnvInvalidRuntime{}) {
				t.Fatalf(
					"expected error to equal '%+v', got '%+v'",
					ErrEnvInvalidRuntime{},
					err,
				)
			}

			typedError := err.(ErrEnvInvalidRuntime)

			if typedError.Runtime != tc.expectedInvalidRuntime {
				t.Fatalf(
					"expected invalid runtime to equal '%+v', got '%+v'",
					tc.expectedInvalidRuntime,
					typedError.Runtime,
				)
			}
		})
	}
}

func TestParseEnvRuntimesWithInvalidRuntimeVersions(t *testing.T) {
	testCases := []struct {
		test                   string
		runtimes               []string
		expectedInvalidRuntime string
		expectedInvalidVersion string
	}{
		{
			test:                   "with invalid PHP version",
			runtimes:               []string{"go@latest", "node", "php@8.1.3", "ruby@3.1.0"},
			expectedInvalidRuntime: "php",
			expectedInvalidVersion: "8.1.3",
		},

		{
			test:                   "with invalid Docker version",
			runtimes:               []string{"go@latest", "node", "php@8.1", "ruby@3.1.0", "docker@4.3.0"},
			expectedInvalidRuntime: "docker",
			expectedInvalidVersion: "4.3.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := ParseEnvRuntimes(
				tc.runtimes,
			)

			if err == nil {
				t.Fatalf("expected error, got nothing")
			}

			if !errors.As(err, &ErrEnvInvalidRuntimeVersion{}) {
				t.Fatalf(
					"expected error to equal '%+v', got '%+v'",
					ErrEnvInvalidRuntimeVersion{},
					err,
				)
			}

			typedError := err.(ErrEnvInvalidRuntimeVersion)

			if typedError.Runtime != tc.expectedInvalidRuntime {
				t.Fatalf(
					"expected invalid runtime to equal '%+v', got '%+v'",
					tc.expectedInvalidRuntime,
					typedError.Runtime,
				)
			}

			if typedError.RuntimeVersion != tc.expectedInvalidVersion {
				t.Fatalf(
					"expected invalid runtime version to equal '%+v', got '%+v'",
					tc.expectedInvalidVersion,
					typedError.RuntimeVersion,
				)
			}
		})
	}
}
