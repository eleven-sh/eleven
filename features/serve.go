package features

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/eleven-sh/eleven/actions"
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
)

type ServeInput struct {
	EnvName                   string
	ReservedPorts             []string
	Port                      string
	PortBinding               string
	DomainReachabilityChecker entities.DomainReachabilityChecker
}

type ServeOutput struct {
	Error   error
	Content *ServeOutputContent
	Stepper stepper.Stepper
}

type ServeOutputContent struct {
	Cluster     *entities.Cluster
	Env         *entities.Env
	Port        string
	PortBinding string
}

type ServeOutputHandler interface {
	HandleOutput(ServeOutput) error
}

type ServeFeature struct {
	stepper             stepper.Stepper
	outputHandler       ServeOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewServeFeature(
	stepper stepper.Stepper,
	outputHandler ServeOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) ServeFeature {

	return ServeFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (s ServeFeature) Execute(input ServeInput) error {
	handleError := func(err error) error {
		s.outputHandler.HandleOutput(ServeOutput{
			Stepper: s.stepper,
			Error:   err,
		})

		return err
	}

	step := fmt.Sprintf(
		"Serving port \"%s\"",
		input.Port,
	)

	if len(input.PortBinding) > 0 {
		step = fmt.Sprintf(
			"Serving port \"%s\" as \"%s\"",
			input.Port,
			input.PortBinding,
		)
	}

	s.stepper.StartTemporaryStep(step)

	err := entities.CheckPortValidity(
		input.Port,
		input.ReservedPorts,
	)

	if err != nil {
		return handleError(err)
	}

	if len(input.PortBinding) > 0 {
		// Bindings are always domains for now
		err := entities.CheckDomainValidity(input.PortBinding)

		if err != nil {
			return handleError(err)
		}
	}

	cloudService, err := s.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	elevenConfig, err := cloudService.LookupElevenConfig(
		s.stepper,
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
		return handleError(entities.ErrServeRemovingEnv{
			EnvName: envName,
		})
	}

	if env.Status == entities.EnvStatusCreating {
		return handleError(entities.ErrServeCreatingEnv{
			EnvName: envName,
		})
	}

	portBinding := input.PortBinding
	if len(portBinding) == 0 {
		// if no binding was passed,
		// the user wants port to be served
		// using same port number
		portBinding = input.Port
	}

	portBindingAlreadyUsed := env.DoesServedPortBindingExist(portBinding)

	if !portBindingAlreadyUsed && govalidator.IsPort(portBinding) {
		err = actions.OpenPort(
			s.stepper,
			cloudService,
			elevenConfig,
			cluster,
			env,
			portBinding,
		)

		if err != nil {
			return handleError(err)
		}
	}

	redirPortBindingHTTPS := false

	if !govalidator.IsPort(portBinding) {
		step := fmt.Sprintf(
			"Checking that \"%s\" resolves to your sandbox's public IP address",
			portBinding,
		)

		s.stepper.StartTemporaryStep(step)

		reachable, redirToHTTPS, err := input.DomainReachabilityChecker.Check(
			env,
			portBinding,
		)

		if err != nil {
			return handleError(err)
		}

		if !reachable {
			return handleError(entities.ErrUnresolvableDomain{
				Domain:       portBinding,
				EnvIPAddress: env.InstancePublicIPAddress,
			})
		}

		redirPortBindingHTTPS = redirToHTTPS
	}

	if portBindingAlreadyUsed {
		env.RemoveServedPortBinding(
			entities.EnvServedPort(input.Port),
			portBinding,
		)
	}

	env.AddServedPortBinding(
		entities.EnvServedPort(input.Port),
		portBinding,
		redirPortBindingHTTPS,
	)

	err = actions.UpdateEnvInConfig(
		s.stepper,
		cloudService,
		elevenConfig,
		cluster,
		env,
	)

	if err != nil {
		return handleError(err)
	}

	return s.outputHandler.HandleOutput(ServeOutput{
		Stepper: s.stepper,
		Content: &ServeOutputContent{
			Cluster:     cluster,
			Env:         env,
			Port:        input.Port,
			PortBinding: portBinding,
		},
	})
}
