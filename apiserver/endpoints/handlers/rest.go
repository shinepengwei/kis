package handlers

import (
	"net/http"
	"io/ioutil"
)

func readBody(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	return ioutil.ReadAll(req.Body)
}
