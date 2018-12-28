package gometalint

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	types "github.com/gogo/protobuf/types"
	log "gopkg.in/src-d/go-log.v1"
	"gopkg.in/src-d/lookout-sdk.v0/pb"
)

const artificialSep = "___.___"

// Analyzer for the lookout
type Analyzer struct {
	Version    string
	DataClient pb.DataClient
	Args       []string
}

var _ pb.AnalyzerServer = &Analyzer{}

// function to convert pb.types.Value to string argument
type argumentConstructor func(v *types.Value) string

// map of linters with options and argument constructors
var lintersOptions = map[string]map[string]argumentConstructor{
	"lll": map[string]argumentConstructor{
		"maxLen": func(v *types.Value) string {
			var number int

			switch v.GetKind().(type) {
			case *types.Value_StringValue:
				n, err := strconv.Atoi(v.GetStringValue())
				if err != nil {
					log.Warningf("wrong type for lll:maxLen argument")
					return ""
				}
				number = n
			case *types.Value_NumberValue:
				intpart, frac := math.Modf(v.GetNumberValue())
				if frac != 0 {
					log.Warningf("wrong type for lll:maxLen argument")
					return ""
				}
				number = int(intpart)
			default:
				log.Warningf("wrong type for lll:maxLen argument")
				return ""
			}

			if number < 1 {
				return ""
			}

			return fmt.Sprintf("--line-length=%d", number)
		},
	},
}

func (a *Analyzer) NotifyReviewEvent(ctx context.Context, e *pb.ReviewEvent) (
	*pb.EventResponse, error) {
	changes, err := a.DataClient.GetChanges(ctx, &pb.ChangesRequest{
		Head:             &e.Head,
		Base:             &e.Base,
		WantContents:     true,
		WantUAST:         false,
		ExcludeVendored:  true,
		IncludeLanguages: []string{"go"},
	})
	if err != nil {
		log.Errorf(err, "failed to GetChanges from a DataService")
		return nil, err
	}

	tmp, err := ioutil.TempDir("", "gometalint")
	if err != nil {
		log.Errorf(err, "cannot create tmp dir in %s", os.TempDir())
		return nil, err
	}
	defer os.RemoveAll(tmp)
	log.Debugf("Saving files to '%s'", tmp)

	found, saved := 0, 0
	for {
		change, err := changes.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(err, "failed to get a file from DataServer")
			continue
		}

		if change.Head == nil {
			continue
		}

		file := change.Head
		if err = saveTo(file, tmp); err != nil {
			log.Errorf(err, "failed to write file %q", file.Path)
		} else {
			saved++
		}
		found++
	}

	if saved < found {
		log.Warningf("%d/%d Golang files saved. analyzer won't run on non-saved ones", saved, found)
	}
	if saved == 0 {
		log.Debugf("no Golang files to work on. skip running gometalinter")
		return &pb.EventResponse{AnalyzerVersion: a.Version}, nil
	}
	log.Debugf("%d Golang files to work on. running gometalinter", saved)

	withArgs := append(append(a.Args, tmp), a.linterArguments(e.Configuration)...)
	comments := RunGometalinter(withArgs)
	var allComments []*pb.Comment
	for _, comment := range comments {
		origPathFile := revertOriginalPath(comment.file, tmp)
		origPathText := revertOriginalPathIn(comment.text, tmp)
		newComment := pb.Comment{
			File: origPathFile,
			Line: comment.lino,
			Text: origPathText,
		}
		allComments = append(allComments, &newComment)
		log.Debugf("Get comment %v", newComment)
	}

	log.Infof("%d comments created", len(allComments))
	return &pb.EventResponse{
		AnalyzerVersion: a.Version,
		Comments:        allComments,
	}, nil
}

// flattenPath flattens relative path and puts it inside tmp.
func flattenPath(file string, tmp string) string {
	nFile := strings.Join(strings.Split(file, string(os.PathSeparator)), artificialSep)
	nPath := path.Join(tmp, nFile)
	return nPath
}

// revertOriginalPath reverses origina path from a flat one.
func revertOriginalPath(file string, tmp string) string {
	//TrimLeft(, tmp) but works for rel paths
	noTmpfile := file[strings.Index(file, tmp)+len(tmp):]
	origPathFile := strings.TrimLeft(
		path.Join(strings.Split(noTmpfile, artificialSep)...),
		string(os.PathSeparator))
	return origPathFile
}

// revertOriginalPathIn a given text, recovers original path in words
// that have 'artificialSep'.
func revertOriginalPathIn(text string, tmp string) string {
	if strings.LastIndex(text, artificialSep) < 0 {
		return text
	}
	var words []string
	for _, word := range strings.Fields(text) {
		if strings.Index(word, artificialSep) >= 0 {
			word = revertOriginalPath(word, tmp)
		}
		words = append(words, word)
	}
	return strings.Join(words, " ")
}

// saveTo saves a file to given dir, preserving it's original path.
// In case of error it is returned. All files saved this way will
// be in the root of the same dir.
func saveTo(file *pb.File, tmp string) error {
	flatPath := flattenPath(file.Path, tmp)
	return ioutil.WriteFile(flatPath, file.Content, 0644)
}

func (a *Analyzer) NotifyPushEvent(ctx context.Context, e *pb.PushEvent) (*pb.EventResponse, error) {
	return &pb.EventResponse{}, nil
}

func (a *Analyzer) linterArguments(s types.Struct) []string {
	config := s.GetFields()
	if config == nil {
		return nil
	}

	clStruct, ok := config["linters"]
	if !ok || clStruct == nil {
		return nil
	}

	lintersListValue := clStruct.GetListValue()
	if lintersListValue == nil {
		return nil
	}

	var args []string

	for _, v := range lintersListValue.GetValues() {
		if v == nil {
			continue
		}

		sv := v.GetStructValue()
		if sv == nil {
			continue
		}

		fields := sv.GetFields()
		nameV, ok := fields["name"]
		if !ok || nameV == nil {
			continue
		}

		name := nameV.GetStringValue()
		correctLinter := false
		for linter := range lintersOptions {
			if name == linter {
				correctLinter = true
			}
		}

		if !correctLinter {
			log.Warningf("unknown linter %s", name)
			continue
		}

		linterOpts := lintersOptions[name]
		for optionName := range linterOpts {
			optV, ok := fields[optionName]
			if !ok || optV == nil {
				continue
			}

			arg := linterOpts[optionName](optV)
			if arg != "" {
				args = append(args, arg)
			}
		}
	}

	return args
}
