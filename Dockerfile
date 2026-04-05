FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o booking-service ./cmd/Server/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o seed ./cmd/Script/seed.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/booking-service .
COPY --from=builder /app/seed .
COPY --from=builder /app/migrations ./migrations

RUN chmod +x /app/booking-service /app/seed

EXPOSE 8080
CMD ["./booking-service"]