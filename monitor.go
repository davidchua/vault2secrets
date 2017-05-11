package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	namespace                      = os.Getenv("NAMESPACE")
	apiHost                        = "http://127.0.0.1:8001"
	customSecretsWithWatchEndpoint = fmt.Sprintf("/apis/cubiclerebels.com/v1/namespaces/%s/customsecrets?watch=true", namespace)
	customSecretsEndpoint          = fmt.Sprintf("/apis/cubiclerebels.com/v1/namespaces/%s/customsecrets", namespace)
)

func pollSecrets() <-chan VaultEvent {
	events := make(chan VaultEvent)

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
