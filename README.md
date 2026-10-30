# gonet

![ci](https://github.com/AAbbey01/gonet/actions/workflows/ci.yml/badge.svg) [![Go](https://img.shields.io/github/go-mod/go-version/AAbbey01/gonet?logo=go&logoColor=white)](https://go.dev) [![Docker](https://img.shields.io/badge/Docker-distroless-2496ED?logo=docker&logoColor=white)](https://github.com/AAbbey01/gonet/blob/main/Dockerfile)

gonet checks a list of network targets (TCP connect, HTTP status, DNS lookup) concurrently on a timer and exports the results as Prometheus metrics. It ships as a ~20 MB distroless Docker image, with tests and CI on GitHub Actions.

## Configuration

gonet reads its targets from `targets.yaml`:

```yaml
interval: 30s   # how often each target is checked (default 30s)
timeout: 5s     # how long one check may take (default 5s)
targets:
  - name: cloudflare
    kind: tcp
    address: 1.1.1.1:443
  - name: github
    kind: http
    address: https://github.com
  - name: example
    kind: dns
    address: example.com
  - name: closed-port
    kind: tcp
    address: 127.0.0.1:9
```

| Kind | `address` format | Up when |
| --- | --- | --- |
| `tcp` | `host:port` | the TCP connection opens |
| `http` | full URL, including `http://` or `https://` | the response status is below 400 |
| `dns` | hostname | the name resolves to at least one address |

`closed-port` fails on purpose, so the output always shows a target that is down. gonet refuses to start if the file has no targets or a target has an unknown `kind`.

## Running it

With Go installed:

```sh
go run ./cmd/gonet
```

With Docker:

```sh
docker build -t gonet .
docker run --rm -p 8080:8080 gonet
```

Then view the metrics:

```sh
curl http://localhost:8080/metrics
```

| Flag | Default | Purpose |
| --- | --- | --- |
| `-config` | `targets.yaml` | path to the targets file |
| `-addr` | `:8080` | listen address for `/metrics` |

To use your own targets file with Docker, mount it over the built-in one (Windows cmd shown; use `$(pwd)` on Linux/macOS):

```sh
docker run --rm -p 8080:8080 -v "%cd%\targets.yaml:/targets.yaml" gonet
```

Stop it with Ctrl+C or `docker stop`; gonet shuts down cleanly either way.

## Sample output

From one run on my laptop (3 check rounds, 30 s apart):

```text
# HELP gonet_up 1 if the last check of the target succeeded, 0 if it failed.
# TYPE gonet_up gauge
gonet_up{kind="dns",target="example"} 1
gonet_up{kind="http",target="github"} 1
gonet_up{kind="tcp",target="closed-port"} 0
gonet_up{kind="tcp",target="cloudflare"} 1

gonet_latency_seconds_sum{kind="dns",target="example"} 0.1054735
gonet_latency_seconds_count{kind="dns",target="example"} 3
gonet_latency_seconds_sum{kind="http",target="github"} 0.1966048
gonet_latency_seconds_count{kind="http",target="github"} 3
gonet_latency_seconds_sum{kind="tcp",target="closed-port"} 0.0012824
gonet_latency_seconds_count{kind="tcp",target="closed-port"} 3
gonet_latency_seconds_sum{kind="tcp",target="cloudflare"} 0.06258459999999999
gonet_latency_seconds_count{kind="tcp",target="cloudflare"} 3
```

| Target | Kind | Up | Average check time (`_sum` ÷ `_count`) |
| --- | --- | --- | --- |
| example.com | DNS lookup | yes | 35 ms |
| github.com | HTTP GET | yes | 66 ms |
| 1.1.1.1:443 (Cloudflare) | TCP connect | yes | 21 ms |
| 127.0.0.1:9 | TCP connect | no (refused, by design) | 0.4 ms |

`gonet_up` is 1 when the last check succeeded and 0 when it failed. `gonet_latency_seconds` is a histogram, so Prometheus can also compute percentiles from its `_bucket` series (omitted above for length).

The final Docker image is about 20.4 MB because it holds only the static Go binary on a distroless base, with no shell or package manager.

## Design

- **One goroutine per target.** Each target runs on its own ticker, so a slow or hanging target never delays the others.
- **A deadline on every check.** Each check gets a `context.WithTimeout`, so a hung connection, request or DNS lookup is cancelled after `timeout`.
- **A custom Prometheus registry** instead of the global default, so each test gets a fresh set of metrics.
- **Clean shutdown.** `signal.NotifyContext` stops the checkers and the HTTP server on Ctrl+C or `docker stop`.
- **Fail fast on bad config.** Invalid YAML, an empty target list or an unknown kind stops gonet at startup with a clear error.
  - This was clear from a personal addition to the `targets.yaml`, where I specified a http connection to a local server but it failed instantly. This was because the code checked for that `http`, and threw an error when it couldn't find it from the direct connection method.

## Tests

```sh
go vet ./...
go test -cover ./...
```

| Package | Coverage | What the tests cover |
| --- | --- | --- |
| `internal/check` | 77.8% | TCP to an open vs closed port, HTTP 200 vs 500, a 50 ms deadline against a slow server, config defaults and validation |
| `internal/metrics` | 100% | `gonet_up` set to 1 on success and 0 on failure |
| `cmd/gonet` | 0% | wiring only (flags, goroutines, HTTP server); no tests |

The tests use local throwaway servers (`httptest`, `net.Listen` on `127.0.0.1`), so they pass offline. GitHub Actions runs `go vet`, the tests with the race detector, a build and a Docker build on every push.

## How it was built

Built as an introduction to Go, network health checks, Docker and GitHub Actions. I followed a step-by-step guide made with Claude that supplied the project layout and function signatures, then wrote the check, config, metrics and scheduling code myself. The `internal/check` tests were written with help from Claude, and I debugged them alongside my code. Human checked and edited for safety and consistency
