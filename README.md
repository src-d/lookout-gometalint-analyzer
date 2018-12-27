# [![Build Status](https://travis-ci.org/src-d/lookout-gometalint-analyzer.svg)](https://travis-ci.org/src-d/lookout-gometalint-analyzer) lookout analyzer: gometalint

A [lookout](https://github.com/src-d/lookout/) analyzer implementation that
uses [gometalinter](https://github.com/alecthomas/gometalinter).

**Disclaimer:** This is not an official product, but can be used to verify that
your lookout installation is working.

This analyzer only enables the gometalinter checks that are file-level, and
skips directory- and package-level checks.  The currently-enabled linters are
(from [gometalint.go](gometalint.go)):

* `gofmt`
* http://godoc.org/github.com/client9/misspell/cmd/misspell
* http://godoc.org/github.com/mibk/dupl
* http://godoc.org/github.com/securego/gosec/cmd/gosec
* http://godoc.org/github.com/walle/lll/cmd/lll
* http://godoc.org/golang.org/x/tools/cmd/goimports

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

$ lookout-sdk review --log-level=debug \
    --from c99dcdff172f1cb5505603a45d054998cb4dd606 \
    --to 3a9d78bdd1139c929903885ecb8f811931b8aa70
```

# Configuration

| Variable | Default | Description |
| -- | -- | -- |
| `GOMETALINT_HOST` | `0.0.0.0` | IP address to bind the gRCP serve |
| `GOMETALINT_PORT` | `9930` | Port to bind the gRPC server |
| `GOMETALINT_DATA_SERVICE_URL` | `ipv4://localhost:10301` | gRPC URL of the [Data service](https://github.com/src-d/lookout/tree/master/docs#components)
| `GOMETALINT_LOG_LEVEL` | `info` | Logging level ("info", "debug", "warning" or "error") |


# License

AGPLv3, see [LICENSE](LICENSE)
