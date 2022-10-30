package entities

type ErrEnvInvalidRuntime struct {
	Runtime string
}

func (ErrEnvInvalidRuntime) Error() string {
	return "ErrEnvInvalidRuntime"
}

type ErrEnvDuplicatedRuntimes struct {
	Runtime string
}

func (ErrEnvDuplicatedRuntimes) Error() string {
	return "ErrEnvDuplicatedRuntimes"
}

type ErrEnvInvalidRuntimeVersion struct {
	Runtime                string
	RuntimeVersion         string
	RuntimeVersionExamples []string
}

func (ErrEnvInvalidRuntimeVersion) Error() string {
	return "ErrEnvInvalidRuntimeVersion"
}
