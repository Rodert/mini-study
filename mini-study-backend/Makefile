run:
	air

build:
	go build -o bin/app ./cmd/server

test:
	go test ./... -cover

swagger:
	swag init -g cmd/server/main.go -o internal/docs

docker:
	docker build -t mini-study-backend .

compose:
	docker-compose up -d
