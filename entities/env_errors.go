package entities

type ErrEnvNotExists struct {
	ClusterName string
	EnvName     string
}

func (ErrEnvNotExists) Error() string {
	return "ErrEnvNotExists"
}

type ErrInitRemovingEnv struct {
	EnvName string
}

func (ErrInitRemovingEnv) Error() string {
	return "ErrInitRemovingEnv"
}

type ErrEditRemovingEnv struct {
	EnvName string
}

func (ErrEditRemovingEnv) Error() string {
	return "ErrEditRemovingEnv"
}

type ErrEditCreatingEnv struct {
	EnvName string
}

func (ErrEditCreatingEnv) Error() string {
	return "ErrEditCreatingEnv"
}

type ErrInvalidEnvName struct {
	EnvName          string
	EnvNameRegExp    string
	EnvNameMaxLength int
}

func (ErrInvalidEnvName) Error() string {
	return "ErrInvalidEnvName"
}

type ErrUpdateInstanceTypeCreatingEnv struct {
	EnvName string
}

func (ErrUpdateInstanceTypeCreatingEnv) Error() string {
	return "ErrUpdateInstanceTypeCreatingEnv"
}

type ErrEnvCloudInitError struct {
	Logs         string
	ErrorMessage string
}

func (ErrEnvCloudInitError) Error() string {
	return "ErrEnvCloudInitError"
}
