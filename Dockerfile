FROM golang:1.22-alpine as builder

WORKDIR /app
RUN apk add --update --no-cache ca-certificates git

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

ARG VERSION=0.0.0-docker
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X 'main.versionString=${VERSION}'" -o /go/bin/testDummy

FROM gcr.io/distroless/static
COPY --from=builder /go/bin/testDummy /go/bin/testDummy
EXPOSE 8000
CMD ["/go/bin/testDummy"]