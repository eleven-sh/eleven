package actions

import (
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

func UpdateClusterInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	err := elevenConfig.SetCluster(cluster)

	if err != nil {
		return err
	}

	return cloudService.SaveElevenConfig(
		stepper,
		elevenConfig,
	)
}

func RemoveClusterInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	err := elevenConfig.RemoveCluster(cluster.Name)

	if err != nil {
		return err
	}

	return cloudService.SaveElevenConfig(
		stepper,
		elevenConfig,
	)
}

func UpdateEnvInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
) error {

	err := elevenConfig.SetEnv(cluster.Name, env)

	if err != nil {
		return err
	}

	return cloudService.SaveElevenConfig(
		stepper,
		elevenConfig,
	)
}

func RemoveEnvInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
) error {

	err := elevenConfig.RemoveEnv(cluster.Name, env.Name)

	if err != nil {
		return err
	}

	return cloudService.SaveElevenConfig(
		stepper,
		elevenConfig,
	)
}
