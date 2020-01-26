# Introduction

TestDummy is a rudimentary utility service to assist in testing platform schedulers such as Kubernetes and Nomad. It can be used as a simple "hello world" service or it can be used to test a couple of common failure modes (a humble "chaos monkey" in a way).

# Usage
```
docker run -it --rm -p 8000:8000 tshak/testdummy:latest
```

Docker images are hosted on [Docker Hub](https://hub.docker.com/r/tshak/testdummy). See the below sections for additional options.

## Environment variables

The following environment variables can be used to configure the service.

| Name | Default | Description
| -| -| -
| `TESTDUMMY_BIND_ADDRESS` | `localhost:8000` | The address to bind the service to
| `TESTDUMMY_HEALTHY` | `true` | Boolean. When `true` `/health` returns an HTTP 200 response. When `false` `/health`  returns an empty HTTP 500 response.
| `TESTDUMMY_HEALTHY_AFTER_SECONDS` | | When set, sets the health status to `false` (regardless of `TESTDUMMY_HEALTHY`) until the specified number of seconds. This is useful for testing post deployment "warmup" scenarios.
| `TESTDUMMY_PANIC_SECONDS` | | When set, specifies the number of seconds to wait before panicking. This is useful for testing crash recovery scenarios.


## API

| Path | Description
|  -| -
| `/` _or_ `/ping` | Returns `pong`
| `/echo` | Returns request body
| `/env` | Returns all environment variables in the format `NAME=VALUE`, one per line.
| `/exit?code={exitCode}` | Causes the process to exit with a default exit code of `1`. Use the `code` parameter to customize the exit code.
| <code>/health?healthy={true&#124;false}</code> | Returns an empty HTTP 200 if healthy, otherwise an empty HTTP 500.  The healthy status can be changed using the optional `healthy` query parameter. This status persists for the lifetime of the process. See the [Environment Variables](#environment-variables) section for related options.
| `/status?status={statusCode}` | Returns the supplied HTTP status code
| `/version` | Returns the testdummy version

