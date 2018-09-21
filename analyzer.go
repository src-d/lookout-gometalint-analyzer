package gometalint

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	types "github.com/gogo/protobuf/types"
	"github.com/src-d/lookout"
	log "gopkg.in/src-d/go-log.v1"
)

const artificialSep = "___.___"

// Analyzer for the lookout
type Analyzer struct {
	Version    string
	DataClient *lookout.DataClient
	Args       []string
}

var _ lookout.AnalyzerServer = &Analyzer{}

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
				number = int(v.GetNumberValue())
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

func (a *Analyzer) NotifyReviewEvent(ctx context.Context, e *lookout.ReviewEvent) (
	*lookout.EventResponse, error) {
	changes, err := a.DataClient.GetChanges(ctx, &lookout.ChangesRequest{
		Head:            &e.Head,
		Base:            &e.Base,
		WantContents:    true,
		WantLanguage:    true,
		WantUAST:        false,
		ExcludeVendored: true,
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

	saved := 0
	for changes.Next() {
		change := changes.Change()
		if change.Head == nil {
			continue
		}

		// analyze only changes in Golang
		if strings.HasPrefix(strings.ToLower(change.Head.Language), "go") {
			tryToSaveTo(change.Head, tmp)
			saved++
		}
	}
	if changes.Err() != nil {
		log.Errorf(changes.Err(), "failed to get a file from DataServer")
	}

	if saved == 0 {
		log.Debugf("no Golang files found. skip running gometalinter")
		return &lookout.EventResponse{AnalyzerVersion: a.Version}, nil
	}
	log.Debugf("%d Golang files found. running gometalinter", saved)

	withArgs := append(append(a.Args, tmp), a.linterArguments(e.Configuration)...)
	comments := RunGometalinter(withArgs)
	var allComments []*lookout.Comment
	for _, comment := range comments {
		origPathFile := revertOriginalPath(comment.file, tmp)
		origPathText := revertOriginalPathIn(comment.text, tmp)
		newComment := lookout.Comment{
			File: origPathFile,
			Line: comment.lino,
			Text: origPathText,
		}
		allComments = append(allComments, &newComment)
		log.Debugf("Get comment %v", newComment)
	}

	log.Infof("%d comments created", len(allComments))
	return &lookout.EventResponse{
		AnalyzerVersion: a.Version,
		Comments:        allComments,
	}, nil
}

// faltternPath flattens relative path and puts it inside tmp.
func faltternPath(file string, tmp string) string {
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

// tryToSaveTo saves a file to given dir, preserving it's original path.
// It only loggs any errors and does not fail. All files saved this way will
// be in the root of the same dir.
func tryToSaveTo(file *lookout.File, tmp string) {
	flatPath := faltternPath(file.Path, tmp)
	err := ioutil.WriteFile(flatPath, file.Content, 0644)
	if err != nil {
		log.Errorf(err, "failed to write a file %q", flatPath)
	}
}
func (a *Analyzer) NotifyPushEvent(ctx context.Context, e *lookout.PushEvent) (*lookout.EventResponse, error) {
	return &lookout.EventResponse{}, nil
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
