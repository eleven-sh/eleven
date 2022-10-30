package entities

import "errors"

func (c *Config) SetEnv(clusterName string, env *Env) error {
	if !c.ClusterExists(clusterName) {
		return ErrClusterNotExists{
			ClusterName: clusterName,
		}
	}

	if env == nil {
		return errors.New("passed env is nil")
	}

	c.Clusters[clusterName].Envs[env.Name] = env

	return nil
}

func (c *Config) EnvExists(clusterName, envName string) bool {
	if !c.ClusterExists(clusterName) {
		return false
	}

	_, envExists := c.Clusters[clusterName].Envs[envName]

	return envExists
}

func (c *Config) GetEnv(clusterName, envName string) (*Env, error) {
	if !c.EnvExists(clusterName, envName) {
		return nil, ErrEnvNotExists{
			ClusterName: clusterName,
			EnvName:     envName,
		}
	}

	env := c.Clusters[clusterName].Envs[envName]

	return env, nil
}

func (c *Config) RemoveEnv(clusterName, envName string) error {
	if !c.EnvExists(clusterName, envName) {
		return ErrEnvNotExists{
			ClusterName: clusterName,
			EnvName:     envName,
		}
	}

	delete(c.Clusters[clusterName].Envs, envName)

	return nil
}

func (c *Config) CountEnvsInCluster(clusterName string) (int, error) {
	if !c.ClusterExists(clusterName) {
		return 0, ErrClusterNotExists{
			ClusterName: clusterName,
		}
	}

	return len(c.Clusters[clusterName].Envs), nil
}
