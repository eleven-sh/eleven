package entities

type ErrInvalidPort struct {
	InvalidPort string
}

func (ErrInvalidPort) Error() string {
	return "ErrInvalidPort"
}

type ErrReservedPort struct {
	ReservedPort string
}

func (ErrReservedPort) Error() string {
	return "ErrReservedPort"
}

type ErrServeRemovingEnv struct {
	EnvName string
}

func (ErrServeRemovingEnv) Error() string {
	return "ErrServeRemovingEnv"
}

type ErrServeCreatingEnv struct {
	EnvName string
}

func (ErrServeCreatingEnv) Error() string {
	return "ErrServeCreatingEnv"
}

type ErrUnserveRemovingEnv struct {
	EnvName string
}

func (ErrUnserveRemovingEnv) Error() string {
	return "ErrUnserveRemovingEnv"
}

type ErrUnserveCreatingEnv struct {
	EnvName string
}

func (ErrUnserveCreatingEnv) Error() string {
	return "ErrUnserveCreatingEnv"
}

type ErrInvalidDomain struct {
	Domain string
}

func (ErrInvalidDomain) Error() string {
	return "ErrInvalidDomain"
}

type ErrUnresolvableDomain struct {
	Domain       string
	EnvIPAddress string
}

func (ErrUnresolvableDomain) Error() string {
	return "ErrUnresolvableDomain"
}

type ErrCloudflareSSLFull struct {
	Domain string
}

func (ErrCloudflareSSLFull) Error() string {
	return "ErrCloudflareSSLFull"
}

type ErrProxyForceHTTPS struct {
	Domain string
}

func (ErrProxyForceHTTPS) Error() string {
	return "ErrProxyForceHTTPS"
}

type ErrLetsEncryptTimedOut struct {
	Domain        string
	ReturnedError error
}

func (ErrLetsEncryptTimedOut) Error() string {
	return "ErrLetsEncryptTimedOut"
}
