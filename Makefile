APP=echo
APP_VERSION=0.1
APP_COMMIT=$(shell git rev-parse HEAD)
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

LOCAL_CONFIG_FILE=local.env
DOCKER_REGISTRY_USER_NAME=nsnikhil
HTTP_SERVE_COMMAND=http-serve

deps:
	go mod download

tidy:
	go mod tidy

check: fmt vet lint

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	golint $(ALL_PACKAGES)

compile:
	mkdir -p out/
	go build -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(APP_COMMIT)" -o $(APP_EXECUTABLE) *.go

build: deps compile

http-serve: build
	$(APP_EXECUTABLE) $(HTTP_SERVE_COMMAND)

docker-build:
	docker build -t $(DOCKER_REGISTRY_USER_NAME)/$(APP):$(APP_VERSION) .
	docker rmi -f $$(docker images -f "dangling=true" -q) || true

docker-push: docker-build
	docker push $(DOCKER_REGISTRY_USER_NAME)/$(APP):$(APP_VERSION)

clean:
	rm -rf out/

test:
	go clean -testcache
	go test ./...