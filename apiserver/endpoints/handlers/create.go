package handlers

import (
	"net/http"
	"kis/apiserver/registry/rest"
	"k8s.io/kubernetes/src/github.com/json-iterator/go"
	"fmt"
)

func CreateResource(r rest.Creater) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		obj := r.New()
		body, _ := readBody(req)
		fmt.Println("request create body:"+string(body[:]))
		//fmt.Println(jsoniter.Unmarshal(body))
		//json.Unmarshal(body, obj)
		jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body,obj)
		ctx := req.Context()
		r.Create(ctx,obj)
		w.Header().Set("Content-Type","application/text")
		w.WriteHeader(http.StatusOK)
		out := []byte("success")
		w.Write(out)
	}
}