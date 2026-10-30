# build stage: match the Go version in go.mod
FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gonet ./cmd/gonet

# run stage: just the binary, running as a non-root user
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /gonet /gonet
COPY targets.yaml /targets.yaml
EXPOSE 8080
ENTRYPOINT ["/gonet"]
CMD ["-config", "/targets.yaml"]