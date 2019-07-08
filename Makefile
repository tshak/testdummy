IMG ?= tshak/testdummy

build: fmt lint
	go build .

fmt:
	go fmt ./...

lint:
	golangci-lint run

docker-build: check-version
	docker build --build-arg VERSION=${VERSION} . -t  testdummy:${VERSION}
	docker tag testdummy:${VERSION} ${IMG}:${VERSION}
	docker tag testdummy:${VERSION} ${IMG}:latest

docker-push: check-version
	docker push ${IMG}:${VERSION}
	docker push ${IMG}:latest

docker-run: docker-build
	docker run -it --rm -p 8000:8000 testdummy

check-version:
	@if test -z "${VERSION}"; then echo "Usage: make docker-builder VERSION=<version>"; exit 1; fi


