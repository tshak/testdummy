IMG ?= tshak/testdummy

build: fmt lint
	go build .

fmt:
	go fmt ./...

lint:
	golangci-lint run -e SA5004

docker-run: docker-build
	docker run -it --rm -p 8000:8000 testdummy



