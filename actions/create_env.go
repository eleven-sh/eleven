package actions

import (
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

func CreateEnv(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
) error {

	createEnvErr := cloudService.CreateEnv(
		stepper,
		elevenConfig,
		cluster,
		env,
	)

	// "createEnvErr" is not handled first
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

	return createEnvErr
}
