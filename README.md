# Order Tracking & Analytics Microservices (Go)
Skeleton project for Order Service, Stock Service, and Analytics Service.
Each service is a minimal Go module with gRPC proto, basic main, Dockerfile and placeholders.

curl "http://localhost:8080/search?product=sku123"
curl "http://localhost:8080/search?customer=John"
curl "http://localhost:8080/search?status=CREATED"


curl "http://localhost:8080/agg/status"
curl "http://localhost:8080/agg/customer"

