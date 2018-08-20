package gometalint

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	log "gopkg.in/src-d/go-log.v1"
)

var (
	bin         = "gometalinter.v2"
	DefaultArgs = []string{
		"--disable-all", "--enable=dupl", "--enable=gas",
		"--enable=gofmt", "--enable=goimports", "--enable=lll", "--enable=misspell",
	}
)

type comment struct {
	level string
	file  string
	lino  int
	text  string
}

func RunGometalinter(args []string) []comment {
	dArgs := append([]string(nil), DefaultArgs...)
	args = append(dArgs, args...)
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
	log.Debugf("Done. %d issues found\n", len(comments))
	return comments
}
