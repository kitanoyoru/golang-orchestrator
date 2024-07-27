.PHONY: test

gomod:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run -c .golangci.yml

test:
	gotestsum --format standard-verbose -- --covermode atomic --coverpkg ./... --count 1 --race ./...

mockgen:
		mockgen -source=task/service/docker/docker.go -destination=task/service/docker/mock/docker.go -package=dockermock
