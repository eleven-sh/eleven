package entities

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"golang.org/x/crypto/ssh"
)

const (
	EnvNameRegExp    = `^[a-z0-9]+(-[a-z0-9]+)*$`
	EnvNameMaxLength = 16
)

type EnvSSHHostKey struct {
	Algorithm   string `json:"algorithm"`
	Fingerprint string `json:"fingerprint"`
}

type EnvStatus string

const (
	EnvStatusCreating EnvStatus = "creating"
	EnvStatusCreated  EnvStatus = "created"
	EnvStatusRemoving EnvStatus = "removing"
)

type Env struct {
	ID                       string          `json:"id"`
	Name                     string          `json:"name"`
	LocalSSHConfigHostname   string          `json:"local_ssh_config_hostname"`
	InfrastructureJSON       string          `json:"infrastructure_json"`
	InstanceType             string          `json:"instance_type"`
	InstancePublicIPAddress  string          `json:"instance_public_ip_address"`
	SSHHostKeys              []EnvSSHHostKey `json:"ssh_host_keys"`
	SSHKeyPairPEMContent     string          `json:"ssh_key_pair_pem_content"`
	Repositories             []EnvRepository `json:"repositories"`
	Runtimes                 EnvRuntimes     `json:"runtimes"`
	ServedPorts              EnvServedPorts  `json:"served_ports"`
	Status                   EnvStatus       `json:"status"`
	AdditionalPropertiesJSON string          `json:"additional_properties_json"`
	CreatedAtTimestamp       int64           `json:"created_at_timestamp"`
}

func NewEnv(
	envName string,
	localSSHCfgDupHostCt int,
	instanceType string,
	repositories []EnvRepository,
	runtimes EnvRuntimes,
) *Env {

	return &Env{
		ID:   uuid.NewString(),
		Name: envName,
		LocalSSHConfigHostname: buildLocalSSHCfgHostnameForEnv(
			envName,
			localSSHCfgDupHostCt,
		),
		InstanceType:       instanceType,
		SSHHostKeys:        []EnvSSHHostKey{},
		Repositories:       repositories,
		Runtimes:           runtimes,
		ServedPorts:        EnvServedPorts{},
		Status:             EnvStatusCreating,
		CreatedAtTimestamp: time.Now().Unix(),
	}
}

func (e *Env) GetNameSlug() string {
	return BuildSlugForEnv(e.Name)
}

func (e *Env) GetSSHKeyPairName() string {
	return BuildSlugForEnv(e.LocalSSHConfigHostname)
}

func (e *Env) SetInfrastructureJSON(infrastructure interface{}) error {
	infrastructureJSON, err := json.Marshal(infrastructure)

	if err != nil {
		return err
	}

	e.InfrastructureJSON = string(infrastructureJSON)

	return nil
}

func (e *Env) SetAdditionalPropertiesJSON(additionalProps interface{}) error {
	additionalPropsJSON, err := json.Marshal(additionalProps)

	if err != nil {
		return err
	}

	e.AdditionalPropertiesJSON = string(additionalPropsJSON)

	return nil
}

func BuildInitialLocalSSHCfgHostnameForEnv(envName string) string {
	return "eleven/" + BuildSlugForEnv(envName)
}

func buildLocalSSHCfgHostnameForEnv(
	envName string,
	duplicateHostnamesCount int,
) string {

	prefix := BuildInitialLocalSSHCfgHostnameForEnv(envName)

	if duplicateHostnamesCount == 0 {
		return prefix
	}

	return prefix + "-" + fmt.Sprintf("%d", duplicateHostnamesCount)
}

func BuildSlugForEnv(rawString string) string {
	return strings.ReplaceAll(slug.Make(rawString), "_", "-")
}

func ParseSSHHostKeysForEnv(hostKeysContent string) ([]EnvSSHHostKey, error) {
	parsedHostKeys := []EnvSSHHostKey{}

	for len(hostKeysContent) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey([]byte(hostKeysContent))

		if err != nil {
			return nil, err
		}

		parsedHostKeys = append(parsedHostKeys, EnvSSHHostKey{
			Algorithm:   pubKey.Type(),
			Fingerprint: base64.StdEncoding.EncodeToString(pubKey.Marshal()),
		})

		hostKeysContent = string(rest)
	}

	return parsedHostKeys, nil
}

func CheckEnvNameValidity(envName string) error {
	validEnvName := govalidator.Matches(
		envName,
		EnvNameRegExp,
	)

	if !validEnvName || len(envName) > EnvNameMaxLength {
		return ErrInvalidEnvName{
			EnvName:          envName,
			EnvNameRegExp:    EnvNameRegExp,
			EnvNameMaxLength: EnvNameMaxLength,
		}
	}

	return nil
}
