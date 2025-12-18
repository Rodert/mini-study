# syntax=docker/dockerfile:1
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mini-study ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./scripts/migrate.go

FROM golang:1.23
WORKDIR /app
COPY --from=builder /app/mini-study ./mini-study
COPY --from=builder /app/migrate ./migrate
COPY configs configs
COPY storage storage
EXPOSE 8080
ENTRYPOINT ["/app/mini-study"]
