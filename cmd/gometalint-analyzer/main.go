package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

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

	bin         = "gometalinter.v2"
	defaultArgs = []string{
		"--disable-all", "--enable=dupl", "--enable=gas",
		"--enable=gofmt", "--enable=goimports", "--enable=lll", "--enable=misspell",
	}
)

type config struct {
	Host       string `envconfig:"HOST" default:"0.0.0.0"`
	Port       int    `envconfig:"PORT" default:"2001"`
	DataServer string `envconfig:"DATA_SERVER_URL" default:"ipv4://localhost:10301"`
}

type comment struct {
	level string
	file  string
	lino  int
	text  string
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

	//TODO(bzz):
	//get changes
	//  for each change
	//    saveFileToTmp(change.File, tmp)
	withArgs := append([]string(nil), os.Args[1:]...)
	withArgs = append(withArgs, tmp)
	runGometalinter(withArgs)
}

func runGometalinter(args []string) {
	args = append(defaultArgs, args...)
	log.Infof("Running '%s %v'\n", bin, args)
	out, _ := exec.Command(bin, args...).Output()
	// ignoring err, as it's always not nil if anything found

	var comments []comment
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() { //scan stdout for results
		sp := strings.SplitN(s.Text(), ":", 5)
		if len(sp) != 5 {
			log.Warningf("failed to parse string %s\n", s.Text())
			continue
		}

		file, line, _, severity, msg := sp[0], sp[1], sp[2], sp[3], sp[4]
		c := comment{
			level: severity,
			file:  file,
			text:  msg,
		}
		lino, err := strconv.Atoi(line)
		if err != nil {
			log.Warningf("failed to parse line number from '%s' in '%s'\n", line, sp)
			continue
		}

		c.lino = lino
		comments = append(comments, c)
	}
	log.Infof("Done. %d issues found\n", len(comments))
}
