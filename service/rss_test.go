package service

import (
	"testing"

	"github.com/benpate/data/expression"
	"github.com/stretchr/testify/assert"
)

func TestRss(t *testing.T) {

	streamService := getTestStreamService()
	factory := streamService.factory

	rss := factory.RSS()

	feed, err := rss.Feed(expression.All())

	assert.Nil(t, err)
	assert.Equal(t, 3, len(feed.Items))
	assert.Equal(t, "My First Stream", feed.Items[0].Title)

	/*
		spew.Dump(err)
		spew.Dump(feed)

		t.Fail()
	*/
}
