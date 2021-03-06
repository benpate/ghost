package domain

import (
	"testing"

	"github.com/benpate/ghost/config"
	"github.com/benpate/ghost/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestService(t *testing.T) {

	factory, err := NewFactory(config.Domain{
		ConnectString: "mongodb://127.0.0.1/ghost",
		DatabaseName:  "ghost",
	})

	require.Nil(t, err)

	streamService := factory.Stream()

	streamID, err := primitive.ObjectIDFromHex("5f84e964e49c4c226eb51097")
	require.Nil(t, err)

	it, err := streamService.ListByParent(streamID)

	require.Nil(t, err)

	stream := model.Stream{}

	for it.Next(&stream) {
		spew.Dump(stream)
	}
}
