# Introduction

TestDummy is a rudimentary utility service to assist in testing platform schedulers such as Kubernetes and Nomad. It can be used as a simple "hello world" service or it can be used to test a couple of common failure modes (a humble "chaos monkey" in a way).

# Usage
```
docker run -it --rm -p 8000:8000 ghcr.io/tshak/testdummy:latest
```

Docker images are hosted on [GitHub](https://github.com/tshak/testdummy/pkgs/container/testdummy).

## Environment variables

The following environment variables can be used to configure the service.

| Name | Default | Description
| -| -| -
| `TESTDUMMY_BIND_ADDRESS` | `:8000` | The address to bind the service to
| `TESTDUMMY_HEALTHY` | `true` | Boolean. When `true` `/health` returns an HTTP 200 response. When `false` `/health`  returns an empty HTTP 500 response.
| `TESTDUMMY_HEALTHY_AFTER_SECONDS` | | When set, sets the health status to `false` (regardless of `TESTDUMMY_HEALTHY`) until the specified number of seconds. This is useful for testing post deployment "warmup" scenarios.
| `TESTDUMMY_PANIC_SECONDS` | | When set, specifies the number of seconds to wait before panicking. This is useful for testing crash recovery scenarios.
| `TESTDUMMY_ENABLE_REQUEST_LOGGING` | `false` | Logs all requests to stdout
| `TESTDUMMY_ENABLE_ENV` | `false` | Enable the `/env` endpoint which dumps env vars. :warning: This can be a security risk so enable with caution
| `TESTDUMMY_ROOT_PATH` | `/` | The root path for all routes
| `TESTDUMMY_STRESS_CPU_DURATION` | `0s` | Ping endpoints will perform a naive CPU stress test on all cores for the supplied duration


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

