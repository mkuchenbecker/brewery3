package datasink

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type ColValueMap map[string]interface{}

type RowList []ColValueMap

type FirestoreClient interface {
	Send(ctx context.Context, collection string, docName string, doc map[string]interface{}) error
	Get(ctx context.Context, collection string, id string) (RowList, error)
}

type firestoreClient struct {
	client *firestore.Client
}

func NewFirestoreClient(client *firestore.Client) FirestoreClient {
	return &firestoreClient{client: client}
}

func (c *firestoreClient) Send(ctx context.Context, collection string, docName string, doc map[string]interface{}) error {
	_, err := c.client.Collection(collection).Doc(docName).Set(ctx, doc)
	return err
}

func (c *firestoreClient) Get(ctx context.Context, collection string, id string) (RowList, error) {
	iter := c.client.Collection(collection).
		Where("id", "==", id).
		Documents(ctx)

	var lom RowList = make([]ColValueMap, 0)
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
