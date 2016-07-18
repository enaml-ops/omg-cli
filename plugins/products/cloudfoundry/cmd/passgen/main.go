package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func RandomString(strlen int) string {
	const chars = "abcdefghipqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

const passLength = 20

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	passVault := PasswordVault{
		Router:             RandomString(passLength),
		Nats:               RandomString(passLength),
		MysqlAdmin:         RandomString(passLength),
		MysqlBootstrap:     RandomString(passLength),
		MysqlProxyAPI:      RandomString(passLength),
		CCStagingUpload:    RandomString(passLength),
		CCBulkAPI:          RandomString(passLength),
		CCInternalAPI:      RandomString(passLength),
		DBUaa:              RandomString(passLength),
		DBCCDB:             RandomString(passLength),
		DBConsole:          RandomString(passLength),
		DiegoDB:            RandomString(passLength),
		UAALdapUser:        RandomString(passLength),
		Admin:              RandomString(passLength),
		PushAppsManager:    RandomString(passLength),
		SmokeTests:         RandomString(passLength),
		SystemServices:     RandomString(passLength),
		SystemVerification: RandomString(passLength),
		SystemPasswords:    RandomString(passLength),
	}
	b, _ := json.Marshal(passVault)
	fmt.Println(string(b))
}

type PasswordVault struct {
	Router             string `json:"router-pass"`
	Nats               string `json:"nats-pass"`
	MysqlAdmin         string `json:"mysql-admin-password"`
	MysqlBootstrap     string `json:"mysql-bootstrap-password"`
	MysqlProxyAPI      string `json:"mysql-proxy-api-password"`
	CCStagingUpload    string `json:"cc-staging-upload-password"`
	CCBulkAPI          string `json:"cc-bulk-api-password"`
	CCInternalAPI      string `json:"cc-internal-api-password"`
	DBUaa              string `json:"db-uaa-password"`
	DBCCDB             string `json:"db-ccdb-password"`
	DBConsole          string `json:"db-console-password"`
	DiegoDB            string `json:"diego-db-passphrase"`
	UAALdapUser        string `json:"uaa-ldap-user-password"`
	Admin              string `json:"admin-password"`
	PushAppsManager    string `json:"push-apps-manager-password"`
	SmokeTests         string `json:"smoke-tests-password"`
	SystemServices     string `json:"system-services-password"`
	SystemVerification string `json:"system-verification-password"`
	SystemPasswords    string `json:"system-passwords-client-secret"`
}
