# [![Build Status](https://travis-ci.org/src-d/lookout-gometalint-analyzer.svg)](https://travis-ci.org/src-d/lookout-gometalint-analyzer) lookout analyzer: gometalint

A [lookout](https://github.com/src-d/lookout/) analyzer implementation that uses [gometalinter](https://github.com/alecthomas/gometalinter).

It only applies 6 checks from gometalinter that are file-level, and skips dir and package level ones.

_Disclamer: this is not official product, but only serves the purpose of testing the lookout._


# Build

```
make packages
```

will produce binaries for multiple architectures under `./build`.

# Dependencies

Requires stable version of `gometalinter.v2` binary avabilable in PATH.

To install, do
```
go get -u gopkg.in/alecthomas/gometalinter.v2
gometalinter.v2 --install
```
This will also install a number of linter binaries, vendored by gometalinter.

# Example of utilization

With `lookout-sdk` binary from the latest release of [SDK](https://github.com/src-d/lookout/releases)

```
$ lookout-gometalint

$ lookout-sdk review -v ipv4://localhost:2001 \
    --from c99dcdff172f1cb5505603a45d054998cb4dd606 \
    --to 3a9d78bdd1139c929903885ecb8f811931b8aa70
```

# Configuration

| Variable | Default | Description |
| -- | -- | -- |
| `GOMETALINT_HOST` | `0.0.0.0` | IP address to bind the gRCP serve |
| `GOMETALINT_PORT` | `2001` | Port to bind the gRPC server |
| `GOMETALINT_SERVER_URL` | `ipv4://localhost:10302` | gRPC URL of the [Data service](https://github.com/src-d/lookout/tree/master/docs#components)
| `GOMETALINT_LOG_LEVEL` | `info` | Logging level (info, debug, warning or error) |


# Licens

AGPLv3
