TAG=$(shell git describe --tags --always)
VERSION=$(TAG:v%=%)
REPO=jukeizu/weather
GO=GO111MODULE=on go
BUILD=GOARCH=amd64 $(GO) build -ldflags="-s -w -X main.Version=$(VERSION)" 
PROTOFILES=$(wildcard api/protobuf-spec/*/*.proto)
PBFILES=$(patsubst %.proto,%.pb.go, $(PROTOFILES))
BUILD_IMAGE=$(REPO):build
IMAGE=$(REPO):$(VERSION)

.PHONY: all deps test proto build clean $(PROTOFILES)

all: deps test build 
deps:
	$(GO) mod download

test:
	$(GO) vet ./...
	$(GO) test -v -race ./...

build:
	$(BUILD) -o bin/weather-$(VERSION) .

build-linux:
	CGO_ENABLED=0 GOOS=linux $(BUILD) -a -installsuffix cgo -o bin/weather .

docker-pull:
	docker pull $(BUILD_IMAGE) || true

docker-build:
	docker build --target build --cache-from $(BUILD_IMAGE) -t $(BUILD_IMAGE) .
	docker build --cache-from $(BUILD_IMAGE) -t $(IMAGE) .

docker-deploy:
	docker push $(IMAGE)
	docker push $(BUILD_IMAGE)

proto: $(PBFILES)

%.pb.go: %.proto
	cd $(dir $<) && protoc $(notdir $<) --go_out=plugins=grpc:.

clean:
	@find bin -type f ! -name '*.toml' -delete -print
