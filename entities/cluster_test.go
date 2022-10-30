package entities

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestClusterGetNameSlug(t *testing.T) {
	testCases := []struct {
		test         string
		clusterName  string
		expectedSlug string
	}{
		{
			test:         "with underscore",
			clusterName:  "cluster_name",
			expectedSlug: "cluster-name",
		},

		{
			test:         "with default cluster name",
			clusterName:  DefaultClusterName,
			expectedSlug: DefaultClusterName,
		},

		{
			test:         "with spaces",
			clusterName:  "this is the cluster name",
			expectedSlug: "this-is-the-cluster-name",
		},

		{
			test:         "with invalid characters",
			clusterName:  "this is !() the cluster ^`$ name",
			expectedSlug: "this-is-the-cluster-name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			cluster := NewCluster(
				tc.clusterName,
				"default_instance_type",
				true,
			)

			slug := cluster.GetNameSlug()

			if tc.expectedSlug != slug {
				t.Fatalf(
					"expected slug to equal '%+v', got '%+v'",
					tc.expectedSlug,
					slug,
				)
			}
		})
	}
}

func TestClusterSetInfrastructureJSON(t *testing.T) {
	type clusterInfra struct {
		VPCID    string
		SubnetID string
	}

	expectedInfra := clusterInfra{
		VPCID:    "vpc_id",
		SubnetID: "subnet_id",
	}

	cluster := NewCluster(
		"cluster_name",
		"default_instance_type",
		true,
	)

	err := cluster.SetInfrastructureJSON(make(chan struct{}))

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	err = cluster.SetInfrastructureJSON(expectedInfra)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	var setInfra clusterInfra
	err = json.Unmarshal([]byte(cluster.InfrastructureJSON), &setInfra)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if !reflect.DeepEqual(setInfra, expectedInfra) {
		t.Fatalf(
			"expected infra to equal '%+v', got '%+v'",
			expectedInfra,
			setInfra,
		)
	}
}
