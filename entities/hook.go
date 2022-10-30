package entities

type HookRunner interface {
	Run(
		cloudService CloudService,
		config *Config,
		cluster *Cluster,
		env *Env,
	) error
}

type DomainReachabilityChecker interface {
	Check(
		env *Env,
		domain string,
	) (reachable bool, redirToHTTPS bool, err error)
}
