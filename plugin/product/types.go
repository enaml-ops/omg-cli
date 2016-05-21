package product

import (
	"log"
	"net/rpc"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/go-plugin"
	"github.com/enaml-ops/enaml"
)

type Meta struct {
	Name       string
	Properties map[string]interface{}
}

// ProductDeployer is the interface that we will expose for cloud config
// plugins
type ProductDeployer interface {
	GetMeta() Meta
	GetFlags() []cli.Flag
	GetProduct(args []string) enaml.DeploymentManifest
}

// ProductRPC - Here is an implementation that talks over RPC
type ProductRPC struct{ client *rpc.Client }

func (s *ProductRPC) GetMeta() Meta {
	var resp Meta
	err := s.client.Call("Plugin.GetMeta", new(interface{}), &resp)

	if err != nil {
		panic(err)
	}
	return resp
}

func (s *ProductRPC) GetProduct(args []string) enaml.DeploymentManifest {
	var resp enaml.DeploymentManifest
	log.Println("calling rpc client getcloudconfig")
	err := s.client.Call("Plugin.GetProduct", args, &resp)
	log.Println("call failed:", err)
	if err != nil {
		panic(err)
	}
	return resp
}

func (s *ProductRPC) GetFlags() []cli.Flag {
	var resp []cli.Flag
	err := s.client.Call("Plugin.GetFlags", new(interface{}), &resp)

	if err != nil {
		panic(err)
	}
	return resp
}

//ProductRPCServer - Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type ProductRPCServer struct {
	Impl ProductDeployer
}

func (s *ProductRPCServer) GetFlags(args interface{}, resp *[]cli.Flag) error {
	*resp = s.Impl.GetFlags()
	return nil
}

func (s *ProductRPCServer) GetMeta(args interface{}, resp *Meta) error {
	*resp = s.Impl.GetMeta()
	return nil
}

func (s *ProductRPCServer) GetProduct(args []string, resp *enaml.DeploymentManifest) error {
	*resp = s.Impl.GetProduct(args)
	return nil
}

func NewProductPlugin(plg ProductDeployer) ProductPlugin {
	return ProductPlugin{
		Plugin: plg,
	}
}

type ProductPlugin struct {
	Plugin ProductDeployer
}

func (s ProductPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ProductRPCServer{Impl: s.Plugin}, nil
}

func (s ProductPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ProductRPC{client: c}, nil
}
