package entities

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestEnvDoesServedPortExist(t *testing.T) {
	env := NewEnv(
		"test_env",
		0,
		"test_instance_type",
		[]EnvRepository{},
		EnvRuntimes{},
	)

	servedPortToSearch := EnvServedPort("8000")

	if exists := env.DoesServedPortExist(servedPortToSearch); exists {
		t.Fatalf("expected served port to not exist")
	}

	env.ServedPorts[servedPortToSearch] = []EnvServedPortBinding{
		{
			Type:  EnvServedPortBindingTypeDomain,
			Value: "api.eleven.sh",
		},
	}

	if exists := env.DoesServedPortExist(servedPortToSearch); !exists {
		t.Fatalf("expected served port to exist")
	}
}

func TestEnvAddServedPortBinding(t *testing.T) {

	type bindings struct {
		values          []string
		redirectToHTTPS bool
	}

	testCases := []struct {
		test                string
		bindings            map[string]bindings
		expectedServedPorts EnvServedPorts
	}{
		{
			test: "with mixed port bindings",

			bindings: map[string]bindings{
				"8000": bindings{
					values: []string{
						"8000",
						"api.eleven.sh",
					},
					redirectToHTTPS: true,
				},

				"6000": bindings{
					values: []string{
						"6000",
						"4000",
					},
					redirectToHTTPS: false,
				},

				"11000": bindings{
					values: []string{
						"test.eleven.sh",
					},
					redirectToHTTPS: true,
				},

				"12000": bindings{
					values: []string{
						"12000",
					},
					redirectToHTTPS: false,
				},

				"14000": bindings{
					values: []string{
						"a.eleven.sh",
						"b.eleven.sh",
					},
					redirectToHTTPS: true,
				},
			},

			expectedServedPorts: EnvServedPorts{
				"8000": {
					{
						Value:           "8000",
						Type:            EnvServedPortBindingTypePort,
						RedirectToHTTPS: true,
					},

					{
						Value:           "api.eleven.sh",
						Type:            EnvServedPortBindingTypeDomain,
						RedirectToHTTPS: true,
					},
				},

				"6000": {
					{
						Value:           "6000",
						Type:            EnvServedPortBindingTypePort,
						RedirectToHTTPS: false,
					},

					{
						Value:           "4000",
						Type:            EnvServedPortBindingTypePort,
						RedirectToHTTPS: false,
					},
				},

				"11000": {
					{
						Value:           "test.eleven.sh",
						Type:            EnvServedPortBindingTypeDomain,
						RedirectToHTTPS: true,
					},
				},

				"12000": {
					{
						Value:           "12000",
						Type:            EnvServedPortBindingTypePort,
						RedirectToHTTPS: false,
					},
				},

				"14000": {
					{
						Value:           "a.eleven.sh",
						Type:            EnvServedPortBindingTypeDomain,
						RedirectToHTTPS: true,
					},

					{
						Value:           "b.eleven.sh",
						Type:            EnvServedPortBindingTypeDomain,
						RedirectToHTTPS: true,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			env := NewEnv(
				"test_env",
				0,
				"test_instance_type",
				[]EnvRepository{},
				EnvRuntimes{},
			)

			for servedPort, bindings := range tc.bindings {
				for _, binding := range bindings.values {
					env.AddServedPortBinding(
						EnvServedPort(servedPort),
						binding,
						bindings.redirectToHTTPS,
					)
				}
			}

			if !reflect.DeepEqual(env.ServedPorts, tc.expectedServedPorts) {
				t.Fatalf(
					"expected served ports to equal '%+v', got '%+v'",
					tc.expectedServedPorts,
					env.ServedPorts,
				)
			}
		})
	}
}

func TestEnvDoesServedPortBindingExist(t *testing.T) {
	testCases := []struct {
		test              string
		bindings          map[string][]string
		wantedBinding     string
		expectedExistence bool
	}{
		{
			test: "with port binding",
			bindings: map[string][]string{
				"8000": {
					"8000",
					"api.eleven.sh",
				},
			},
			wantedBinding:     "8000",
			expectedExistence: true,
		},

		{
			test: "with domain binding",
			bindings: map[string][]string{
				"8000": {
					"8000",
				},

				"6000": {
					"api.eleven.sh",
				},
			},
			wantedBinding:     "api.eleven.sh",
			expectedExistence: true,
		},

		{
			test: "without binding",
			bindings: map[string][]string{
				"8000": {
					"8000",
					"api.eleven.sh",
				},
			},
			wantedBinding:     "6000",
			expectedExistence: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			env := NewEnv(
				"test_env",
				0,
				"test_instance_type",
				[]EnvRepository{},
				EnvRuntimes{},
			)

			for servedPort, bindings := range tc.bindings {
				for _, binding := range bindings {
					env.AddServedPortBinding(
						EnvServedPort(servedPort),
						binding,
						false,
					)
				}
			}

			exists := env.DoesServedPortBindingExist(tc.wantedBinding)

			if exists != tc.expectedExistence {
				if tc.expectedExistence {
					t.Fatalf("expected port binding to exist")
				}

				t.Fatalf("expected port binding to not exist")
			}
		})
	}
}

func TestEnvRemoveServedPortBinding(t *testing.T) {
	testCases := []struct {
		test                string
		bindings            map[string][]string
		bindingsToRemove    map[string][]string
		expectedServedPorts EnvServedPorts
	}{
		{
			test: "remove existing port binding",

			bindings: map[string][]string{
				"8000": {
					"8000",
					"api.eleven.sh",
				},
			},

			bindingsToRemove: map[string][]string{
				"8000": {
					"8000",
				},
			},

			expectedServedPorts: EnvServedPorts{
				"8000": {
					{
						Value: "api.eleven.sh",
						Type:  EnvServedPortBindingTypeDomain,
					},
				},
			},
		},

		{
			test: "remove existing domain binding",

			bindings: map[string][]string{
				"8000": {
					"8000",
					"api.eleven.sh",
				},
			},

			bindingsToRemove: map[string][]string{
				"8000": {
					"api.eleven.sh",
				},
			},

			expectedServedPorts: EnvServedPorts{
				"8000": {
					{
						Value: "8000",
						Type:  EnvServedPortBindingTypePort,
					},
				},
			},
		},

		{
			test: "remove unexisting port binding",

			bindings: map[string][]string{
				"8000": {
					"8000",
					"api.eleven.sh",
				},
			},

			bindingsToRemove: map[string][]string{
				"8000": {
					"6000",
				},
			},

			expectedServedPorts: EnvServedPorts{
				"8000": {
					{
						Value: "8000",
						Type:  EnvServedPortBindingTypePort,
					},

					{
						Value: "api.eleven.sh",
						Type:  EnvServedPortBindingTypeDomain,
					},
				},
			},
		},

		{
			test: "remove unexisting domain binding",

			bindings: map[string][]string{
				"8000": {
					"8000",
					"api.eleven.sh",
				},
			},

			bindingsToRemove: map[string][]string{
				"8000": {
					"test.eleven.sh",
				},
			},

			expectedServedPorts: EnvServedPorts{
				"8000": {
					{
						Value: "8000",
						Type:  EnvServedPortBindingTypePort,
					},

					{
						Value: "api.eleven.sh",
						Type:  EnvServedPortBindingTypeDomain,
					},
				},
			},
		},

		{
			test: "remove unexisting served port",

			bindings: map[string][]string{
				"8000": {
					"8000",
					"api.eleven.sh",
				},
			},

			bindingsToRemove: map[string][]string{
				"6000": {
					"api.eleven.sh",
				},
			},

			expectedServedPorts: EnvServedPorts{
				"8000": {
					{
						Value: "8000",
						Type:  EnvServedPortBindingTypePort,
					},

					{
						Value: "api.eleven.sh",
						Type:  EnvServedPortBindingTypeDomain,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			env := NewEnv(
				"test_env",
				0,
				"test_instance_type",
				[]EnvRepository{},
				EnvRuntimes{},
			)

			for servedPort, bindings := range tc.bindings {
				for _, binding := range bindings {
					env.AddServedPortBinding(
						EnvServedPort(servedPort),
						binding,
						false,
					)
				}
			}

			for servedPort, bindings := range tc.bindingsToRemove {
				for _, binding := range bindings {
					env.RemoveServedPortBinding(
						EnvServedPort(servedPort),
						binding,
					)
				}
			}

			if !reflect.DeepEqual(env.ServedPorts, tc.expectedServedPorts) {
				t.Fatalf(
					"expected served ports to equal '%+v', got '%+v'",
					tc.expectedServedPorts,
					env.ServedPorts,
				)
			}
		})
	}
}

func TestEnvRemoveServedPort(t *testing.T) {
	env := NewEnv(
		"test_env",
		0,
		"test_instance_type",
		[]EnvRepository{},
		EnvRuntimes{},
	)

	if len(env.ServedPorts) != 0 {
		t.Fatalf(
			"expected no served ports, got '%+v'",
			env.ServedPorts,
		)
	}

	servedPortToRemove := EnvServedPort("8000")

	env.RemoveServedPort(servedPortToRemove)

	if len(env.ServedPorts) != 0 {
		t.Fatalf(
			"expected no served ports, got '%+v'",
			env.ServedPorts,
		)
	}

	env.ServedPorts[servedPortToRemove] = []EnvServedPortBinding{
		{
			Type:  EnvServedPortBindingTypeDomain,
			Value: "api.eleven.sh",
		},
	}

	remainingServedPort := EnvServedPort("6000")
	remainingServedPortBindings := []EnvServedPortBinding{
		{
			Type:  EnvServedPortBindingTypeDomain,
			Value: "api.eleven.sh",
		},
	}

	env.ServedPorts[remainingServedPort] = remainingServedPortBindings

	env.RemoveServedPort(servedPortToRemove)

	if len(env.ServedPorts) != 1 || !reflect.DeepEqual(
		env.ServedPorts[remainingServedPort],
		remainingServedPortBindings,
	) {

		t.Fatalf(
			"expected served ports to equal '%+v', got '%+v'",
			remainingServedPortBindings,
			env.ServedPorts,
		)
	}

	env.RemoveServedPort(remainingServedPort)

	if len(env.ServedPorts) != 0 {
		t.Fatalf(
			"expected no served ports, got '%+v'",
			env.ServedPorts,
		)
	}
}

func TestEnvCheckPortValidityWithValidPorts(t *testing.T) {
	testCases := []struct {
		test          string
		port          string
		reservedPorts []string
	}{
		{
			test:          "with valid port",
			port:          "8080",
			reservedPorts: []string{"2200"},
		},

		{
			test:          "with minimum port",
			port:          "1",
			reservedPorts: []string{},
		},

		{
			test:          "with maximum port",
			port:          "65535",
			reservedPorts: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckPortValidity(
				tc.port,
				tc.reservedPorts,
			)

			if err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}
		})
	}
}

func TestEnvCheckPortValidityWithInvalidPorts(t *testing.T) {
	testCases := []struct {
		test          string
		port          string
		reservedPorts []string
	}{
		{
			test:          "with invalid port",
			port:          "invalid_port",
			reservedPorts: []string{"2200"},
		},

		{
			test:          "with less than minimum port",
			port:          "0",
			reservedPorts: []string{},
		},

		{
			test:          "with more than maximum port",
			port:          "65536",
			reservedPorts: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckPortValidity(
				tc.port,
				tc.reservedPorts,
			)

			if err == nil {
				t.Fatalf("expected error, got nothing")
			}

			if !errors.As(err, &ErrInvalidPort{}) {
				t.Fatalf("expected invalid port error, got '%+v'", err)
			}
		})
	}
}

func TestEnvCheckPortValidityWithReservedPorts(t *testing.T) {
	testCases := []struct {
		test          string
		port          string
		reservedPorts []string
	}{
		{
			test:          "with reserved port",
			port:          "2200",
			reservedPorts: []string{"2200", "8", "100"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckPortValidity(
				tc.port,
				tc.reservedPorts,
			)

			if err == nil {
				t.Fatalf("expected error, got nothing")
			}

			if !errors.As(err, &ErrReservedPort{}) {
				t.Fatalf("expected reserved port error, got '%+v'", err)
			}
		})
	}
}

func TestEnvCheckDomainValidityWithValidDomains(t *testing.T) {
	testCases := []struct {
		test   string
		domain string
	}{
		{
			test:   "with one level domain",
			domain: "test.com",
		},

		{
			test:   "with two levels domain",
			domain: "www.test.com",
		},

		{
			test:   "with three levels domain",
			domain: "bar.www.test.com",
		},

		{
			test:   "with four levels domain",
			domain: "baz.bar.www.test.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckDomainValidity(
				tc.domain,
			)

			if err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}
		})
	}
}

func TestEnvCheckDomainValidityWithInvalidDomains(t *testing.T) {
	testCases := []struct {
		test   string
		domain string
	}{
		{
			test:   "without ext",
			domain: "localhost",
		},

		{
			test:   "with invalid characters",
			domain: "invalid_domain.com",
		},

		{
			test:   "with invalid length",
			domain: strings.Repeat("a", 255) + ".com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckDomainValidity(
				tc.domain,
			)

			if err == nil {
				t.Fatalf("expected error, got nothing")
			}

			if !errors.As(err, &ErrInvalidDomain{}) {
				t.Fatalf("expected invalid domain error, got '%+v'", err)
			}
		})
	}
}
