# Order Tracking Microservices

This repository contains a sample microservices application for order tracking. The application is composed of three services: `order-service`, `stock-service`, and `analytics-service`.

## Project Overview

The order tracking microservices application is designed to handle order creation, stock management, and data analytics. It uses a combination of technologies including Go, gRPC, RabbitMQ, and Elasticsearch to provide a scalable and resilient system.

### Services

*   **order-service**: This service is responsible for creating and managing orders. It communicates with the `stock-service` to check for product availability and publishes order creation events to RabbitMQ.
*   **stock-service**: This service manages the product inventory. It listens for order creation events from RabbitMQ and updates the stock accordingly.
*   **analytics-service**: This service consumes order creation events from RabbitMQ, indexes the data in Elasticsearch, and provides an API for querying the aggregated data.

## Getting Started

To get started with this project, you will need to have Docker and Docker Compose installed on your machine.

### Prerequisites

*   [Docker](https://docs.docker.com/get-docker/)
*   [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/HETIC-MT-P2026/order-tracking-microservices.git
    ```
2.  Navigate to the project directory:
    ```sh
    cd order-tracking-microservices
    ```
3.  Start the services using Docker Compose:
    ```sh
    docker-compose up -d
    ```

This will build the Docker images for each service and start the containers in detached mode.

## Usage

Once the services are up and running, you can interact with them using a gRPC client or by sending HTTP requests to the `analytics-service`.

### Creating an Order

To create an order, you can use the `dummy_client.go` file in the `order` service directory. This will create a dummy order and send it to the `order-service`.

### Querying Analytics Data

The `analytics-service` provides an API for querying the aggregated order data from Elasticsearch. You can send a GET request to `http://localhost:8080/search?q=<query>` to search for orders.

For example, to search for orders with the product name "Laptop", you can use the following URL:

```
http://localhost:8080/search?q=product_name:Laptop
```

This will return a JSON response with the search results.