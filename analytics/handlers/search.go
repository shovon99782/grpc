package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/example/analytics-service/internal/elastic"
)

func SearchOrders(w http.ResponseWriter, r *http.Request) {
	product := r.URL.Query().Get("product")
	customer := r.URL.Query().Get("customer")
	status := r.URL.Query().Get("status")

	var filter []map[string]interface{}

	if product != "" {
		filter = append(filter, map[string]interface{}{
			"match": map[string]interface{}{
				"items.sku": product,
			},
		})
	}

	if customer != "" {
		filter = append(filter, map[string]interface{}{
			"match": map[string]interface{}{
				"customer_name": customer,
			},
		})
	}

	if status != "" {
		filter = append(filter, map[string]interface{}{
			"match": map[string]interface{}{
				"status": status,
			},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": filter,
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
