package gometalint

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"strings"

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

	withArgs := append(a.Args, tmp)
	comments := RunGometalinter(withArgs)
	var allComments []*lookout.Comment
	for _, comment := range comments {
		//TrimLeft(, tmp) but \w rel path
		file := comment.file[strings.LastIndex(comment.file, tmp)+len(tmp):]
		newComment := lookout.Comment{
			File: strings.TrimLeft(
				path.Join(strings.Split(file, artificialSep)...),
				string(os.PathSeparator)),
			Line: comment.lino,
			Text: comment.text,
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

// tryToSaveTo saves a file to given dir, preserving it's original path.
// It only loggs any errors and does not fail. All files saved this way will
// be in the root of the same dir.
func tryToSaveTo(file *lookout.File, tmp string) {
	nFile := strings.Join(strings.Split(file.Path, string(os.PathSeparator)), artificialSep)
	nPath := path.Join(tmp, nFile)
	log.Debugf("Saving file '%s', as '%s'", file.Path, nPath)
	err := ioutil.WriteFile(nPath, file.Content, 0644)
	if err != nil {
		log.Errorf(err, "failed to write a file %s", nPath)
	}
}
func (a *Analyzer) NotifyPushEvent(ctx context.Context, e *lookout.PushEvent) (*lookout.EventResponse, error) {
	return &lookout.EventResponse{}, nil
}
