package entities

import (
	"errors"
	"testing"
)

func TestConfigSetEnv(t *testing.T) {
	config := NewConfig()
	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)
	env := NewEnv(
		"env_name",
		0,
		"instance_type",
		[]EnvRepository{},
		EnvRuntimes{},
	)

	err := config.SetEnv(cluster.Name, env)

	if err == nil || !errors.As(err, &ErrClusterNotExists{}) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			ErrClusterNotExists{},
			err,
		)
	}

	typedError := err.(ErrClusterNotExists)

	if typedError.ClusterName != cluster.Name {
		t.Fatalf(
			"expected error cluster name to equal '%s', got '%s'",
			cluster.Name,
			typedError.ClusterName,
		)
	}

	config.Clusters[cluster.Name] = cluster

	err = config.SetEnv(cluster.Name, nil)

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	err = config.SetEnv(cluster.Name, env)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	_, envSet := config.Clusters[cluster.Name].Envs[env.Name]

	if !envSet {
		t.Fatalf("expected env to be set")
	}
}

func TestConfigEnvExists(t *testing.T) {
	config := NewConfig()
	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)
	env := NewEnv(
		"env_name",
		0,
		"instance_type",
		[]EnvRepository{},
		EnvRuntimes{},
	)

	envExists := config.EnvExists(cluster.Name, env.Name)

	if envExists {
		t.Fatalf("expected env to not exist")
	}

	config.Clusters[cluster.Name] = cluster
	envExists = config.EnvExists(cluster.Name, env.Name)

	if envExists {
		t.Fatalf("expected env to not exist")
	}

	config.Clusters[cluster.Name].Envs[env.Name] = env
	envExists = config.EnvExists(cluster.Name, env.Name)

	if !envExists {
		t.Fatalf("expected env to exist")
	}
}

func TestConfigGetEnv(t *testing.T) {
	config := NewConfig()
	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)
	env := NewEnv(
		"env_name",
		0,
		"instance_type",
		[]EnvRepository{},
		EnvRuntimes{},
	)

	returnedEnv, err := config.GetEnv(cluster.Name, env.Name)

	if returnedEnv != nil {
		t.Fatalf("expected env to not exist")
	}

	if err == nil || !errors.As(err, &ErrEnvNotExists{}) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			ErrEnvNotExists{},
			err,
		)
	}

	typedError := err.(ErrEnvNotExists)

	if typedError.ClusterName != cluster.Name {
		t.Fatalf(
			"expected error cluster name to equal '%s', got '%s'",
			cluster.Name,
			typedError.ClusterName,
		)
	}

	if typedError.EnvName != env.Name {
		t.Fatalf(
			"expected error env name to equal '%s', got '%s'",
			env.Name,
			typedError.EnvName,
		)
	}

	config.Clusters[cluster.Name] = cluster
	returnedEnv, err = config.GetEnv(cluster.Name, env.Name)

	if returnedEnv != nil {
		t.Fatalf("expected env to not exist")
	}

	if err == nil || !errors.As(err, &ErrEnvNotExists{}) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			ErrEnvNotExists{},
			err,
		)
	}

	typedError = err.(ErrEnvNotExists)

	if typedError.ClusterName != cluster.Name {
		t.Fatalf(
			"expected error cluster name to equal '%s', got '%s'",
			cluster.Name,
			typedError.ClusterName,
		)
	}

	if typedError.EnvName != env.Name {
		t.Fatalf(
			"expected error env name to equal '%s', got '%s'",
			env.Name,
			typedError.EnvName,
		)
	}

	config.Clusters[cluster.Name].Envs[env.Name] = env
	returnedEnv, err = config.GetEnv(cluster.Name, env.Name)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if returnedEnv == nil {
		t.Fatalf("expected env to exist")
	}
}

func TestConfigRemoveEnv(t *testing.T) {
	config := NewConfig()
	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)
	env := NewEnv(
		"env_name",
		0,
		"instance_type",
		[]EnvRepository{},
		EnvRuntimes{},
	)

	err := config.RemoveEnv(cluster.Name, env.Name)

	if err == nil || !errors.As(err, &ErrEnvNotExists{}) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			ErrEnvNotExists{},
			err,
		)
	}

	typedError := err.(ErrEnvNotExists)

	if typedError.ClusterName != cluster.Name {
		t.Fatalf(
			"expected error cluster name to equal '%s', got '%s'",
			cluster.Name,
			typedError.ClusterName,
		)
	}

	if typedError.EnvName != env.Name {
		t.Fatalf(
			"expected error env name to equal '%s', got '%s'",
			env.Name,
			typedError.EnvName,
		)
	}

	config.Clusters[cluster.Name] = cluster
	err = config.RemoveEnv(cluster.Name, env.Name)

	if err == nil || !errors.As(err, &ErrEnvNotExists{}) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			ErrEnvNotExists{},
			err,
		)
	}

	typedError = err.(ErrEnvNotExists)

	if typedError.ClusterName != cluster.Name {
		t.Fatalf(
			"expected error cluster name to equal '%s', got '%s'",
			cluster.Name,
			typedError.ClusterName,
		)
	}

	if typedError.EnvName != env.Name {
		t.Fatalf(
			"expected error env name to equal '%s', got '%s'",
			env.Name,
			typedError.EnvName,
		)
	}

	config.Clusters[cluster.Name].Envs[env.Name] = env
	err = config.RemoveEnv(cluster.Name, env.Name)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	_, envExists := config.Clusters[cluster.Name].Envs[env.Name]
	if envExists {
		t.Fatalf("expected env to not exist")
	}
}

func TestConfigCountEnvsInCluster(t *testing.T) {
	config := NewConfig()
	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)
	env := NewEnv(
		"env_name",
		0,
		"instance_type",
		[]EnvRepository{},
		EnvRuntimes{},
	)

	count, err := config.CountEnvsInCluster(cluster.Name)

	if err == nil || !errors.As(err, &ErrClusterNotExists{}) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			ErrClusterNotExists{},
			err,
		)
	}

	typedError := err.(ErrClusterNotExists)

	if typedError.ClusterName != cluster.Name {
		t.Fatalf(
			"expected error cluster name to equal '%s', got '%s'",
			cluster.Name,
			typedError.ClusterName,
		)
	}

	if count != 0 {
		t.Fatalf(
			"expected count to equal '0', got '%d'",
			count,
		)
	}

	config.Clusters[cluster.Name] = cluster

	count, err = config.CountEnvsInCluster(cluster.Name)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if count != 0 {
		t.Fatalf(
			"expected count to equal '0', got '%d'",
			count,
		)
	}

	config.Clusters[cluster.Name].Envs[env.Name] = env

	count, err = config.CountEnvsInCluster(cluster.Name)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if count != 1 {
		t.Fatalf(
			"expected count to equal '1', got '%d'",
			count,
		)
	}
}
