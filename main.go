package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

var (
	api         = "http://127.0.0.1:8001"
	namespace   = os.Getenv("NAMESPACE")
	vaultClient *VaultClient
	yellow      = color.New(color.FgYellow).SprintFunc()
	green       = color.New(color.FgGreen).SprintFunc()
	red         = color.New(color.FgRed).SprintFunc()
)

func main() {

	if namespace == "" {
		namespace = "default"
		log.Fatal("No NAMESPACE Environment Variable set. Please set the variable and re-run this binary")
	}

	log.Println(green("Vault2Secrets initialized"))

	//okChan := make(chan string)
	var wg sync.WaitGroup
	go func() {
		http.ListenAndServe("0.0.0.0:8888", nil)
	}()

	done := make(chan struct{})

	log.Println(yellow("Waiting ..."))

	wg.Add(2)
	go func() {
		// go get customSecrets
		secretEvent := pollSecrets()
		for {
			select {
			case event := <-secretEvent:
				// process event
				processEvent(event)
			case <-done:
				wg.Done()
				log.Println("Done")
			}

		}
		log.Println("Processing ...")
		wg.Done()
		log.Println("Done ...")
	}()

	go func() {
		// sync and reconsile
		for {
			log.Println("Syncing and Checking CustomSecrets")
			syncSecrets()
			time.Sleep(5 * time.Second)
		}

		wg.Done()

	}()

	wg.Wait()

	log.Println("Everything Done")

	//	abc := <-okChan

	//	fmt.Println(abc)

}
