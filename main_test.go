package main

import (
	"testing"

	"gotest.tools/assert"

	"github.com/kyoukaya/catte/internal/fflogs"
	"github.com/kyoukaya/catte/internal/xivdata"
)

const (
	fflogsClientID = ""
	fflogsToken    = ""
)

func TestF(t *testing.T) {
	// TODO: fix up tests with mocks
	fflogsClient, err := fflogs.NewClient(fflogsClientID, fflogsToken)
	assert.NilError(t, err)
	app := &App{
		xivds:  xivdata.NewDataSource(),
		fflogs: fflogsClient,
	}
	app.MessageHandler(nil, nil)
}
