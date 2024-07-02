IMG ?= tshak/testdummy

build: fmt lint
	go build .

fmt:
	go fmt ./...

lint:
	golangci-lint run

docker-run: docker-build
	docker run -it --rm -p 8000:8000 testdummy

test:
	go test -v ./...

# compile and run unit tests on change.
# requires gotestsum
watch:
	gotestsum --watch



