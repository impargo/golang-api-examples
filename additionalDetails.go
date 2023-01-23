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
query GetCompanyOrders($query: CompanyOrderQuery, $paginate: Pagination, $sortBy: CompanyOrderSortBy) {
  companyOrders(query: $query, paginate: $paginate, sortBy: $sortBy) {
    items {
      ...CompanyOrderOverview
      __typename
    }
    hasNext
    __typename
  }
}

fragment CompanyOrderOverview on CompanyOrder {
  _id
  archived
  vehicleDetails {
    _id
    name
    __typename
  }
  subcontractorDetails {
    ... on CostprofilesSubcontractor {
      _id
      name
      email
      __typename
    }
    __typename
  }
  driverDetails {
    _id
    name
    __typename
  }
  customerDetails {
    _id
    name
    email
    __typename
  }
  quote {
    quote {
      costsPerKm
      total
      toll
      route {
        distance
        __typename
      }
      __typename
    }
    __typename
  }
  quotexSettings {
    pricePerKm
    __typename
  }
  order {
    ...OrderOverview
    __typename
  }
  __typename
}

fragment OrderOverview on Order {
  _id
  load {
    length
    weight
    bodyType
    price
    __typename
  }
  reference
  status
  additionalDetails {
    value
    label
    customFieldId
  }
  route {
    distance
    time
    __typename
  }
  ...OrderFirstLastStopData
  __typename
}

fragment OrderFirstLastStopData on Order {
  firstStop {
    location {
      country
      city
      zipcode
      coordinates {
        lat
        lon
        __typename
      }
      __typename
    }
    date {
      from
      to
      __typename
    }
    time {
      from
      to
      __typename
    }
    __typename
  }
  lastStop {
    location {
      country
      city
      zipcode
      coordinates {
        lat
        lon
        __typename
      }
      __typename
    }
    date {
      from
      to
      __typename
    }
    time {
      from
      to
      __typename
    }
    __typename
  }
  __typename
}
`

func main() {
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
			"paginate": map[string]interface{}{
				"limit": 20,
			},
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
		fmt.Println("Successfully read orders:\n", string(respBytes), data)
	} else {
		fmt.Println("Error:", respBody["errors"])
	}
}
