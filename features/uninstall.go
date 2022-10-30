package features

import (
	"errors"

	"github.com/eleven-sh/eleven/actions"
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

type UninstallInput struct {
	SuccessMessage            string
	AlreadyUninstalledMessage string
}

type UninstallOutput struct {
	Error   error
	Content *UninstallOutputContent
	Stepper stepper.Stepper
}

type UninstallOutputContent struct {
	ElevenAlreadyUninstalled  bool
	SuccessMessage            string
	AlreadyUninstalledMessage string
}

type UninstallOutputHandler interface {
	HandleOutput(UninstallOutput) error
}

type UninstallFeature struct {
	stepper             stepper.Stepper
	outputHandler       UninstallOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewUninstallFeature(
	stepper stepper.Stepper,
	outputHandler UninstallOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) UninstallFeature {

	return UninstallFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (u UninstallFeature) Execute(input UninstallInput) error {
	handleError := func(err error) error {
		u.outputHandler.HandleOutput(UninstallOutput{
			Stepper: u.stepper,
			Error:   err,
		})

		return err
	}

	u.stepper.StartTemporaryStep("Uninstalling Eleven")

	cloudService, err := u.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	elevenConfig, err := cloudService.LookupElevenConfig(
		u.stepper,
	)

	if err != nil {
		if errors.Is(err, entities.ErrElevenNotInstalled) {
			return u.outputHandler.HandleOutput(UninstallOutput{
				Stepper: u.stepper,
				Content: &UninstallOutputContent{
					ElevenAlreadyUninstalled:  true,
					SuccessMessage:            input.SuccessMessage,
					AlreadyUninstalledMessage: input.AlreadyUninstalledMessage,
				},
			})
		}

		return handleError(err)
	}

	clusterName := entities.DefaultClusterName
	cluster, err := elevenConfig.GetCluster(clusterName)

	// In case of error the Eleven config storage
	// could be created but without cluster
	if err != nil && !errors.As(err, &entities.ErrClusterNotExists{}) {
		return handleError(err)
	}

	if cluster != nil {
		nbOfEnvsInCluster, err := elevenConfig.CountEnvsInCluster(clusterName)

		if err != nil {
			return handleError(err)
		}

		if nbOfEnvsInCluster > 0 {
			return handleError(entities.ErrUninstallExistingEnvs)
		}

		err = actions.RemoveCluster(
			u.stepper,
			cloudService,
			elevenConfig,
			cluster,
		)

		if err != nil {
			return handleError(err)
		}
	}

	err = cloudService.RemoveElevenConfigStorage(
		u.stepper,
	)

	if err != nil {
		return handleError(err)
	}

	return u.outputHandler.HandleOutput(UninstallOutput{
		Stepper: u.stepper,
		Content: &UninstallOutputContent{
			ElevenAlreadyUninstalled:  false,
			SuccessMessage:            input.SuccessMessage,
			AlreadyUninstalledMessage: input.AlreadyUninstalledMessage,
		},
	})
}
