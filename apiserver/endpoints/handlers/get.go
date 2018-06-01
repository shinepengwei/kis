package handlers

import (
	"net/http"
	"fmt"
	"kis/apiserver/registry/rest"
	"encoding/json"
)

func ListResource(r rest.Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		result,_ := r.List(ctx)
		fmt.Println(result)
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		data,_ := json.Marshal(result)
		w.Write(data)

	}
}
