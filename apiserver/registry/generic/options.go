package generic

import (
	"kis/apiserver/runtime/schema"
	"time"
	"kis/apiserver/storage/storagebackend"
)

// RESTOptions is set of configuration options to generic registries.
type RESTOptions struct {
	StorageConfig *storagebackend.Config
	Decorator     StorageDecorator

	EnableGarbageCollection bool
	DeleteCollectionWorkers int
	ResourcePrefix          string
	CountMetricPollPeriod   time.Duration
}

type RESTOptionsGetter interface {
	GetRESTOptions(resource schema.GroupResource) (RESTOptions)
}
