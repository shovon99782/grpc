FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /analytics-service ./cmd
EXPOSE 8080
CMD ["/analytics-service"]
