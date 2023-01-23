package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const query = `
mutation importOrder($data: OrderImportInput!){
  importOrder(data:$data) {
    _id
    order {
      route {
        distance
        time
        routeDetails {
          tolls {
            summary {
              amount
            }
            byCountryAndTollSystem {
              name
              amount
            }
          }
        }
      }
    }
  }
}
`

func main() {
	data, err := ioutil.ReadFile("./data/simpleOrder.json")
	if err != nil {
		panic(err)
	}

	endpoint := os.Getenv("ENDPOINT")
	if endpoint == "" {
		endpoint = "https://backend.impargo.eu/"
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN not set")
	}

	reqBody := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"data": json.RawMessage(data),
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", token)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		panic(err)
	}

	if data, ok := respBody["data"]; ok {
		fmt.Println("Tolls details of order:\n", string(respBytes), data)
	} else {
		fmt.Println("Error:", respBody["errors"])
	}
}
