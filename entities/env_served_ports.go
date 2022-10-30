package entities

import (
	"github.com/asaskevich/govalidator"
)

type EnvServedPort string
type EnvServedPorts map[EnvServedPort][]EnvServedPortBinding

type EnvServedPortBinding struct {
	Value           string                   `json:"value"`
	Type            EnvServedPortBindingType `json:"type"`
	RedirectToHTTPS bool                     `json:"redirect_to_https"`
}

type EnvServedPortBindingType string

const (
	EnvServedPortBindingTypePort   EnvServedPortBindingType = "port"
	EnvServedPortBindingTypeDomain EnvServedPortBindingType = "domain"
)

func (e *Env) DoesServedPortExist(servedPort EnvServedPort) bool {
	_, servedPortExists := e.ServedPorts[servedPort]
	return servedPortExists
}

func (e *Env) AddServedPortBinding(
	servedPort EnvServedPort,
	binding string,
	redirectToHTTPS bool,
) {

	if !e.DoesServedPortExist(servedPort) {
		e.ServedPorts[servedPort] = []EnvServedPortBinding{}
	}

	servedPortBinding := EnvServedPortBinding{
		Value:           binding,
		Type:            EnvServedPortBindingTypeDomain,
		RedirectToHTTPS: redirectToHTTPS,
	}

	if govalidator.IsPort(binding) {
		servedPortBinding.Type = EnvServedPortBindingTypePort
	}

	e.ServedPorts[servedPort] = append(e.ServedPorts[servedPort], servedPortBinding)
}

func (e *Env) DoesServedPortBindingExist(targetBinding string) bool {
	for _, bindings := range e.ServedPorts {
		for _, binding := range bindings {
			if targetBinding == binding.Value {
				return true
			}
		}
	}

	return false
}

func (e *Env) RemoveServedPortBinding(servedPort EnvServedPort, targetBinding string) {
	if !e.DoesServedPortExist(servedPort) {
		return
	}

	servedPortBindings := []EnvServedPortBinding{}

	for _, binding := range e.ServedPorts[servedPort] {
		if binding.Value == targetBinding {
			continue
		}

		servedPortBindings = append(servedPortBindings, binding)
	}

	e.ServedPorts[servedPort] = servedPortBindings
}

func (e *Env) RemoveServedPort(servedPort EnvServedPort) {
	delete(e.ServedPorts, servedPort)
}

func CheckPortValidity(port string, reservedPorts []string) error {
	if !govalidator.IsPort(port) {
		return ErrInvalidPort{
			InvalidPort: port,
		}
	}

	for _, reservedPort := range reservedPorts {
		if reservedPort == port {
			return ErrReservedPort{
				ReservedPort: port,
			}
		}
	}

	return nil
}

func CheckDomainValidity(domain string) error {
	validDomain := govalidator.Matches(
		domain,
		`^([a-z0-9]{1}[a-z0-9-]{0,62}){1}(\.[a-z0-9]{1}[a-z0-9-]{0,62})+$`,
	)

	if !validDomain || len(domain) > 253 {
		return ErrInvalidDomain{
			Domain: domain,
		}
	}

	return nil
}
