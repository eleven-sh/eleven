package features

import (
	"fmt"

	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

type EditInput struct {
	EnvName string
}

type EditOutput struct {
	Error   error
	Content *EditOutputContent
	Stepper stepper.Stepper
}

type EditOutputContent struct {
	Cluster *entities.Cluster
	Env     *entities.Env
}

type EditOutputHandler interface {
	HandleOutput(EditOutput) error
}

type EditFeature struct {
	stepper             stepper.Stepper
	outputHandler       EditOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewEditFeature(
	stepper stepper.Stepper,
	outputHandler EditOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) EditFeature {

	return EditFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (e EditFeature) Execute(input EditInput) error {
	handleError := func(err error) error {
		e.outputHandler.HandleOutput(EditOutput{
			Stepper: e.stepper,
			Error:   err,
		})

		return err
	}

	envName := input.EnvName

	e.stepper.StartTemporaryStep(
		fmt.Sprintf("Editing the sandbox \"%s\"", envName),
	)

	cloudService, err := e.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	elevenConfig, err := cloudService.LookupElevenConfig(
		e.stepper,
	)

	if err != nil {
		return handleError(err)
	}

	clusterName := entities.DefaultClusterName
	cluster, err := elevenConfig.GetCluster(clusterName)

	if err != nil {
		return handleError(err)
	}

	env, err := elevenConfig.GetEnv(cluster.Name, envName)

	if err != nil {
		return handleError(err)
	}

	if env.Status == entities.EnvStatusRemoving {
		return handleError(entities.ErrEditRemovingEnv{
			EnvName: envName,
		})
	}

	if env.Status == entities.EnvStatusCreating {
		return handleError(entities.ErrEditCreatingEnv{
			EnvName: envName,
		})
	}

	return e.outputHandler.HandleOutput(EditOutput{
		Stepper: e.stepper,
		Content: &EditOutputContent{
			Cluster: cluster,
			Env:     env,
		},
	})
}
