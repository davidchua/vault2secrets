package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	newNamespace                   = &namespace
	apiHost                        = "http://127.0.0.1:8001"
	customSecretsWithWatchEndpoint = fmt.Sprintf("/apis/cubiclerebels.com/v1/namespaces/%v/customsecrets?watch=true", *newNamespace)
	customSecretsEndpoint          = fmt.Sprintf("/apis/cubiclerebels.com/v1/namespaces/%v/customsecrets", *newNamespace)
)

func pollSecrets() <-chan VaultEvent {
	log.Println(*newNamespace)
	log.Println(fmt.Sprintf("/apis/%v/", *newNamespace))
	events := make(chan VaultEvent)
	log.Println(customSecretsWithWatchEndpoint)

	go func() {
		for {
			log.Println("Polling")
			// Add CallKubernetes API endpoint here
			resp, err := http.Get(apiHost + customSecretsWithWatchEndpoint)
			if err != nil {
				log.Fatal(err)
			}

			decoder := json.NewDecoder(resp.Body)
			for {
				var vaultevent VaultEvent
				err = decoder.Decode(&vaultevent)
				if err != nil {
					log.Println(customSecretsWithWatchEndpoint)
					log.Fatal(err)
				}
				events <- vaultevent

			}
			time.Sleep(5 * time.Second)
		}

	}()
	return events

}
