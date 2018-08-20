package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/sanity-io/litter"
	log "gopkg.in/src-d/go-log.v1"
)

var usageMessage = fmt.Sprintf(`usage: %s [-version]

%s is a lookout analyzer implementation, based on https://github.com/alecthomas/gometalinter.
`, name, name)

func usage() {
	fmt.Printf(usageMessage)
	os.Exit(2)
}

var (
	name        = "gometalint-analyzer"
	version     string
	build       string
	versionFlag = flag.Bool("version", false, "show version")

	bin  = "gometalinter.v2"
	args = []string{
		"--disable-all", "--enable=dupl", "--enable=gas",
		"--enable=gofmt", "--enable=goimports", "--enable=lll", "--enable=misspell",
	}
)

type config struct {
	Host       string `envconfig:"HOST" default:"0.0.0.0"`
	Port       int    `envconfig:"PORT" default:"2001"`
	DataServer string `envconfig:"DATA_SERVER_URL" default:"ipv4://localhost:10301"`
}

func main() {
	litter.Config.Compact = true
	flag.Usage = usage
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s %s built on %s\n", name, version, build)
		return
	}

	var conf config
	envconfig.MustProcess("GOMETALINT", &conf)

	log.Infof("Starting %s, %s", name, litter.Sdump(conf))
}
