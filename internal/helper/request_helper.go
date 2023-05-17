package helper

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

func RequestClone(req *http.Request) *http.Request {
	r2 := req.Clone(context.Background())
	if r2.Body != nil {
		bts, _ := ioutil.ReadAll(r2.Body)
		req.Body = ioutil.NopCloser(bytes.NewReader(bts))
		r2.Body = ioutil.NopCloser(bytes.NewReader(bts))
	}
	return r2
}
