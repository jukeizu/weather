VERSION=$(shell git describe --tags)
BUILD=GOARCH=amd64 go build -v

.PHONY: all deps test proto build clean

all: deps test build

deps:
	go get -t -v ./...

test:
	go vet ./...
	go test -v -race ./...

proto:
	cd api/weather && protoc weather.proto --go_out=plugins=grpc:.
	cd api/geocoding && protoc geocoding.proto --go_out=plugins=grpc:.

build:
	for CMD in `ls cmd/services`; do $(BUILD) -o bin/$$CMD-service-$(VERSION) ./cmd/services/$$CMD; done
	# for CMD in `ls cmd/listeners`; do $(BUILD) -o bin/$$CMD-listener-$(VERSION) ./cmd/listeners/$$CMD; done

clean:
	find bin -type f ! -name '*.toml' -delete -print
