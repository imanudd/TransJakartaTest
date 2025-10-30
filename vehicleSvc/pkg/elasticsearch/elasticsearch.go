package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const (
	BOOK_DETAILS = "book_details"
)

type ElasticsearchImpl interface {
	Save(ctx context.Context, index string, document interface{}) error
}

type elasticsearch struct {
	client *elastic.Client
}

func New(client *elastic.Client) ElasticsearchImpl {
	return &elasticsearch{
		client: client,
	}
}

func (es *elasticsearch) Save(ctx context.Context, index string, document interface{}) error {
	data, err := json.Marshal(document)
	if err != nil {
		return errors.New("error when marshaling document")
	}

	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error indexing document: %s", res.String())
		return fmt.Errorf("error when indexing document : %s", res.String())
	}

	log.Println("Document indexed successfully")

	return nil
}
