package main

import (
	"os"
	"strings"
	"testing"

	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	require := require.New(t)
	os.Args = []string{"-h"}

	var stdout, stderr string

	stdout = capturer.CaptureStdout(func() {
		stderr = capturer.CaptureStderr(func() {
			main()
		})
	})

	require.Empty(stderr)
	require.True(strings.HasPrefix("usage: ", stdout))
}
