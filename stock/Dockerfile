FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /stock-service ./cmd
EXPOSE 50052
CMD ["/stock-service"]
