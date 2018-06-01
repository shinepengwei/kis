package generic

import (
	"kis/apiserver/runtime"
	"context"
	"kis/apiserver/storage"
	"kis/apiserver/runtime/schema"
	apis "kis/apis/core"
)

// Store implements pkg/api/rest.StandardStorage. It's intended to be
// embeddable and allows the consumer to implement any non-generic functions
// that are required. This object is intended to be copyable so that it can be
// used in different ways but share the same underlying behavior.
//
// All fields are required unless specified.
//
// The intended use of this type is embedding within a Kind specific
// RESTStorage implementation. This type provides CRUD semantics on a Kubelike
// resource, handling details like conflict detection with ResourceVersion and
// semantics. The RESTCreateStrategy, RESTUpdateStrategy, and
// RESTDeleteStrategy are generic across all backends, and encapsulate logic
// specific to the API.
//
// TODO: make the default exposed methods exactly match a generic RESTStorage
type Store struct {
	// NewFunc returns a new instance of the type this registry returns for a
	// GET of a single object, e.g.:
	//
	// curl GET /apis/group/version/namespaces/my-ns/myresource/name-of-object
	NewFunc func() runtime.Object

	// NewListFunc returns a new list of the type this registry; it is the
	// type returned when the resource is listed, e.g.:
	//
	// curl GET /apis/group/version/namespaces/my-ns/myresource
	NewListFunc func() runtime.Object

	KeyFunc func(ctx context.Context)(string)//返回存到etcd中的key

	// Storage is the interface for the underlying storage for the resource.
	Storage storage.Interface
}

// CompleteWithOptions updates the store with the provided options and
// defaults common fields.
func (e *Store) CompleteWithOptions(optsGetter RESTOptionsGetter) error {
	opts:= optsGetter.GetRESTOptions(schema.GroupResource{"v1","pods"})
	e.Storage,_ = opts.Decorator(opts.StorageConfig)
	return nil
}


// New implements RESTStorage.New.
func (e *Store) New() runtime.Object {
	return e.NewFunc()
}

// Create inserts a new item according to the unique key from the object.
func (e *Store) Create(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	out := e.NewFunc()
	e.Storage.Create(ctx, getKey(obj), obj, out)
	return out,nil
}


// NewList implements rest.Lister.
func (e *Store) NewList() runtime.Object {
	return e.NewListFunc()
}


// List returns a list of items matching labels and field according to the
// store's PredicateFunc.
func (e *Store) List(ctx context.Context) (runtime.Object, error) {
	out, _ := e.ListPredicate(ctx)
	return out, nil
}

// ListPredicate returns a list of all the items matching the given
// SelectionPredicate.
func (e *Store) ListPredicate(ctx context.Context) (runtime.Object, error) {
	list := e.NewListFunc()
	err := e.Storage.List(ctx, getKey(list), "v1",list)
	return list,err
}

func getKey(obj interface{}) string {
	switch t := obj.(type) {
	case *apis.Pod:
		return "PODS/"+t.Name
	case *apis.PodList:
		return "PODS/"
	}
	return ""
}