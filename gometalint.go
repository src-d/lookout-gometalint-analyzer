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
	defaultArgs = []string{
		"--disable-all", "--enable=dupl", "--enable=gas",
		"--enable=gofmt", "--enable=goimports", "--enable=lll", "--enable=misspell",
	}
)

// Comment as returned by gometalint
type Comment struct {
	level string
	file  string
	lino  int32
	text  string
}

// RunGometalinter execs gometalint binary \w pre-configured set of linters
func RunGometalinter(args []string) []Comment {
	dArgs := append([]string(nil), defaultArgs...)
	args = append(dArgs, args...)
	log.Infof("Running '%s %v'\n", bin, args)
	out, _ := exec.Command(bin, args...).Output() // nolint: gas
	// ignoring err, as it's always not nil if anything found

	var comments []Comment
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() { //scan stdout for results
		sp := strings.SplitN(s.Text(), ":", 5)
		if len(sp) != 5 {
			log.Warningf("failed to parse string %s\n", s.Text())
			continue
		}

		file, line, _, severity, msg := sp[0], sp[1], sp[2], sp[3], sp[4]
		c := Comment{
			level: severity,
			file:  file,
			text:  msg,
		}
		lino, err := strconv.Atoi(line)
		if err != nil {
			log.Warningf("failed to parse line number from '%s' in '%s'\n", line, sp)
			continue
		}

		c.lino = int32(lino)
		comments = append(comments, c)
	}
	log.Debugf("Done. %d issues found\n", len(comments))
	return comments
}
