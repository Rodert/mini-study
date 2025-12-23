REGISTRY       ?= crpi-4otucz63tm2q5dhq.cn-beijing.personal.cr.aliyuncs.com
IMAGE_NAME     ?= project-shiyu/mini-study-backend
REGISTRY_USER  ?= shiyuwang7

# 当前 git commit 的前 8 位作为镜像版本
GIT_SHA        := $(shell git rev-parse --short=8 HEAD)
IMAGE_TAG      := $(GIT_SHA)
IMAGE_VERSION  := $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
IMAGE_LATEST   := $(REGISTRY)/$(IMAGE_NAME):latest

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

docker-login:
	docker login $(REGISTRY) --username=$(REGISTRY_USER)

docker-release: docker-login
	docker build -t $(IMAGE_VERSION) .
	docker tag $(IMAGE_VERSION) $(IMAGE_LATEST)
	docker push $(IMAGE_VERSION)
	docker push $(IMAGE_LATEST)
