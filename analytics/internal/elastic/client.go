package elastic

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/example/analytics-service/config"
)

var ES *elasticsearch.Client

func InitElasticsearch() {
	cfg := config.LoadConfig()
	escfg := elasticsearch.Config{
		Addresses: []string{cfg.ElasticUrl},
	}

	client, err := elasticsearch.NewClient(escfg)
	if err != nil {
		log.Fatalf("❌ Failed to create Elasticsearch client: %v", err)
	}

	ES = client
	log.Println("✅ Elasticsearch connected")
}
