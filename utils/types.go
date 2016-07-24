package utils

import (
	"net/http"

	"github.com/enaml-ops/enaml/enamlbosh"
)

//HttpClientDoer - interface for a http.Client.Doer
type HttpClientDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type BoshClientCaller interface {
	GetInfo() (*enamlbosh.BoshInfo, error)
}
