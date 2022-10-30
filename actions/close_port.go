package actions

import (
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

func ClosePort(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
	portToClose string,
) error {

	closePortErr := cloudService.ClosePort(
		stepper,
		elevenConfig,
		cluster,
		env,
		portToClose,
	)

	// "closePortErr" is not handled first
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

	return closePortErr
}
