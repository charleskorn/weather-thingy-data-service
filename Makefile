all: test analyse

clean:
	go clean

setup:
	go get -t -v
	go get -v github.com/jteeuwen/go-bindata/...

generate:
	go-bindata -pkg main -o bindata.go db/migrations/

build: generate
	go build -o weather-thingy-data-service

test:
	go test

analyse:
	go tool vet -all -shadow .

docker-build:
	CGO_ENABLED=0 GOOS=linux go build -o weather-thingy-data-service-amd64-linux -a -installsuffix cgo .
	docker build --tag=charleskorn/weather-thingy-data-service .
