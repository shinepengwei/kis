package endpoints

import (
	"github.com/emicklei/go-restful"
	"sort"
	"kis/apiserver/registry/rest"
	"fmt"
	"strings"
	"kis/apiserver/endpoints/handlers"
)

type APIInstaller struct {
	group                        *APIGroupVersion
	prefix                       string // Path prefix where API resources are to be registered.
	//minRequestTimeout            time.Duration
	//enableAPIResponseCompression bool
}

// Struct capturing information about an action ("GET", "POST", "WATCH", "PROXY", etc).
type action struct {
	Verb          string               // Verb identifying the action ("GET", "POST", "WATCH", "PROXY", etc).
	Path          string               // The path of the action
	Params        []*restful.Parameter // List of parameters associated with the action.
	//Namer         handlers.ScopeNamer
	AllNamespaces bool // true iff the action is namespaced but works on aggregate result for all namespaces
}


func (a *APIInstaller) Install() *restful.WebService{
	ws := a.newWebService()
	paths := make([]string, len(a.group.Storage))
	var i int = 0
	for path := range a.group.Storage{
		paths[i] = path
		i++
	}
	sort.Strings(paths)
	for _, path := range paths{
		a.registerResourceHandlers(path, a.group.Storage[path], ws)
	}
	return ws
}

func(a *APIInstaller) newWebService() *restful.WebService{
	ws := new(restful.WebService)
	ws.Path(a.prefix)
	ws.Doc("API at" + a.prefix)
	return ws
}

func (a *APIInstaller) registerResourceHandlers(path string, storage rest.RestStorage, ws *restful.WebService) {
	resource, _, _ := splitSubresource(path)
	//namespaceParamName := "namespaces"
	namespaceParam := ws.PathParameter("namespace", "object name and auth scope, such as for teams and projects").DataType("string")
	namespacedPath := resource// + "/{" + "namespace" + "}/" + resource
	namespaceParams := []*restful.Parameter{namespaceParam}
	resourcePath := namespacedPath
	resourceParams := namespaceParams
	//nameParam := ws.PathParameter("name", "name of the "+kind).DataType("string")
	//pathParam := ws.PathParameter("path", "path to the resource").DataType("string")

	//itemPath := namespacedPath + "/{name}"
	//nameParams := append(namespaceParams, nameParam)
	//proxyParams := append(nameParams, pathParam)


	creater, isCreater := storage.(rest.Creater)//creater
	_, isGetter := storage.(rest.Getter)//getter
	lister, isLister := storage.(rest.Lister)
	actions := []action{}
	if isCreater{
		actions = append(actions, action{"POST", resourcePath, resourceParams, false})
	}
	if isGetter {
		actions = append(actions, action{"GET", resourcePath, resourceParams, false})
	}
	if isLister{
		actions = append(actions, action{"LIST", resourcePath, resourceParams, false})
	}

	routes := []*restful.RouteBuilder{}

	for _, action := range actions{
		switch action.Verb {
		case "LIST": // List all resources of a kind.
			handler := restfulListResource(lister)
			route := ws.GET(action.Path).To(handler)
			routes = append(routes, route)
		case "POST": // Create a resource.
			handler := restfulCreateResource(creater)
			route := ws.POST(action.Path).To(handler)
			routes = append(routes, route)
		}

	}

	for _, route := range routes{
		ws.Route(route)
	}

}

func restfulCreateResource(r rest.Creater) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.CreateResource(r)(res.ResponseWriter, req.Request)
	}
}

func restfulListResource(r rest.Lister) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.ListResource(r)(res.ResponseWriter, req.Request)
	}
}

func splitSubresource(path string) (string, string, error) {
	var resource, subresource string
	switch parts := strings.Split(path, "/"); len(parts) {
	case 2:
		resource, subresource = parts[0], parts[1]
	case 1:
		resource = parts[0]
	default:
		// TODO: support deeper paths
		return "", "", fmt.Errorf("api_installer allows only one or two segment paths (resource or resource/subresource)")
	}
	return resource, subresource, nil
}