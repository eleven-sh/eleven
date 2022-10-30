package actions

import (
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

func OpenPort(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
	portToOpen string,
) error {

	openPortErr := cloudService.OpenPort(
		stepper,
		elevenConfig,
		cluster,
		env,
		portToOpen,
	)

	// "openPortErr" is not handled first
	// in order to be able to save partial infrastructure
	err := UpdateEnvInConfig(
		stepper,
		cloudService,
		elevenConfig,
		cluster,
		env,
	)

	if err != nil {
		return err
	}

	return openPortErr
}
