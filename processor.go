package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

type FullSecret struct {
	Kind       string            `json:"kind"`
	APIVersion string            `json:"apiVersion"`
	Metadata   map[string]string `json:"metadata"`
	Data       map[string]string `json:"data"`
	Type       string            `json:"type"`
}

var secretsEndpoint = fmt.Sprintf("/api/v1/namespaces/%s/secrets", namespace)

func processEvent(event VaultEvent) {

	log.Println("Processing VaultEvent")
	vaultObj := event.Object

	//		vaultObj := VaultObject{event.Object.(VaultObject)}
	token := vaultObj.Spec.TokenRef
	url := vaultObj.Spec.Url
	path := vaultObj.Spec.Path

	if event.Type == "ADDED" {
		log.Println(yellow("Getting Token ", token, " URL ", url, " Path ", path))
		vClient, _ := createVaultClient(token, url)
		secret, err := vClient.read(path)
		if err != nil {
			return

		}
		processSecrets(vaultObj.Spec.Secret, secret)

	} else if event.Type == "DELETED" {
		// delete Secret
		ok, err := deleteSecret(vaultObj.Spec.Secret)
		if ok {
			log.Println(red("["+vaultObj.Spec.Secret, "] Secret Deleted"))

		}

		if err != nil {
			log.Println(err)

		}

	}

}

type customSecretList struct {
	Items []VaultObject `json:"items"`
}

var (
	allCustomSecrets *customSecretList
)

func deleteSecret(secretName string) (bool, error) {
	var b []byte
	secretContent := bytes.NewBuffer(b)
	req, err := http.NewRequest("DELETE", apiHost+secretsEndpoint+"/"+secretName, secretContent)
	if err != nil {
		log.Println(red(err))
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		return false, nil
	}
	if err != nil {
		return false, err

	}
	return true, nil
}

func syncSecrets() {
	// iterate through all customSecrets and Secret to ensure they are similar
	// open up CustomSecret
	resp, err := http.Get(apiHost + customSecretsEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&allCustomSecrets)
	//	log.Println(allCustomSecrets.Items)
	for _, item := range allCustomSecrets.Items {
		token := item.Spec.TokenRef
		url := item.Spec.Url
		path := item.Spec.Path
		vClient, err := createVaultClient(token, url)
		if err != nil {
			log.Println(err)
			break

		}
		secret, err := vClient.read(path)
		if err != nil {
			break

		}
		err = processSecrets(item.Spec.Secret, secret)
		if err != nil {
			log.Fatal("Error!")
		}
	}

}

func processSecrets(secret string, data map[string]interface{}) error {
	pdata := make(map[string]string)
	for k, v := range data {
		pdata[k] = base64.StdEncoding.EncodeToString([]byte(v.(string)))

	}
	metadata := make(map[string]string)
	metadata["name"] = secret

	// pull secrets
	var secretTmp = &FullSecret{
		APIVersion: "v1",
		Data:       pdata,
		Kind:       "Secret",
		Metadata:   metadata,
		Type:       "Opaque",
	}

	var currentSecret *FullSecret

	resp, err := http.Get(apiHost + secretsEndpoint + "/" + secret)
	if err != nil {
		return err
	}

	// compare secrets with data
	if resp.StatusCode == 200 {
		// replace Secret
		content := json.NewDecoder(resp.Body)
		content.Decode(&currentSecret)

		if !reflect.DeepEqual(secretTmp.Data, currentSecret.Data) {
			log.Println(green("["+secret, "] Vault Secret has changed."))
			modifySecret(secretTmp)

		} else {
			log.Println("Vault Secret has not changed. Not replacing.")
		}
	} else if resp.StatusCode == 404 {
		log.Println(red("["+secret, "] Vault Secret does not exist"))
		writeSecret(secretTmp)

	}
	// if secrets different, update secret with new data
	return nil
}
func writeSecret(secret *FullSecret) {
	var b []byte
	secretContent := bytes.NewBuffer(b)
	err := json.NewEncoder(secretContent).Encode(secret)
	req, err := http.NewRequest("POST", apiHost+secretsEndpoint, secretContent)
	if err != nil {
		log.Println(err)

	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil || resp.StatusCode != 201 {
		log.Println("Unable to write secret.")
	} else {
		log.Println(green("Secret written."))

	}

}
func modifySecret(secret *FullSecret) {
	var b []byte
	secretContent := bytes.NewBuffer(b)
	err := json.NewEncoder(secretContent).Encode(secret)
	req, err := http.NewRequest("PUT", apiHost+secretsEndpoint+"/"+secret.Metadata["name"], secretContent)
	if err != nil {
		log.Println(err)

	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil || resp.StatusCode != 200 {
		log.Println("Unable to modify secret.")

	} else {
		log.Println(green("Secret modified"))

	}
	//output, _ := ioutil.ReadAll(resp.Body)

}
