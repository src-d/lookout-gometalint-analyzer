package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bzz/lookout-gometalint-analyzer"

	"github.com/kelseyhightower/envconfig"
	"github.com/sanity-io/litter"
	log "gopkg.in/src-d/go-log.v1"
)

var usageMessage = fmt.Sprintf(`usage: %s [-version] [OPTIONS]

%s is a lookout analyzer implementation, based on https://github.com/alecthomas/gometalinter.

OPTIONS - any of the supported by gometalinter.
`, name, name)

var (
	name        = "gometalint-analyzer"
	version     string
	build       string
	versionFlag = flag.Bool("version", false, "show version")
)

type config struct {
	Host       string `envconfig:"HOST" default:"0.0.0.0"`
	Port       int    `envconfig:"PORT" default:"2001"`
	DataServer string `envconfig:"DATA_SERVER_URL" default:"ipv4://localhost:10301"`
}

func main() {
	litter.Config.Compact = true
	flag.Usage = func() {
		fmt.Printf(usageMessage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s %s built on %s\n", name, version, build)
		return
	}

	var conf config
	envconfig.MustProcess("GOMETALINT", &conf)

	log.Infof("Starting %s, %s", name, litter.Sdump(conf))
	tmp, err := ioutil.TempDir("", "gometalint")
	if err != nil {
		log.Errorf(err, "cannot create tmp dir in %s", os.TempDir())
		return
	}
	defer os.RemoveAll(tmp)

	//TODO(bzz): move to Analyzer
	//get changes
	//  for each change
	//    saveFileToTmp(change.File, tmp)
	withArgs := append([]string(nil), os.Args[1:]...)
	withArgs = append(withArgs, tmp)
	_ = gometalint.RunGometalinter(withArgs)
}
