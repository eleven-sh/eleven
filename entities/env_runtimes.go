package entities

import (
	"strings"

	"github.com/asaskevich/govalidator"
)

const envRuntimesLatestFlag = "latest"

type envAvailableRuntime struct {
	versionRegExp   string
	versionExamples []string
}

var (
	envRuntimesSemverRegExp     = `^((?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?|latest)$`
	envRuntimesLatestRegExp     = "^" + envRuntimesLatestFlag + "$"
	envRuntimesPHPVersionRegExp = `^((?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)|latest)$`

	envAvailableRuntimes = map[string]envAvailableRuntime{
		"clang": {
			versionRegExp:   envRuntimesLatestRegExp,
			versionExamples: []string{"latest"},
		},
		"docker": {
			versionRegExp:   envRuntimesLatestRegExp,
			versionExamples: []string{"latest"},
		},
		"go": {
			versionRegExp:   envRuntimesSemverRegExp,
			versionExamples: []string{"latest", "1.19.2", "1.19.0", "1.0.0"},
		},
		"java": {
			versionRegExp:   envRuntimesLatestRegExp,
			versionExamples: []string{"latest"},
		},
		"node": {
			versionRegExp:   envRuntimesSemverRegExp,
			versionExamples: []string{"latest", "18.11.0", "18.0.0", "16.18.0"},
		},
		"php": {
			versionRegExp:   envRuntimesPHPVersionRegExp,
			versionExamples: []string{"latest", "8.11", "8.0", "7.4"},
		},
		"python": {
			versionRegExp:   envRuntimesSemverRegExp,
			versionExamples: []string{"latest", "3.10.8", "3.0.0", "2.7.0"},
		},
		"ruby": {
			versionRegExp:   envRuntimesSemverRegExp,
			versionExamples: []string{"latest", "3.1.2", "3.0.0", "2.7.0"},
		},
		"rust": {
			versionRegExp:   envRuntimesSemverRegExp,
			versionExamples: []string{"latest", "1.64.0", "1.62.1", "1.0.0"},
		},
	}
)

type EnvRuntimes map[string]string

func ParseEnvRuntimes(runtimes []string) (EnvRuntimes, error) {
	parsedRuntimes := EnvRuntimes{}

	for _, runtime := range runtimes {
		runtimeName, runtimeVersion, _ := strings.Cut(runtime, "@")

		if _, runtimeAvailable := envAvailableRuntimes[runtimeName]; !runtimeAvailable {
			return nil, ErrEnvInvalidRuntime{
				Runtime: runtimeName,
			}
		}

		if _, runtimeAlreadyParsed := parsedRuntimes[runtimeName]; runtimeAlreadyParsed {
			return nil, ErrEnvDuplicatedRuntimes{
				Runtime: runtimeName,
			}
		}

		if len(runtimeVersion) == 0 {
			runtimeVersion = envRuntimesLatestFlag
		}

		validRuntimeVersion := govalidator.Matches(
			runtimeVersion,
			envAvailableRuntimes[runtimeName].versionRegExp,
		)

		if !validRuntimeVersion {
			return nil, ErrEnvInvalidRuntimeVersion{
				Runtime:                runtimeName,
				RuntimeVersion:         runtimeVersion,
				RuntimeVersionExamples: envAvailableRuntimes[runtimeName].versionExamples,
			}
		}

		parsedRuntimes[runtimeName] = runtimeVersion
	}

	return parsedRuntimes, nil
}
