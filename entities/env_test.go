package entities

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestEnvGetNameSlug(t *testing.T) {
	testCases := []struct {
		test         string
		envName      string
		expectedSlug string
	}{
		{
			test:         "with underscore",
			envName:      "env_name",
			expectedSlug: "env-name",
		},

		{
			test:         "with spaces",
			envName:      "this is the env name",
			expectedSlug: "this-is-the-env-name",
		},

		{
			test:         "with invalid characters",
			envName:      "this is !() the env ^`$ name",
			expectedSlug: "this-is-the-env-name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			env := NewEnv(
				tc.envName,
				0,
				"instance_type",
				[]EnvRepository{},
				EnvRuntimes{},
			)

			slug := env.GetNameSlug()

			if tc.expectedSlug != slug {
				t.Fatalf(
					"expected slug to equal '%+v', got '%+v'",
					tc.expectedSlug,
					slug,
				)
			}
		})
	}
}

func TestEnvGetSSHKeyPairName(t *testing.T) {
	testCases := []struct {
		test                   string
		envName                string
		duplicateSSHHostsCount int
		expectedSSHKeyPairName string
	}{
		{
			test:                   "with underscore",
			envName:                "env_name",
			duplicateSSHHostsCount: 0,
			expectedSSHKeyPairName: "eleven-env-name",
		},

		{
			test:                   "with spaces",
			envName:                "this is the env name",
			duplicateSSHHostsCount: 1,
			expectedSSHKeyPairName: "eleven-this-is-the-env-name-1",
		},

		{
			test:                   "with invalid characters",
			envName:                "this is !() the env ^`$ name",
			duplicateSSHHostsCount: 2,
			expectedSSHKeyPairName: "eleven-this-is-the-env-name-2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			env := NewEnv(
				tc.envName,
				tc.duplicateSSHHostsCount,
				"instance_type",
				[]EnvRepository{},
				EnvRuntimes{},
			)

			sshKeyPairName := env.GetSSHKeyPairName()

			if tc.expectedSSHKeyPairName != sshKeyPairName {
				t.Fatalf(
					"expected SSH key pair name to equal '%+v', got '%+v'",
					tc.expectedSSHKeyPairName,
					sshKeyPairName,
				)
			}
		})
	}
}

func TestEnvSetInfrastructureJSON(t *testing.T) {
	type envInfra struct {
		InstanceID string
	}

	expectedInfra := envInfra{
		InstanceID: "instance_id",
	}

	env := NewEnv(
		"env_name",
		0,
		"instance_type",
		[]EnvRepository{},
		EnvRuntimes{},
	)

	err := env.SetInfrastructureJSON(make(chan struct{}))

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	err = env.SetInfrastructureJSON(expectedInfra)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	var setInfra envInfra
	err = json.Unmarshal([]byte(env.InfrastructureJSON), &setInfra)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if !reflect.DeepEqual(setInfra, expectedInfra) {
		t.Fatalf(
			"expected infra to equal '%+v', got '%+v'",
			expectedInfra,
			setInfra,
		)
	}
}

func TestEnvSetAdditionalPropertiesJSON(t *testing.T) {
	type additionalProps struct {
		Prop string
	}

	expectedAdditionalProps := additionalProps{
		Prop: "prop",
	}

	env := NewEnv(
		"env_name",
		0,
		"instance_type",
		[]EnvRepository{},
		EnvRuntimes{},
	)

	err := env.SetAdditionalPropertiesJSON(make(chan struct{}))

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	err = env.SetAdditionalPropertiesJSON(expectedAdditionalProps)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	var setAdditionalProps additionalProps
	err = json.Unmarshal([]byte(env.AdditionalPropertiesJSON), &setAdditionalProps)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if !reflect.DeepEqual(setAdditionalProps, expectedAdditionalProps) {
		t.Fatalf(
			"expected additional props to equal '%+v', got '%+v'",
			expectedAdditionalProps,
			setAdditionalProps,
		)
	}
}

func TestBuildLocalSSHConfigHostnameForEnv(t *testing.T) {
	testCases := []struct {
		test                    string
		envName                 string
		duplicateHostnamesCount int
		expectedHostname        string
	}{
		{
			test:                    "without duplicate hostname",
			envName:                 "env-name",
			duplicateHostnamesCount: 0,
			expectedHostname:        "eleven/env-name",
		},

		{
			test:                    "with duplicate hostnames",
			envName:                 "env-name",
			duplicateHostnamesCount: 4,
			expectedHostname:        "eleven/env-name-4",
		},

		{
			test:                    "without duplicate hostname and invalid characters",
			envName:                 "env_name_4837''",
			duplicateHostnamesCount: 0,
			expectedHostname:        "eleven/env-name-4837",
		},

		{
			test:                    "with duplicate hostnames and invalid characters",
			envName:                 "env_name_4837''",
			duplicateHostnamesCount: 8,
			expectedHostname:        "eleven/env-name-4837-8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			hostname := buildLocalSSHCfgHostnameForEnv(
				tc.envName,
				tc.duplicateHostnamesCount,
			)

			if tc.expectedHostname != hostname {
				t.Fatalf(
					"expected hostname to equal '%+v', got '%+v'",
					tc.expectedHostname,
					hostname,
				)
			}
		})
	}
}

func TestBuildInitialLocalSSHCfgHostnameForEnv(t *testing.T) {
	testCases := []struct {
		test             string
		envName          string
		expectedHostname string
	}{
		{
			test:             "without duplicate hostname",
			envName:          "env-name",
			expectedHostname: "eleven/env-name",
		},

		{
			test:             "without duplicate hostname and invalid characters",
			envName:          "env_name_4837''",
			expectedHostname: "eleven/env-name-4837",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			hostname := BuildInitialLocalSSHCfgHostnameForEnv(
				tc.envName,
			)

			if tc.expectedHostname != hostname {
				t.Fatalf(
					"expected hostname to equal '%+v', got '%+v'",
					tc.expectedHostname,
					hostname,
				)
			}
		})
	}
}

func TestBuildSlugForEnv(t *testing.T) {
	testCases := []struct {
		test         string
		envName      string
		expectedSlug string
	}{
		{
			test:         "with underscore",
			envName:      "env_name",
			expectedSlug: "env-name",
		},

		{
			test:         "with spaces",
			envName:      "this is the env name",
			expectedSlug: "this-is-the-env-name",
		},

		{
			test:         "with invalid characters",
			envName:      "this is !() the env ^`$ name",
			expectedSlug: "this-is-the-env-name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			slug := BuildSlugForEnv(tc.envName)

			if tc.expectedSlug != slug {
				t.Fatalf(
					"expected slug to equal '%+v', got '%+v'",
					tc.expectedSlug,
					slug,
				)
			}
		})
	}
}

func TestParseValidSSHHostKeysForEnv(t *testing.T) {
	sshHostKeysContent := `ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQB/nAmOjTmezNUDKYvEeIRf2YnwM9/uUG1d0BYsc8/tRtx+RGi7N2lUbp728MXGwdnL9od4cItzky/zVdLZE2cycOa18xBK9cOWmcKS0A8FYBxEQWJ/q9YVUgZbFKfYGaGQxsER+A0w/fX8ALuk78ktP31K69LcQgxIsl7rNzxsoOQKJ/CIxOGMMxczYTiEoLvQhapFQMs3FL96didKr/QbrfB1WT6s3838SEaXfgZvLef1YB2xmfhbT9OXFE3FXvh2UPBfN+ffE7iiayQf/2XR+8j4N4bW30DiPtOQLGUrH1y5X/rpNZNlWW2+jGIxqZtgWg7lTy3mXy5x836Sj/6L jje.levy@gmail.com
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJqmVkvKmywIYkfXOWWPya3I1zAbWGwOGu9Q870Zh49v jeremylevy@macbook-pro-de-jeremy.home`
	expectedSSHHostKeys := []EnvSSHHostKey{
		{
			Algorithm:   "ssh-rsa",
			Fingerprint: "AAAAB3NzaC1yc2EAAAABJQAAAQB/nAmOjTmezNUDKYvEeIRf2YnwM9/uUG1d0BYsc8/tRtx+RGi7N2lUbp728MXGwdnL9od4cItzky/zVdLZE2cycOa18xBK9cOWmcKS0A8FYBxEQWJ/q9YVUgZbFKfYGaGQxsER+A0w/fX8ALuk78ktP31K69LcQgxIsl7rNzxsoOQKJ/CIxOGMMxczYTiEoLvQhapFQMs3FL96didKr/QbrfB1WT6s3838SEaXfgZvLef1YB2xmfhbT9OXFE3FXvh2UPBfN+ffE7iiayQf/2XR+8j4N4bW30DiPtOQLGUrH1y5X/rpNZNlWW2+jGIxqZtgWg7lTy3mXy5x836Sj/6L",
		},

		{
			Algorithm:   "ssh-ed25519",
			Fingerprint: "AAAAC3NzaC1lZDI1NTE5AAAAIJqmVkvKmywIYkfXOWWPya3I1zAbWGwOGu9Q870Zh49v",
		},
	}

	returnedSSHHostKeys, err := ParseSSHHostKeysForEnv(sshHostKeysContent)

	if err != nil {
		t.Fatalf(
			"expected no error, got %s",
			err,
		)
	}

	if !reflect.DeepEqual(expectedSSHHostKeys, returnedSSHHostKeys) {
		t.Fatalf(
			"expected SSH host keys to equal '%+v', got '%+v'",
			expectedSSHHostKeys,
			returnedSSHHostKeys,
		)
	}
}

func TestParseInvalidSSHHostKeysForEnv(t *testing.T) {
	sshHostKeysContent := "host_keys_content"
	returnedSSHHostKeys, err := ParseSSHHostKeysForEnv(sshHostKeysContent)

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	if returnedSSHHostKeys != nil {
		t.Fatalf(
			"expected no SSH host keys, got '%+v'",
			returnedSSHHostKeys,
		)
	}
}

func TestCheckEnvNameValidityWithValidNames(t *testing.T) {
	testCases := []struct {
		test    string
		envName string
	}{
		{
			test:    "without dash",
			envName: "environment",
		},

		{
			test:    "with one dash",
			envName: "env-name",
		},

		{
			test:    "with multiple dashes",
			envName: "my-env-name",
		},

		{
			test:    "with max length",
			envName: strings.Repeat("b", EnvNameMaxLength),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckEnvNameValidity(
				tc.envName,
			)

			if err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}
		})
	}
}

func TestCheckEnvNameValidityWithInvalidNames(t *testing.T) {
	testCases := []struct {
		test    string
		envName string
	}{
		{
			test:    "with invalid characters",
			envName: "env984_()'à-name",
		},

		{
			test:    "with underscores",
			envName: "env_name",
		},

		{
			test:    "with more than max length",
			envName: strings.Repeat("b", EnvNameMaxLength+1),
		},

		{
			test:    "starting with dash",
			envName: "-env-name",
		},

		{
			test:    "ending with dash",
			envName: "env-",
		},

		{
			test:    "starting and ending with dash",
			envName: "-env-name-",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckEnvNameValidity(
				tc.envName,
			)

			if err == nil || !errors.As(err, &ErrInvalidEnvName{}) {
				t.Fatalf(
					"expected error to equal '%+v', got '%+v'",
					ErrInvalidEnvName{},
					err,
				)
			}

			typedError := err.(ErrInvalidEnvName)

			if typedError.EnvName != tc.envName {
				t.Fatalf(
					"expected error env name to equal '%s', got '%s'",
					tc.envName,
					typedError.EnvName,
				)
			}

			if typedError.EnvNameRegExp != EnvNameRegExp {
				t.Fatalf(
					"expected error env name regexp to equal '%s', got '%s'",
					EnvNameRegExp,
					typedError.EnvNameRegExp,
				)
			}

			if typedError.EnvNameMaxLength != EnvNameMaxLength {
				t.Fatalf(
					"expected error env name max length to equal '%d', got '%d'",
					EnvNameMaxLength,
					typedError.EnvNameMaxLength,
				)
			}
		})
	}
}
