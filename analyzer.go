package gometalint

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/src-d/lookout"
	log "gopkg.in/src-d/go-log.v1"
)

type Analyzer struct {
	Version    string
	DataClient *lookout.DataClient
	Args       []string
}

var _ lookout.AnalyzerServer = &Analyzer{}

func (a *Analyzer) NotifyReviewEvent(ctx context.Context, e *lookout.ReviewEvent) (
	*lookout.EventResponse, error) {
	tmp, err := ioutil.TempDir("", "gometalint")
	if err != nil {
		log.Errorf(err, "cannot create tmp dir in %s", os.TempDir())
		return nil, err
	}
	defer os.RemoveAll(tmp)

	changes, err := a.DataClient.GetChanges(ctx, &lookout.ChangesRequest{
		Head:         &e.Head,
		Base:         &e.Base,
		WantContents: true,
		WantLanguage: true,
	})
	if err != nil {
		log.Errorf(err, "failed on GetChanges from the DataService")
	}

	for changes.Next() {
		//    saveFileToTmp(change.Head.File, tmp)
		change := changes.Change()
		file := path.Join(tmp, path.Base(change.Head.Path))
		err = ioutil.WriteFile(file, change.Head.Content, 0644)
		if err != nil {
			log.Errorf(err, "failed to write a file %s", file)
		}
		log.Infof("Saved file:'%s'", file)
	}
	if changes.Err() != nil {
		log.Errorf(changes.Err(), "failed to get a file from DataServer")
	}

	withArgs := append(a.Args, tmp)
	comments := RunGometalinter(withArgs)
	var allComments []*lookout.Comment
	for _, comment := range comments {
		newComment := lookout.Comment{
			File: comment.file,
			Line: comment.lino,
			Text: comment.text,
		}
		allComments = append(allComments, &newComment)
	}

	return &lookout.EventResponse{
		AnalyzerVersion: a.Version,
		Comments:        allComments,
	}, nil
}

func (a *Analyzer) NotifyPushEvent(ctx context.Context, e *lookout.PushEvent) (*lookout.EventResponse, error) {
	return &lookout.EventResponse{}, nil
}
