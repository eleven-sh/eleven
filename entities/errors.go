package entities

import "errors"

var (
	ErrElevenNotInstalled    = errors.New("ErrElevenNotInstalled")
	ErrUninstallExistingEnvs = errors.New("ErrUninstallExistingEnvs")
)
