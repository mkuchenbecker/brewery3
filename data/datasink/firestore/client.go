// +build integration

package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/mkuchenbecker/brewery3/data/datasink"
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
