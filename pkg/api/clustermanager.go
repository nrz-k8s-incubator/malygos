package api

type ClusterManager interface {
	Create(cluster *Cluster) (*Cluster, error)
	Delete(id string) error
	List() ([]*Cluster, error)
	Get(id string) (*Cluster, error)
	ListSubscriptions(id string) ([]*CatalogComponent, error)
}
