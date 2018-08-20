package main

import (
	"os"

	"github.com/bzz/lookout-gometalint-analyzer"

	log "gopkg.in/src-d/go-log.v1"
)

func main() {
	withArgs := append([]string(nil), os.Args[1:]...)
	comments := gometalint.RunGometalinter(withArgs)
	log.Infof("%d issues found\n", len(comments))
}
