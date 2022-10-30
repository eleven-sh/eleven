package features

import (
	"fmt"

	"github.com/eleven-sh/eleven/actions"
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

type UnserveInput struct {
	EnvName       string
	ReservedPorts []string
	Port          string
}

type UnserveOutput struct {
	Error   error
	Content *UnserveOutputContent
	Stepper stepper.Stepper
}

type UnserveOutputContent struct {
	Cluster *entities.Cluster
	Env     *entities.Env
	Port    string
}

type UnserveOutputHandler interface {
	HandleOutput(UnserveOutput) error
}

type UnserveFeature struct {
	stepper             stepper.Stepper
	outputHandler       UnserveOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewUnserveFeature(
	stepper stepper.Stepper,
	outputHandler UnserveOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) UnserveFeature {

	return UnserveFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (u UnserveFeature) Execute(input UnserveInput) error {
	handleError := func(err error) error {
		u.outputHandler.HandleOutput(UnserveOutput{
			Stepper: u.stepper,
			Error:   err,
		})

		return err
	}

	u.stepper.StartTemporaryStep(
		fmt.Sprintf(
			"Unserving port \"%s\"",
			input.Port,
		),
	)

	err := entities.CheckPortValidity(
		input.Port,
		input.ReservedPorts,
	)

	if err != nil {
		return handleError(err)
	}

	cloudService, err := u.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	elevenConfig, err := cloudService.LookupElevenConfig(
		u.stepper,
	)

	if err != nil {
		return handleError(err)
	}

	clusterName := entities.DefaultClusterName
	cluster, err := elevenConfig.GetCluster(clusterName)

	if err != nil {
		return handleError(err)
	}

	envName := input.EnvName
	env, err := elevenConfig.GetEnv(cluster.Name, envName)

	if err != nil {
		return handleError(err)
	}

	if env.Status == entities.EnvStatusRemoving {
		return handleError(entities.ErrUnserveRemovingEnv{
			EnvName: envName,
		})
	}

	if env.Status == entities.EnvStatusCreating {
		return handleError(entities.ErrUnserveCreatingEnv{
			EnvName: envName,
		})
	}

	portAlreadyUnserved := !env.DoesServedPortExist(
		entities.EnvServedPort(input.Port),
	)

	if !portAlreadyUnserved {
		for _, binding := range env.ServedPorts[entities.EnvServedPort(input.Port)] {
			if binding.Type != entities.EnvServedPortBindingTypePort {
				continue
			}

			err = actions.ClosePort(
				u.stepper,
				cloudService,
				elevenConfig,
				cluster,
				env,
				binding.Value,
			)

			if err != nil {
				return handleError(err)
			}

			env.RemoveServedPortBinding(
				entities.EnvServedPort(input.Port),
				binding.Value,
			)

			err := actions.UpdateEnvInConfig(
				u.stepper,
				cloudService,
				elevenConfig,
				cluster,
				env,
			)

			if err != nil {
				return handleError(err)
			}
		}

		env.RemoveServedPort(entities.EnvServedPort(input.Port))

		err := actions.UpdateEnvInConfig(
			u.stepper,
			cloudService,
			elevenConfig,
			cluster,
			env,
		)

		if err != nil {
			return handleError(err)
		}
	}

	return u.outputHandler.HandleOutput(UnserveOutput{
		Stepper: u.stepper,
		Content: &UnserveOutputContent{
			Cluster: cluster,
			Env:     env,
			Port:    input.Port,
		},
	})
}
