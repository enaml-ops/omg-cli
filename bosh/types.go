package bosh

import "github.com/enaml-ops/enaml/enamlbosh"

type BoshClientCaller interface {
	GetInfo() (*enamlbosh.BoshInfo, error)
}
