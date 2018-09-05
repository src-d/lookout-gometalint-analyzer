package gometalint

import (
	"testing"

	types "github.com/gogo/protobuf/types"
	"github.com/src-d/lookout/pb"
	"github.com/stretchr/testify/require"
)

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
	for _, input := range inputs {
		require.Len(a.linterArguments(input), 0)
	}
}

func TestArgsCorrect(t *testing.T) {
	a := Analyzer{}
	require.Equal(t, []string{"--line-length=120"}, a.linterArguments(*pb.ToStruct(map[string]interface{}{
		"linters": []map[string]interface{}{
			{
				"name":   "lll",
				"maxLen": "120",
			},
		},
	})))

	require.Equal(t, []string{"--line-length=120"}, a.linterArguments(*pb.ToStruct(map[string]interface{}{
		"linters": []map[string]interface{}{
			{
				"name":   "lll",
				"maxLen": 120,
			},
		},
	})))
}
