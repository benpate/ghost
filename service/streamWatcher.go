package service

import (
	"context"
	"fmt"

	"github.com/benpate/derp"
	"github.com/benpate/ghost/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewStreamWatcher initiates a mongodb change stream to on every updates to Stream data objects
func NewStreamWatcher(collection *mongo.Collection) chan model.Stream {

	fmt.Println("templatesource.NewStreamWatcher: attempting to create a new watcher.")
	result := make(chan model.Stream)

	ctx := context.Background()

	cs, err := collection.Watch(ctx, mongo.Pipeline{})

	if err != nil {
		derp.Report(derp.Wrap(err, "ghost.service.Watcher", "Unable to open Mongodb Change Stream"))
		return result
	}

	go func() {

		for cs.Next(ctx) {

			var event struct {
				Stream model.Stream `bson:"fullDocument"`
			}

			if err := cs.Decode(&event); err != nil {
				derp.Report(err)
				continue
			}

			// Skip "zero" sreams
			if event.Stream.StreamID.IsZero() {
				continue
			}

			result <- event.Stream
		}
	}()

	return result
}
