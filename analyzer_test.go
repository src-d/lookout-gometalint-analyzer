package gometalint

import (
	"testing"

	types "github.com/gogo/protobuf/types"
	"github.com/src-d/lookout/util/grpchelper"
	"github.com/stretchr/testify/require"
)

func TestArgsEmpty(t *testing.T) {
	require := require.New(t)

	inputs := []types.Struct{
		types.Struct{},
		*grpchelper.ToPBStruct(map[string]interface{}{
			"linters": []map[string]interface{}{},
		}),
		*grpchelper.ToPBStruct(map[string]interface{}{
			"linters": []map[string]interface{}{
				{
					"name":   "unknown",
					"maxLen": 120,
				},
			},
		}),
		*grpchelper.ToPBStruct(map[string]interface{}{
			"linters": []map[string]interface{}{
				{
					"name": "lll",
				},
			},
		}),
		*grpchelper.ToPBStruct(map[string]interface{}{
			"linters": []map[string]interface{}{
				{
					"name":   "lll",
					"maxLen": "not a number",
				},
			},
		}),
		*grpchelper.ToPBStruct(map[string]interface{}{
			"linters": []map[string]interface{}{
				{
					"name":   "lll",
					"maxLen": 120.1,
				},
			},
		}),
	}

	a := Analyzer{}
	for _, input := range inputs {
		require.Len(a.linterArguments(input), 0)
	}
}

func TestArgsCorrect(t *testing.T) {
	a := Analyzer{}
	require.Equal(t, []string{"--line-length=120"}, a.linterArguments(*grpchelper.ToPBStruct(map[string]interface{}{
		"linters": []map[string]interface{}{
			{
				"name":   "lll",
				"maxLen": "120",
			},
		},
	})))

	require.Equal(t, []string{"--line-length=120"}, a.linterArguments(*grpchelper.ToPBStruct(map[string]interface{}{
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
			flat := faltternPath(tt.in, "/tmp")
			if flat != tt.out {
				t.Errorf("forward: got %q, want %q", flat, tt.out)
			}

			orig := revertOriginalPath(tt.out, "/tmp")
			if orig != tt.in {
				t.Errorf("backward: got %q, want %q", orig, tt.in)
			}
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
