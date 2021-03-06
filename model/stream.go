package model

import (
	"github.com/benpate/convert"
	"github.com/benpate/data/journal"
	"github.com/benpate/datatype"
	"github.com/benpate/derp"
	"github.com/benpate/ghost/content"
	"github.com/benpate/path"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Stream corresponds to a top-level path on any Domain.
type Stream struct {
	StreamID        primitive.ObjectID `json:"streamId"        bson:"_id"`                     // Unique identifier of this Stream.  (NOT USED PUBLICLY)
	ParentID        primitive.ObjectID `json:"parentId"        bson:"parentId"`                // Unique identifier of the "parent" stream. (NOT USED PUBLICLY)
	TemplateID      string             `json:"templateId"      bson:"templateId"`              // Unique identifier (name) of the Template to use when rendering this Stream in HTML.
	StateID         string             `json:"stateId"         bson:"stateId"`                 // Unique identifier of the State this Stream is in.  This is used to populate the State information from the Template service at load time.
	Criteria        Criteria           `json:"criteria"        bson:"criteria"`                // Criteria for which users can access this stream.
	Token           string             `json:"token"           bson:"token"`                   // Unique value that identifies this element in the URL
	URL             string             `json:"url"             bson:"url"`                     // Unique URL of this Stream.  This duplicates the "token" field a bit, but it (hopefully?) makes access easier.
	Label           string             `json:"label"           bson:"label"`                   // Text to display in lists of streams, probably displayed at top of stream page, too.
	Description     string             `json:"description"     bson:"description"`             // Brief summary of this stream, used in lists of streams
	Content         content.Content    `json:"content"         bson:"content"`                 // Content objects for this stream.
	Data            datatype.Map       `json:"data"            bson:"data"`                    // Set of data to populate into the Template.  This is validated by the JSON-Schema of the Template.
	Tags            []string           `json:"tags"            bson:"tags"`                    // Organizational Tags
	AuthorID        primitive.ObjectID `json:"authorId"        bson:"authorId"`                // Unique identifier of the person who created this stream (NOT USED PUBLICLY)
	AuthorName      string             `json:"authorName"      bson:"authorName"`              // Full name of the person who created this stream
	AuthorImage     string             `json:"authorImage"     bson:"authorImage"`             // URL of an image to use for the person who created this stream
	AuthorURL       string             `json:"authorURL"       bson:"authorURL"`               // URL address of the person who created this stream
	ThumbnailImage  string             `json:"thumbnailImage"  bson:"thumbnailImage"`          // Image to display next to the stream in lists.
	BubbleUpdates   bool               `json:"bubbleUpdates"   bson:"bubbleUpdates"`           // If TRUE then updates are sent to the PARENT, instead of THIS stream.  This *should* be controlled by the Template.
	SourceID        primitive.ObjectID `json:"sourceId"        bson:"sourceId,omitempty"`      // Internal identifier of the source configuration that generated this stream
	SourceURL       string             `json:"sourceUrl"       bson:"sourceUrl,omitempty"`     // URL of the original document published by the source server
	SourceUpdated   int64              `json:"sourceUpdated"   bson:"sourceUpdated,omitempty"` // Date the the source updated the original content.
	PublishDate     int64              `json:"publishDate"     bson:"publishDate"`             // Unix timestamp of the date/time when this document is/was/will be first available on the domain.
	UnPublishDate   int64              `json:"unpublishDate"   bson:"unpublishDate"`           // Unix timestemp of the date/time when this document will no longer be available on the domain.
	Rank            int                `json:"Rank"            bson:"Rank"`                    // Rank allows for a manual sort of streams
	journal.Journal `json:"journal" bson:"journal"`
}

// NewStream returns a fully initialized Stream object.
func NewStream() Stream {

	streamID := primitive.NewObjectID()

	return Stream{
		StreamID: streamID,
		Token:    streamID.Hex(),
		StateID:  "new",
		Criteria: NewCriteria(),
		Tags:     make([]string, 0),
		Data:     make(datatype.Map),
		Content: content.Content{
			0: {
				Type: "CONTAINER",
				Refs: []int{1},
				Data: datatype.Map{
					"style": "ROWS",
				},
			},
			1: {
				Type: "WYSIWYG",
			},
		},
	}
}

// ID returns the primary key of this object
func (stream Stream) ID() string {
	return stream.StreamID.Hex()
}

// HasParent returns TRUE if this Stream has a valid parentID
func (stream Stream) HasParent() bool {
	return !stream.ParentID.IsZero()
}

// GetPath implements the path.Getter interface.  It looks up
// data within this Stream and returns it to the caller.
func (stream Stream) GetPath(p path.Path) (interface{}, error) {

	if p.IsEmpty() {
		return nil, derp.New(500, "ghost.model.Stream", "Unrecognized path", p)
	}

	property := p.Head()

	// Properties that can be set
	switch property {

	case "label":
		return stream.Label, nil

	case "description":
		return stream.Description, nil

	case "thumbnailImage":
		return stream.ThumbnailImage, nil

	case "criteria":
		return stream.Criteria.GetPath(p.Tail())

	default:
		return stream.Data[property], nil
	}
}

// SetPath implements the path.Setter interface.  It takes any data value
// and tries to set it to the correct path within this Stream.
func (stream *Stream) SetPath(p path.Path, value interface{}) error {

	if p.IsEmpty() {
		return derp.New(500, "ghost.model.Stream", "Unrecognized path", p)
	}

	property := p.Head()

	// Properties that can be set
	switch property {

	case "label":
		stream.Label = convert.String(value)

	case "description":
		stream.Description = convert.String(value)

	case "thumbnailImage":
		stream.ThumbnailImage = convert.String(value)

	case "criteria":
		return stream.Criteria.SetPath(p.Tail(), value)

	default:
		stream.Data[property] = value
	}

	return nil
}

func (stream *Stream) Roles(authorization *Authorization) []string {

	result := make([]string, 0)

	if stream.Criteria.Public {
		result = append(result, "PUBLIC")
	}

	if authorization == nil {
		return result
	}

	if !stream.AuthorID.IsZero() {
		if authorization.UserID == stream.AuthorID {
			result = append(result, "AUTHOR")
		}
	}

	if authorization.DomainOwner {
		result = append(result, "OWNER")
	}

	if roles, ok := stream.Criteria.Roles[authorization.UserID]; ok {
		result = append(result, roles...)
	}

	for _, groupID := range authorization.GroupIDs {

		if roles, ok := stream.Criteria.Roles[groupID]; ok {
			result = append(result, roles...)
		}
	}

	return result
}
