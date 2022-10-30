package actions

import (
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

func RemoveEnv(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
	preRemoveHook entities.HookRunner,
) error {

	env.Status = entities.EnvStatusRemoving
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

	removeEnvErr := cloudService.RemoveEnv(
		stepper,
		elevenConfig,
		cluster,
		env,
	)

	// "removeEnvErr" is not handled first
	// in order to be able to save partial infrastructure
	err = UpdateEnvInConfig(
		stepper,
		cloudService,
		elevenConfig,
		cluster,
		env,
	)

	if err != nil {
		return err
	}

	if removeEnvErr != nil {
		return removeEnvErr
	}

	if preRemoveHook != nil {
		err = preRemoveHook.Run(
			cloudService,
			elevenConfig,
			cluster,
			env,
		)

		if err != nil {
			return err
		}
	}

	return RemoveEnvInConfig(
		stepper,
		cloudService,
		elevenConfig,
		cluster,
		env,
	)
}
