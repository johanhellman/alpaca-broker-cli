.PHONY: build build-broker build-trader install test generate clean

BROKER_APP_NAME = alpaca-broker
TRADER_APP_NAME = alpaca-trader

build: build-broker build-trader

build-broker:
	go build -o $(BROKER_APP_NAME) cmd/broker/main/main.go

build-trader:
	go build -o $(TRADER_APP_NAME) cmd/trader/main/main.go

install: build
	go install ./cmd/broker/main
	# go install ./cmd/trader/main (uncomment when implemented)

test:
	go test ./...

generate:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -generate types,client -package brokerclient api/openapi.yaml > pkg/brokerclient/client.gen.go

clean:
	rm -f $(BROKER_APP_NAME) $(TRADER_APP_NAME)
