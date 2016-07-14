package pluginutil

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/xchapter7x/lo"
)

type VaultUnmarshaler interface {
	UnmarshalFlags(hash string, flgs []cli.Flag) (err error)
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func NewVaultUnmarshal(domain, token string, client Doer) VaultUnmarshaler {
	return &VaultUnmarshal{
		Domain: domain,
		Token:  token,
		Client: client,
	}
}

type VaultUnmarshal struct {
	Domain string
	Token  string
	Client Doer
}

type VaultJsonObject struct {
	LeaseID       string            `json:"lease_id"`
	Renewable     bool              `json:"renewable"`
	LeaseDuration float64           `json:"lease_duration"`
	Data          map[string]string `json:"data"`
	Warnings      interface{}       `json:"warnings"`
	Auth          interface{}       `json:"auth"`
}

func (s *VaultUnmarshal) UnmarshalFlags(hash string, flgs []cli.Flag) (err error) {
	b := s.getVaultHashValues(hash)
	vaultObj := new(VaultJsonObject)
	json.Unmarshal(b, vaultObj)

	for h, v := range vaultObj.Data {

		for i, f := range flgs {
			if h == f.GetName() {

				switch ft := f.(type) {
				case cli.StringSliceFlag:
					ft.EnvVar = v
					var slice cli.StringSlice = strings.Split(v, ",")
					ft.Value = &slice
					flgs[i] = ft
					fmt.Println("vl: ", ft.Value.Value())

				case cli.StringFlag:
					ft.EnvVar = v
					ft.Value = v
					flgs[i] = ft
					fmt.Println("vl: ", ft.Value)

				default:
					lo.G.Panic("i dont know what field this is in VAULT. exiting the app now")
				}
			}
		}
	}
	return
}

func (s *VaultUnmarshal) getVaultHashValues(hash string) []byte {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/%s", s.Domain, hash), nil)
	s.decorateWithToken(req)
	res, _ := s.Client.Do(req)
	b, _ := ioutil.ReadAll(res.Body)
	return b
}

func (s *VaultUnmarshal) decorateWithToken(req *http.Request) {
	req.Header.Add("X-Vault-Token", s.Token)
}

func DefaultClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}
