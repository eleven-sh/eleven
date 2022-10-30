package entities

import (
	"errors"
	"testing"
)

func TestConfigSetCluster(t *testing.T) {
	config := NewConfig()
	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)

	err := config.SetCluster(nil)

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	err = config.SetCluster(cluster)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	_, clusterSet := config.Clusters[cluster.Name]

	if !clusterSet {
		t.Fatalf("expected cluster to be set")
	}
}

func TestConfigClusterExists(t *testing.T) {
	config := NewConfig()
	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)

	clusterExists := config.ClusterExists(cluster.Name)

	if clusterExists {
		t.Fatalf("expected cluster to not exist")
	}

	config.Clusters[cluster.Name] = cluster
	clusterExists = config.ClusterExists(cluster.Name)

	if !clusterExists {
		t.Fatalf("expected cluster to exist")
	}
}

func TestConfigGetCluster(t *testing.T) {
	config := NewConfig()
	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)

	returnedCluster, err := config.GetCluster(cluster.Name)

	if returnedCluster != nil {
		t.Fatalf("expected cluster to not exist")
	}

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
			"expected cluster name to equal '%s', got '%s'",
			cluster.Name,
			typedError.ClusterName,
		)
	}

	config.Clusters[cluster.Name] = cluster
	returnedCluster, err = config.GetCluster(cluster.Name)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if returnedCluster == nil || returnedCluster != cluster {
		t.Fatalf("expected cluster to exist")
	}
}

func TestConfigRemoveCluster(t *testing.T) {
	config := NewConfig()
	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)

	err := config.RemoveCluster(cluster.Name)

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
			"expected cluster name to equal '%s', got '%s'",
			cluster.Name,
			typedError.ClusterName,
		)
	}

	config.Clusters[cluster.Name] = cluster
	err = config.RemoveCluster(cluster.Name)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	_, clusterExists := config.Clusters[cluster.Name]

	if clusterExists {
		t.Fatalf("expected cluster to not exist")
	}
}
