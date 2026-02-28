.PHONY: build build-broker build-trader install test generate docs clean lint ci

BROKER_APP_NAME = alpaca-broker
TRADER_APP_NAME = alpaca-trader
TARGET_PACKAGES = ./cmd/broker ./cmd/trader

build: build-broker build-trader

build-broker:
	go build -o $(BROKER_APP_NAME) cmd/broker/main/main.go

build-trader:
	go build -o $(TRADER_APP_NAME) cmd/trader/main/main.go

install: build
	@mkdir -p $$(go env GOPATH)/bin
	cp $(BROKER_APP_NAME) $$(go env GOPATH)/bin/
	cp $(TRADER_APP_NAME) $$(go env GOPATH)/bin/

test:
	go test $(TARGET_PACKAGES)

coverage:
	go test -v -coverprofile=coverage.out $(TARGET_PACKAGES)
	go tool cover -func=coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}' | awk '{if ($$1 < 60) {print "Coverage is below 60%"; exit 1} else {print "Coverage is good"}}'

generate:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -generate types,client -package brokerclient api/openapi.yaml > pkg/brokerclient/client.gen.go

docs:
	go run scripts/gen-docs.go

clean:
	rm -f $(BROKER_APP_NAME) $(TRADER_APP_NAME)

	$$(go env GOPATH)/bin/golangci-lint run ./...

ci: lint coverage test
