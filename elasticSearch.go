package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const (
	elasticSearchIndexEnvVar = "ELASTICSEARCH_INDEX"
)

// ElasticSearch Errors
var (
	ErrElasticSearchWriteFail = errors.New("elasticsearch write failed")
)

func openElasticSearch() *elasticsearch.Client {
	var (
		r map[string]interface{}
	)

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating ElasticSearch client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	return es
}

func writeEsPayload(elasticClient *elasticsearch.Client, deliveryID *string, payload []byte) error {
	index := os.Getenv(elasticSearchIndexEnvVar)
	esReq := esapi.IndexRequest{
		Index:      index,
		DocumentID: *deliveryID,

		Body: bytes.NewReader(payload),

		Refresh: "true",
	}

	res, err := esReq.Do(context.Background(), elasticClient)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d\n", res.Status(), deliveryID)

		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("%s\n", err)
		}
		bodyString := string(bodyBytes)
		log.Printf("%s\n", bodyString)

		return ErrElasticSearchWriteFail
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		log.Printf("Error parsing the response body: %s\n", err)
		return ErrElasticSearchWriteFail
	}

	log.Printf("[%s] %s; version=%d\n", res.Status(), responseData["result"], int(responseData["_version"].(float64)))

	return nil
}
