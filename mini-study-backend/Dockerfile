# syntax=docker/dockerfile:1
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mini-study ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/mini-study ./mini-study
COPY configs configs
COPY storage storage
EXPOSE 8080
ENTRYPOINT ["/app/mini-study"]
