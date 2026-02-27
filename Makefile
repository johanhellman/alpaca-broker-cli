.PHONY: build install test generate clean

APP_NAME = alpaca-cli

build:
	go build -o $(APP_NAME) main.go

install:
	go install

test:
	go test ./...

generate:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -generate types,client -package client api/openapi.yaml > pkg/client/client.gen.go

clean:
	rm -f $(APP_NAME)
