package main

import (
	"log"

	"net/http"
)

const (
	path = "/webhook"
)

func main() {
	elasticClient := openElasticSearch()

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		deliveryID, githubPayload, err := receiveGithubWebhook(r)

		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		log.Printf("Received event %s\n", *deliveryID)
		err = writeEsPayload(elasticClient, deliveryID, githubPayload)

		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		w.WriteHeader(204)
	})

	http.ListenAndServe(":3000", nil)
}
