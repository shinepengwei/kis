package generic

import (
	"kis/apiserver/storage/storagebackend"
	"kis/apiserver/storage"

	"kis/apiserver/storage/storagebackend/factory"
	"github.com/golang/glog"
)

type StorageDecorator func(config *storagebackend.Config) (storage.Interface, storagebackend.DestroyFunc)

func UndecoratedStorage(config *storagebackend.Config) (storage.Interface, storagebackend.DestroyFunc) {
	return NewRawStorage(config)
}

// NewRawStorage creates the low level kv storage. This is a work-around for current
// two layer of same storage interface.
// TODO: Once cacher is enabled on all registries (event registry is special), we will remove this method.
func NewRawStorage(config *storagebackend.Config) (storage.Interface, storagebackend.DestroyFunc) {
	s, d, err := factory.CreateStorageBackend(config)
	if err != nil {
		glog.Fatalf("Unable to create storage backend: config (%v), err (%v)", config, err)
	}
	return s, d
}