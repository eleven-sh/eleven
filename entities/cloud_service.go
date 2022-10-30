package entities

import "github.com/eleven-sh/eleven/stepper"

type CloudService interface {
	CreateElevenConfigStorage(stepper.Stepper) error
	RemoveElevenConfigStorage(stepper.Stepper) error

	LookupElevenConfig(stepper.Stepper) (*Config, error)
	SaveElevenConfig(stepper.Stepper, *Config) error

	CreateCluster(stepper.Stepper, *Config, *Cluster) error
	RemoveCluster(stepper.Stepper, *Config, *Cluster) error

	CheckInstanceTypeValidity(stepper.Stepper, string) error

	CreateEnv(stepper.Stepper, *Config, *Cluster, *Env) error
	RemoveEnv(stepper.Stepper, *Config, *Cluster, *Env) error

	OpenPort(stepper.Stepper, *Config, *Cluster, *Env, string) error
	ClosePort(stepper.Stepper, *Config, *Cluster, *Env, string) error
}

type CloudServiceBuilder interface {
	Build() (CloudService, error)
}
