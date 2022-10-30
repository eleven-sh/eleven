package actions

import (
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

func InstallEleven(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	elevenConfig *entities.Config,
) error {

	err := cloudService.CreateElevenConfigStorage(stepper)

	if err != nil {
		return err
	}

	return cloudService.SaveElevenConfig(
		stepper,
		elevenConfig,
	)
}
