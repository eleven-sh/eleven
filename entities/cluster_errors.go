package entities

type ErrClusterAlreadyExists struct {
	ClusterName string
}

func (e ErrClusterAlreadyExists) Error() string {
	return "ErrClusterAlreadyExists"
}

type ErrClusterNotExists struct {
	ClusterName string
}

func (e ErrClusterNotExists) Error() string {
	return "ErrClusterNotExists"
}
