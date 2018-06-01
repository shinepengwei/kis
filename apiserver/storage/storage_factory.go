package storage

import (
	"kis/apiserver/storage/storagebackend"
	"kis/apiserver/runtime/schema"
)

type Backend struct {
	// the url of storage backend like: https://etcd.domain:2379
	Server string
}
type StorageFactory interface {
	// New finds the storage destination for the given group and resource. It will
	// return an error if the group has no storage destination configured.
	NewConfig(groupResource schema.GroupResource) (*storagebackend.Config, error)

	// ResourcePrefix returns the overridden resource prefix for the GroupResource
	// This allows for cohabitation of resources with different native types and provides
	// centralized control over the shape of etcd directories
	ResourcePrefix(groupResource schema.GroupResource) string

	// Backends gets all backends for all registered storage destinations.
	// Used for getting all instances for health validations.
	Backends() []Backend
}

type DefaultStorageFactory struct {
	StorageConfig storagebackend.Config
	DefaultResourcePrefixes map[schema.GroupResource]string

	// DefaultMediaType is the media type used to store resources. If it is not set, "application/json" is used.
	DefaultMediaType string
}

func (s *DefaultStorageFactory)NewConfig(groupResource schema.GroupResource) (*storagebackend.Config, error){

	return &storagebackend.Config{
		Type:"etcd3",
		Prefix:"/",
		ServerList: []string{"http://localhost:2379"},
	},nil
}
func (s *DefaultStorageFactory) ResourcePrefix(groupResource schema.GroupResource) string{
	return ""
}
func (s *DefaultStorageFactory) Backends() []Backend{
	backends := []Backend{}
	return backends
}

func NewStorageFactory(config storagebackend.Config,
	defaultMediaType string,
) *DefaultStorageFactory{
	return &DefaultStorageFactory{
		StorageConfig:           config,
		DefaultMediaType:        defaultMediaType,
	}
}