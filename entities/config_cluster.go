package entities

import "errors"

func (c *Config) SetCluster(cluster *Cluster) error {
	if cluster == nil {
		return errors.New("passed cluster is nil")
	}

	c.Clusters[cluster.Name] = cluster

	return nil
}

func (c *Config) ClusterExists(clusterName string) bool {
	_, clusterExists := c.Clusters[clusterName]
	return clusterExists
}

func (c *Config) GetCluster(clusterName string) (*Cluster, error) {
	if !c.ClusterExists(clusterName) {
		return nil, ErrClusterNotExists{
			ClusterName: clusterName,
		}
	}

	cluster := c.Clusters[clusterName]

	return cluster, nil
}

func (c *Config) RemoveCluster(clusterName string) error {
	if !c.ClusterExists(clusterName) {
		return ErrClusterNotExists{
			ClusterName: clusterName,
		}
	}

	delete(c.Clusters, clusterName)

	return nil
}
