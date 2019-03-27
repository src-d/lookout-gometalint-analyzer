package gometalint

import (
	"testing"

	types "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	log "gopkg.in/src-d/go-log.v1"
	"gopkg.in/src-d/lookout-sdk.v0/pb"
)

var logger = log.New(nil)

func TestArgsEmpty(t *testing.T) {
	require := require.New(t)

	inputs := []types.Struct{
		types.Struct{},
		*pb.ToStruct(map[string]interface{}{
			"linters": []map[string]interface{}{},
		}),
		*pb.ToStruct(map[string]interface{}{
			"linters": []map[string]interface{}{
				{
					"name":   "unknown",
					"maxLen": 120,
				},
			},
		}),
		*pb.ToStruct(map[string]interface{}{
			"linters": []map[string]interface{}{
				{
					"name": "lll",
				},
			},
		}),
		*pb.ToStruct(map[string]interface{}{
			"linters": []map[string]interface{}{
				{
					"name":   "lll",
					"maxLen": "not a number",
				},
			},
		}),
		*pb.ToStruct(map[string]interface{}{
			"linters": []map[string]interface{}{
				{
					"name":   "lll",
					"maxLen": 120.1,
				},
			},
		}),
	}

	a := Analyzer{}
	for i, input := range inputs {
		require.Len(a.linterArguments(logger, input), 0, "test case %d; input: %+v", i, input)
	}
}

func TestArgsCorrect(t *testing.T) {
	a := Analyzer{}
	require.Equal(t, []string{"--line-length=120"}, a.linterArguments(logger, *pb.ToStruct(map[string]interface{}{
		"linters": []map[string]interface{}{
			{
				"name":   "lll",
				"maxLen": "120",
			},
		},
	})))

	require.Equal(t, []string{"--line-length=120"}, a.linterArguments(logger, *pb.ToStruct(map[string]interface{}{
		"linters": []map[string]interface{}{
			{
				"name":   "lll",
				"maxLen": 120,
			},
		},
	})))
}

var pathTests = []struct {
	in  string
	out string
}{
	{"a/b.go", "/tmp/a___.___b.go"},
	{"tmp/a/b.go", "/tmp/tmp___.___a___.___b.go"},
	{"a/b/c/d/e.go", "/tmp/a___.___b___.___c___.___d___.___e.go"},
}

func TestPathTransformations(t *testing.T) {
	for _, tt := range pathTests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.out, flattenPath(tt.in, "/tmp"))
			assert.Equal(t, tt.in, revertOriginalPath(tt.out, "/tmp"))
		})
	}
}

func TestPathInTextTransformations(t *testing.T) {
	tmp := "/var/folders/rx/z9zyr71d70x92zwbn3rrjx4c0000gn/T/gometalint584398570"
	text := "duplicate of /var/folders/rx/z9zyr71d70x92zwbn3rrjx4c0000gn/T/gometalint584398570/provider___.___github___.___poster_test.go:549-554 (dupl)"
	expectedText := "duplicate of provider/github/poster_test.go:549-554 (dupl)"

	newText := revertOriginalPathIn(text, tmp)
	if newText != expectedText {
		t.Fatalf("got %q, want %q", newText, expectedText)
	}
}
