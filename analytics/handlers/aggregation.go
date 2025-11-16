package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/example/analytics-service/internal/elastic"
)

func OrdersByStatus(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"status_count": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "status.keyword",
				},
			},
		},
	}

	body, _ := json.Marshal(query)

	res, err := elastic.ES.Search(
		elastic.ES.Search.WithIndex("orders"),
		elastic.ES.Search.WithBody(bytes.NewReader(body)),
	)

	if err != nil {
		http.Error(w, "Elastic error", 500)
		return
	}

	defer res.Body.Close()

	io.Copy(w, res.Body)
}

func OrdersByCustomer(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"customer_count": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "customer_name.keyword",
				},
			},
		},
	}

	body, _ := json.Marshal(query)

	res, err := elastic.ES.Search(
		elastic.ES.Search.WithIndex("orders"),
		elastic.ES.Search.WithBody(bytes.NewReader(body)),
	)

	if err != nil {
		http.Error(w, "Elastic error", 500)
		return
	}
	defer res.Body.Close()

	io.Copy(w, res.Body)
}
