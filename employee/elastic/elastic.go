package elastic

import (
	"bytes"
	"context"
	conf "employee/config"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

var (
	countSuccessful uint64
	indexName       = "employee-management"
)

// PostDataInSearch pushes data to Elasticsearch
func PostDataInSearch(c conf.Configuration, id string, data interface{}, ctxReq context.Context) {
	esClient, err := generateElasticClient(c)
	if err != nil {
		logrus.Errorf("Unable to create elastic client: %v", err)
		return
	}

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,
		Client:        esClient,
		NumWorkers:    1,
		FlushBytes:    int(5e+6),
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		logrus.Errorf("Error creating bulk indexer: %v", err)
		return
	}

	if indexExists(c, indexName) != 200 {
		res, err := esClient.Indices.Create(indexName)
		if err != nil {
			logrus.Errorf("Error creating index: %v", err)
			return
		}
		defer res.Body.Close()
	}

	putDataInSearch(data, bulkIndexer, id, ctxReq)

	if err := bulkIndexer.Close(ctxReq); err != nil {
		logrus.Errorf("Error closing bulk indexer: %v", err)
	}

	logrus.Infof("Employee data indexed successfully")
}

func putDataInSearch(jsonData interface{}, bulkIndexer esutil.BulkIndexer, id string, ctxReq context.Context) {
	data, err := json.Marshal(jsonData)
	if err != nil {
		logrus.Errorf("JSON marshal failed: %v", err)
		return
	}

	err = bulkIndexer.Add(
		ctxReq,
		esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: id,
			Body:       bytes.NewReader(data),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				atomic.AddUint64(&countSuccessful, 1)
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				logrus.Errorf("Bulk indexing failed: %v", err)
			},
		},
	)

	if err != nil {
		logrus.Errorf("Bulk add failed: %v", err)
	}
}

func indexExists(c conf.Configuration, index string) int {
	esClient, err := generateElasticClient(c)
	if err != nil {
		logrus.Errorf("Elastic client error: %v", err)
		return 500
	}

	resp, err := esClient.Indices.Exists([]string{index})
	if err != nil {
		logrus.Errorf("Index check failed: %v", err)
		return 404
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

// SearchDataInElastic searches by employee ID
func SearchDataInElastic(c conf.Configuration, id string, ctxReq context.Context) map[string]interface{} {
	var buf bytes.Buffer
	result := make(map[string]interface{})

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": id,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logrus.Errorf("Query encode error: %v", err)
		return result
	}

	es, err := generateElasticClient(c)
	if err != nil {
		logrus.Errorf("Elastic client error: %v", err)
		return result
	}

	res, err := es.Search(
		es.Search.WithContext(ctxReq),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		logrus.Errorf("Search failed: %v", err)
		return result
	}
	defer res.Body.Close()

	if res.IsError() {
		logrus.Errorf("Search returned error status: %s", res.Status())
		return result
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logrus.Errorf("Response decode error: %v", err)
	}

	return result
}

// SearchALLDataInElastic returns all documents
func SearchALLDataInElastic(c conf.Configuration, ctxReq context.Context) map[string]interface{} {
	result := make(map[string]interface{})

	es, err := generateElasticClient(c)
	if err != nil {
		logrus.Errorf("Elastic client error: %v", err)
		return result
	}

	res, err := es.Search(
		es.Search.WithContext(ctxReq),
		es.Search.WithIndex(indexName),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		logrus.Errorf("Search failed: %v", err)
		return result
	}
	defer res.Body.Close()

	if res.IsError() {
		logrus.Errorf("Search returned error status: %s", res.Status())
		return result
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logrus.Errorf("Response decode error: %v", err)
	}

	return result
}

// CheckElasticHealth checks Elasticsearch availability
func CheckElasticHealth(c conf.Configuration, ctxReq context.Context) (bool, error) {
	es, err := generateElasticClient(c)
	if err != nil {
		return false, err
	}

	res, err := es.Info(es.Info.WithContext(ctxReq))
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}
