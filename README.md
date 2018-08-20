# [![Build Status](https://travis-ci.org/bzz/lookout-gometalint-analyzer.svg)](https://travis-ci.org/bzz/lookout-gometalint-analyzer) lookout analyzer: gometalint

A [lookout](https://github.com/src-d/lookout/) analyzer implementation that uses [gometalinter](https://github.com/alecthomas/gometalinter).

_Disclamer: this is not official product, but only serves the purpose of testing the lookout._


# Build

```
make packages
```

# Dependencies

Requires stable version of `gometalinter.v2` binary avabilable in PATH.

To install, do
```
go get -u gopkg.in/alecthomas/gometalinter.v2
```
This will also install a number of linter binaries, vendored by gometalinter.

# Example of utilization

With `lookout` binary from the latest release of [SDK](https://github.com/src-d/lookout/releases)

```
./lookout-gometalint

./look review -v ipv4://localhost:2001 \
    --from c99dcdff172f1cb5505603a45d054998cb4dd606 \
    --to 3a9d78bdd1139c929903885ecb8f811931b8aa70
```

# Configuration

| Variable | Default | Description |
| -- | -- | -- |
| `GOMETALINT_HOST` | `0.0.0.0` | IP address to bind the gRCP serve |
| `GOMETALINT_PORT` | `2001` | Port to bind the gRPC server |
| `GOMETALINT_SERVER_URL` | `ipv4://localhost:10302` | gRPC URL of the [Data service](https://github.com/src-d/lookout/tree/master/docs#components)


# Licens

AGPLv3
