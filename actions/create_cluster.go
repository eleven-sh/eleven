package actions

import (
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

func CreateCluser(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	createClusterErr := cloudService.CreateCluster(
		stepper,
		elevenConfig,
		cluster,
	)

	// "createCLusterErr" is not handled first
	// in order to be able to save partial infrastructure
	err := UpdateClusterInConfig(
		stepper,
		cloudService,
		elevenConfig,
		cluster,
	)

	if err != nil {
		return err
	}

	if createClusterErr != nil {
		return createClusterErr
	}

	cluster.Status = entities.ClusterStatusCreated
	return UpdateClusterInConfig(
		stepper,
		cloudService,
		elevenConfig,
		cluster,
	)
}
