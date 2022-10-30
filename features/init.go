package features

import (
	"errors"
	"fmt"

	"github.com/eleven-sh/eleven/actions"
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

type InitInput struct {
	InstanceType         string
	EnvName              string
	LocalSSHCfgDupHostCt int
	Repositories         []entities.EnvRepository
	Runtimes             []string
}

type InitOutput struct {
	Error   error
	Content *InitOutputContent
	Stepper stepper.Stepper
}

type InitOutputContent struct {
	CloudService    entities.CloudService
	ElevenConfig    *entities.Config
	Cluster         *entities.Cluster
	Env             *entities.Env
	EnvCreated      bool
	SetEnvAsCreated func() error
	Runtimes        entities.EnvRuntimes
}

type InitOutputHandler interface {
	HandleOutput(InitOutput) error
}

type InitFeature struct {
	stepper             stepper.Stepper
	outputHandler       InitOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewInitFeature(
	stepper stepper.Stepper,
	outputHandler InitOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) InitFeature {

	return InitFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (i InitFeature) Execute(input InitInput) error {
	handleError := func(err error) error {
		i.outputHandler.HandleOutput(InitOutput{
			Stepper: i.stepper,
			Error:   err,
		})

		return err
	}

	envName := input.EnvName

	step := fmt.Sprintf("Initializing the sandbox \"%s\"", envName)
	i.stepper.StartTemporaryStep(step)

	err := entities.CheckEnvNameValidity(input.EnvName)

	if err != nil {
		return handleError(err)
	}

	runtimes, err := entities.ParseEnvRuntimes(input.Runtimes)

	if err != nil {
		return handleError(err)
	}

	cloudService, err := i.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	err = cloudService.CheckInstanceTypeValidity(
		i.stepper,
		input.InstanceType,
	)

	if err != nil {
		return handleError(err)
	}

	elevenConfig, err := cloudService.LookupElevenConfig(
		i.stepper,
	)

	if err != nil && !errors.Is(err, entities.ErrElevenNotInstalled) {
		return handleError(err)
	}

	if elevenConfig == nil { // Eleven not installed

		i.stepper.StartTemporaryStep("Installing Eleven")

		elevenConfig = entities.NewConfig()

		err = actions.InstallEleven(
			i.stepper,
			cloudService,
			elevenConfig,
		)

		if err != nil {
			return handleError(err)
		}
	}

	clusterName := entities.DefaultClusterName
	cluster, err := elevenConfig.GetCluster(clusterName)

	if err != nil && !errors.As(err, &entities.ErrClusterNotExists{}) {
		return handleError(err)
	}

	if cluster == nil || cluster.Status == entities.ClusterStatusCreating {

		/* Cluster not exists or still
		in creating state after error */

		i.stepper.StartTemporaryStep("Creating default cluster")

		if cluster == nil {
			// Multiple clusters are not implemented for now
			isDefaultCluster := true

			cluster = entities.NewCluster(
				clusterName,
				input.InstanceType,
				isDefaultCluster,
			)
		}

		err = actions.CreateCluser(
			i.stepper,
			cloudService,
			elevenConfig,
			cluster,
		)

		if err != nil {
			return handleError(err)
		}
	}

	env, err := elevenConfig.GetEnv(
		cluster.Name,
		envName,
	)

	if err != nil && !errors.As(err, &entities.ErrEnvNotExists{}) {
		return handleError(err)
	}

	if env != nil && env.Status == entities.EnvStatusRemoving {
		return handleError(entities.ErrInitRemovingEnv{
			EnvName: env.Name,
		})
	}

	envCreated := false

	if env == nil || env.Status == entities.EnvStatusCreating {

		/* Env not exists or still
		in creating state after error */

		if env == nil {
			env = entities.NewEnv(
				envName,
				input.LocalSSHCfgDupHostCt,
				input.InstanceType,
				input.Repositories,
				runtimes,
			)
		} else {
			if env.InstanceType != input.InstanceType {
				return handleError(entities.ErrUpdateInstanceTypeCreatingEnv{
					EnvName: envName,
				})
			}

			env.Repositories = input.Repositories
			env.Runtimes = runtimes
		}

		err = actions.CreateEnv(
			i.stepper,
			cloudService,
			elevenConfig,
			cluster,
			env,
		)

		if err != nil {
			return handleError(err)
		}

		envCreated = true
	}

	// Current step is the last ended infrastructure step.
	// Better UX if we reset to main step here given that
	// the next steps (in GRPC agent) may take some time to start.
	i.stepper.StartTemporaryStep(step)

	setEnvAsCreated := func() error {
		env.Status = entities.EnvStatusCreated

		return actions.UpdateEnvInConfig(
			i.stepper,
			cloudService,
			elevenConfig,
			cluster,
			env,
		)
	}

	return i.outputHandler.HandleOutput(InitOutput{
		Stepper: i.stepper,
		Content: &InitOutputContent{
			CloudService:    cloudService,
			ElevenConfig:    elevenConfig,
			Cluster:         cluster,
			Env:             env,
			EnvCreated:      envCreated,
			SetEnvAsCreated: setEnvAsCreated,
			Runtimes:        runtimes,
		},
	})
}
