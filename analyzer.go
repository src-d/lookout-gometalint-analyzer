package gometalint

import (
	"context"

	"github.com/src-d/lookout"
)

type Analyzer struct {
	Version    string
	DataClient *lookout.DataClient
}

var _ lookout.AnalyzerServer = &Analyzer{}

func (a *Analyzer) NotifyReviewEvent(ctx context.Context, e *lookout.ReviewEvent) (
	*lookout.EventResponse, error) {
	return &lookout.EventResponse{}, nil
}

func (a *Analyzer) NotifyPushEvent(ctx context.Context, e *lookout.PushEvent) (*lookout.EventResponse, error) {
	return &lookout.EventResponse{}, nil
}
