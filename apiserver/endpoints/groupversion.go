package endpoints

import (
	"kis/apiserver/registry/rest"
	"path"
	"github.com/emicklei/go-restful"
)

// APIGroupVersion is a helper for exposing rest.Storage objects as http.Handlers via go-restful
// It handles URLs of the form:
// /${storage_key}[/${object_name}]
// Where 'storage_key' points to a rest.Storage object stored in storage.
// This object should contain all parameterization necessary for running a particular API version
type APIGroupVersion struct {
	Storage map[string] rest.RestStorage

	Root string
	GroupVersion string
}

// InstallREST registers the REST handlers (storage, watch, proxy and redirect) into a restful Container.
// It is expected that the provided path root prefix will serve all operations. Root MUST NOT end
// in a slash.
func (g *APIGroupVersion) InstallREST(container *restful.Container) error {
	prefix := path.Join(g.Root, g.GroupVersion)
	installer := &APIInstaller{
		group: g,
		prefix:prefix,
	}
	ws := installer.Install()
	container.Add(ws)
	return nil
}
