FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /order-service ./cmd
EXPOSE 50051
CMD ["/order-service"]
