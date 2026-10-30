# gonet

![ci](https://github.com/AAbbey01/gonet/actions/workflows/ci.yml/badge.svg) [![Go](https://img.shields.io/github/go-mod/go-version/AAbbey01/gonet?logo=go&logoColor=white)](https://go.dev) [![Docker](https://img.shields.io/badge/Docker-distroless-2496ED?logo=docker&logoColor=white)](https://github.com/AAbbey01/gonet/blob/main/Dockerfile)

This project was a introductory project into network health and the Go programming language, built on docker and using github actions

## Configuration

```targets.yaml``` expects three kinds of connections: ```tcp```, ```http```, and ```dns``` and comes preloaded with 1 test for each as well as a closed tcp port

## Running it

to run as a go project itself, you can use the three following commands:

```go
go mod tidy
go vet ./...
go run ./cmd/gonet
```

to run via docker image, use the following:

```docker
docker build -t gonet .
docker run --rm -p 8080:8080 gonet
```

## Metrics

### Testing metrics of connections

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

`gonet_up` is 1 when the last check succeeded and 0 when it failed. `gonet_latency_seconds` is a histogram, so Prometheus can also compute percentiles from its `_bucket` series (omitted above for length). `closed-port` is a deliberate failure: nothing listens on port 9, so the connection is refused almost instantly.

### Docker image size

The final image is about 20.4 MB because it contains only the static binary on distroless.

## Tests

  cmd/gonet (main) - 0.0% (has no tests)
  internal/check - 77.8%
  internal/metrics - 100.0%
  