package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//	"k8s.io/kubernetes/pkg/api/v1"

	vaultapi "github.com/hashicorp/vault/api"
)

var (
	secretjob *Secret
)

type Secret struct {
	Data SecretToken `json:"data"`
}

type SecretToken struct {
	Token string `json:"token"`
}

type VaultClient struct {
	client *vaultapi.Client
}

func createVaultClient(tokenRef, vaultUrl string) (*VaultClient, error) {

	config := &vaultapi.Config{
		Address: vaultUrl,
	}

	// get secrets from token ref

	resp, err := http.Get(apiHost + fmt.Sprintf("/api/v1/namespaces/default/secrets/"+tokenRef))
	if err != nil {
		log.Println(red("Error pulling secrets", err))
	}
	secretObj := json.NewDecoder(resp.Body)
	err = secretObj.Decode(&secretjob)
	if err != nil {
		log.Println(err)
	}

	token := secretjob.Data.Token
	decodedToken, err := base64.StdEncoding.DecodeString(token)

	client, err := vaultapi.NewClient(config)
	client.SetToken(string(decodedToken))
	if err != nil {
		log.Println("Error creating vault", err)
		return nil, err
	}

	return &VaultClient{client}, nil

}

func (vaultClient *VaultClient) read(secretPath string) map[string]interface{} {
	c := vaultClient.client.Logical()

	secret, err := c.Read(secretPath)
	if err != nil {
		log.Println("VAULT Error getting secret", err)

	}
	return secret.Data
}
