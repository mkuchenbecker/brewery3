package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/mkuchenbecker/brewery3/data/datasink"
	"github.com/mkuchenbecker/brewery3/data/gomodel/data"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type firestoreClient struct {
	client *firestore.Client
}

func NewFirestoreClient(client *firestore.Client) datasink.FirestoreClient {
	return &firestoreClient{client: client}
}

func (c *firestoreClient) Send(ctx context.Context, collection string, docName string, doc map[string]interface{}) error {
	_, err := c.client.Collection(collection).Doc(docName).Set(ctx, doc)
	return err
}

func (c *firestoreClient) Get(ctx context.Context, collection string, id string) (datasink.RowList, error) {
	iter := c.client.Collection(collection).
		Where("id", "==", id).
		Documents(ctx)

	var lom datasink.RowList = make([]datasink.ColValueMap, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		lom = append(lom, doc.Data())
	}
	return lom, nil
}

func NewStore(collection string, client datasink.FirestoreClient) datasink.DataSink {
	return &firestoreSink{collection: collection, clock: &Clock{}, client: client}
}

type firestoreSink struct {
	collection string
	clock      datasink.Clock
	client     datasink.FirestoreClient
}

func (s *firestoreSink) Send(ctx context.Context, in *data.DataObject) (*data.SendResponse, error) {
	row := make(map[string]interface{})
	for key, val := range in.Fields {
		switch t := val.Value.(type) {
		case *data.Value_Bool:
			row[key] = t.Bool
		case *data.Value_Bytes:
			row[key] = t.Bytes
		case *data.Value_Double:
			row[key] = t.Double
		case *data.Value_Float:
			row[key] = t.Float
		case *data.Value_Int32:
			row[key] = t.Int32
		case *data.Value_Int64:
			row[key] = t.Int64
		case *data.Value_String_:
			row[key] = t.String_
		case *data.Value_Uint32:
			row[key] = t.Uint32
		case *data.Value_Uint64:
			row[key] = t.Uint64
		default:
			return &data.SendResponse{}, errors.New("bad data")
		}
	}

	return &data.SendResponse{}, errors.Wrap(s.client.Send(ctx, s.collection, in.Key, row), "error saving to database")
}

func (s *firestoreSink) Get(ctx context.Context, in *data.GetRequest) (*data.GetResponse, error) {
	response := &data.GetResponse{Data: []*data.DataObject{}}
	lom, err := s.client.Get(ctx, s.collection, in.Key)
	if err != nil {
		return response, err
	}

	do := &data.DataObject{Fields: make(map[string]*data.Value)}

	do.Key = in.Key

	for _, mapStringInterface := range lom {
		for k, v := range mapStringInterface {
			switch x := v.(type) {
			case bool:
				do.Fields[k] = &data.Value{Value: &data.Value_Bool{Bool: x}}
			case []byte:
				do.Fields[k] = &data.Value{Value: &data.Value_Bytes{Bytes: x}}
			case float64:
				do.Fields[k] = &data.Value{Value: &data.Value_Double{Double: x}}
			case float32:
				do.Fields[k] = &data.Value{Value: &data.Value_Float{Float: x}}
			case int32:
				do.Fields[k] = &data.Value{Value: &data.Value_Int32{Int32: x}}
			case int64:
				do.Fields[k] = &data.Value{Value: &data.Value_Int64{Int64: x}}
			case string:
				do.Fields[k] = &data.Value{Value: &data.Value_String_{String_: x}}
			case uint32:
				do.Fields[k] = &data.Value{Value: &data.Value_Uint32{Uint32: x}}
			case uint64:
				do.Fields[k] = &data.Value{Value: &data.Value_Uint64{Uint64: x}}
			}
		}
		response.Data = append(response.Data, do)
	}

	return response, nil
}

type Clock struct{}

func (c *Clock) Now() time.Time {
	return time.Now()
}
