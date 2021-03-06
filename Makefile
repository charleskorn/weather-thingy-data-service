COMMIT := $(shell git rev-parse --verify HEAD --short)
DOCKER_REPO := charleskorn/weather-thingy-data-service
DOCKER_IMAGE := $(DOCKER_REPO):git-$(COMMIT)

all: test analyse

clean:
	go clean

setup:
	go get -t -v
	go get -v github.com/jteeuwen/go-bindata/...
	go get -v github.com/golang/mock/mockgen

generate:
	go-bindata -pkg main -o bindata.go db/migrations/
	mockgen -package=main -destination=mock_render.go "github.com/martini-contrib/render" Render
	mockgen -package=main -source=database.go -destination=mock_database.go

build: generate
	go build -o weather-thingy-data-service

test:
	go test

analyse:
	go tool vet -all -shadow .

docker-build:
	CGO_ENABLED=0 GOOS=linux go build -o weather-thingy-data-service-amd64-linux -a -installsuffix cgo .
	docker build --tag=$(DOCKER_IMAGE) .

docker-tag-travis:
ifeq "$(TRAVIS_BUILD_NUMBER)" ""
	$(error TRAVIS_BUILD_NUMBER environment variable not defined)
endif

ifeq "$(TRAVIS_PULL_REQUEST)" ""
	$(error TRAVIS_PULL_REQUEST environment variable not defined)
endif

ifeq "$(TRAVIS_BRANCH)" ""
	$(error TRAVIS_BRANCH environment variable not defined)
endif

	docker tag $(DOCKER_IMAGE) $(DOCKER_REPO):travis-$(TRAVIS_BUILD_NUMBER)

ifeq "$(TRAVIS_PULL_REQUEST)" "false"

	docker tag $(DOCKER_IMAGE) $(DOCKER_REPO):$(TRAVIS_BRANCH)

ifeq "$(TRAVIS_BRANCH)" "master"
	docker tag $(DOCKER_IMAGE) $(DOCKER_REPO):latest
endif
endif

docker-push:
	docker push $(DOCKER_REPO)
