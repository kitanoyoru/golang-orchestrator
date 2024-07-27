.PHONY: test

gomod:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run -c .golangci.yml

test:
	gotestsum --format standard-verbose -- --covermode atomic --coverpkg ./... --count 1 --race ./...

