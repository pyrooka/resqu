package db

import (
	"context"
	"encoding/json"
	"errors"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// BigQuery DB
type BigQuery struct {
	client *bigquery.Client
}

func init() {
	registerDb("bigquery", &BigQuery{})
}

// Init DB interface implementation.
func (bq *BigQuery) Init(config map[string]interface{}) (err error) {
	ctx := context.Background()

	var projectID, serviceAccPath string

	// Check the necessary variables.
	if pID, ok := config["projectId"]; ok {
		projectID = pID.(string)
	} else {
		return errors.New("projectId is missing from the connection config")
	}
	if svcAccPath, ok := config["serviceAccPath"]; ok {
		serviceAccPath = svcAccPath.(string)
	} else {
		return errors.New("serviceAccPath is missing from the connection config")
	}

	bq.client, err = bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(serviceAccPath))

	return
}

// Run DB interface implementation.
func (bq *BigQuery) Run(ctx context.Context, query string) (result []byte, err error) {
	bqQuery := bq.client.Query(query)

	res, err := bqQuery.Read(ctx)
	if err != nil {
		return
	}

	result, err = buildResult(res)

	return
}

func buildResult(res *bigquery.RowIterator) (result []byte, err error) {
	resultRows := []map[string]bigquery.Value{}
	schema := res.Schema

	for {
		var row []bigquery.Value
		err = res.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}

		rowMap := map[string]bigquery.Value{}

		for i, val := range row {
			fieldName := schema[i].Name
			rowMap[fieldName] = val
		}

		resultRows = append(resultRows, rowMap)
	}

	result, err = json.Marshal(resultRows)

	return
}
