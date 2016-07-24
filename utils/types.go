package utils

import "net/http"

//HttpClientDoer - interface for a http.Client.Doer
type HttpClientDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}
