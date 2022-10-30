package actions

import (
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

func RemoveCluster(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	cluster.Status = entities.ClusterStatusRemoving
	err := UpdateClusterInConfig(
		stepper,
		cloudService,
		elevenConfig,
		cluster,
	)

	if err != nil {
		return err
	}

	removeClusterErr := cloudService.RemoveCluster(
		stepper,
		elevenConfig,
		cluster,
	)

	// "removeClusterErr" is not handled first
	// in order to be able to save partial infrastructure
	err = UpdateClusterInConfig(
		stepper,
		cloudService,
		elevenConfig,
		cluster,
	)

	if err != nil {
		return err
	}

	if removeClusterErr != nil {
		return removeClusterErr
	}

	return RemoveClusterInConfig(
		stepper,
		cloudService,
		elevenConfig,
		cluster,
	)
}
