FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o soozan-ws main.go

# Final stage (small & secure runtime)
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/soozan-ws .
# Read the pub key path from env vars
COPY --from=builder /app/public_key.pem .

EXPOSE 8080

CMD ["./soozan-ws"]
