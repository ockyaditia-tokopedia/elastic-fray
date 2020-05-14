package officialclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"

	"github.com/tokopedia/tdk/go/log"

	"github.com/tokopedia/sauron/src/elastic"
)

var (
	once sync.Once
	m    Module
)

func New(c Config) (Method, error) {
	var err error

	once.Do(func() {
		elastic, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{
				c.Config.ElasticSearch.URL,
			},
		})
		if err != nil {
			log.Error(err)
			return
		}

		m = Module{
			config:  c.Config,
			elastic: elastic,
		}
	})

	return m, err
}

func (m Module) GetInfo(ctx context.Context, o ...func(*esapi.InfoRequest)) (*esapi.Response, error) {
	return m.elastic.Info()
}

func (m Module) ProcessSearch(ctx context.Context, so *elastic.SearchOption, o ...func(*esapi.SearchRequest)) error {
	var (
		buffer bytes.Buffer
		size   int64
		result map[string]interface{}
	)

	if so.Size > 0 {
		size = so.Size
	} else {
		size = elastic.MAX_ELASTIC_SIZE
	}

	esq := elastic.ElasticSearchQuery{
		From:  int64(0),
		Size:  size,
		Query: so.Input,
	}

	if so.Sort != nil {
		esq.Sort = so.Sort
	}

	if so.Environment == true {
		if m.config.Server.Environment == "development" {
			m.config.Server.Environment = "staging"
		}

		so.Index = m.config.Server.Environment + "-" + so.Index
	}

	if err := json.NewEncoder(&buffer).Encode(esq); err != nil {
		log.Error(err)
		return err
	}

	resp, err := m.elastic.Search(
		// m.elastic.Search.WithContext(ctx),
		m.elastic.Search.WithIndex(so.Index),
		m.elastic.Search.WithBody(&buffer),
		// m.elastic.Search.WithTrackTotalHits(true),
		// m.elastic.Search.WithPretty(),
	)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			log.Errorf("Error parsing the response body: %v", err)
		} else {
			log.Errorf("[%s] %v: %v",
				resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return errors.New("Error")
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error(err)
		return err
	}

	body, err := json.Marshal(result)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := json.Unmarshal(body, so.Output); err != nil {
		log.Error(err)
	}

	return err
}

func (m Module) ProcessCount(ctx context.Context, so *elastic.SearchOption, o ...func(*esapi.CountRequest)) (int, error) {
	var (
		buffer bytes.Buffer
		result map[string]interface{}
	)

	esq := elastic.ElasticSearchQuery{
		Query: so.Input,
	}

	if so.Environment == true {
		if m.config.Server.Environment == "development" {
			m.config.Server.Environment = "staging"
		}

		so.Index = m.config.Server.Environment + "-" + so.Index
	}

	if err := json.NewEncoder(&buffer).Encode(esq); err != nil {
		log.Error(err)
		return 0, err
	}

	resp, err := m.elastic.Count(
		// m.elastic.Count.WithContext(ctx),
		m.elastic.Count.WithIndex(so.Index),
		m.elastic.Count.WithBody(&buffer),
		// m.elastic.Count.WithPretty(),
	)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			log.Errorf("Error parsing the response body: %v", err)
		} else {
			log.Error(e)
		}

		return 0, errors.New("Error")
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error(err)
		return 0, err
	}

	return int(result["count"].(float64)), err
}

func (m Module) ProcessInsert(ctx context.Context, so *elastic.InsertOption) error {
	if so.Environment == true {
		if m.config.Server.Environment == "development" {
			m.config.Server.Environment = "staging"
		}

		so.Index = m.config.Server.Environment + "-" + so.Index
	}

	body, err := json.Marshal(so.Data)
	if err != nil {
		log.Error(err)
		return err
	}

	req := esapi.IndexRequest{
		Index:      so.Index,
		DocumentID: so.ID,
		Body:       strings.NewReader(string(body)),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, m.elastic)
	if err != nil {
		log.Error(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Errorf("[%s] Error indexing document ID=%d", res.Status(), so.ID)

		return errors.New("Error")
	} else {
		var r map[string]interface{}

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	return err
}

func (m Module) ProcessUpdate(ctx context.Context, so *elastic.InsertOption) error {
	if so.Environment == true {
		if m.config.Server.Environment == "development" {
			m.config.Server.Environment = "staging"
		}

		so.Index = m.config.Server.Environment + "-" + so.Index
	}

	body, err := json.Marshal(so.Data)
	if err != nil {
		log.Error(err)
		return err
	}

	req := esapi.IndexRequest{
		Index:      so.Index,
		DocumentID: so.ID,
		Body:       strings.NewReader(string(body)),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, m.elastic)
	if err != nil {
		log.Error(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Errorf("[%s] Error indexing document ID=%d", res.Status(), so.ID)

		return errors.New("Error")
	} else {
		var r map[string]interface{}

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	return err
}

func (m Module) ProcessDelete(ctx context.Context, id string, so *elastic.DeleteOption, o ...func(*esapi.DeleteRequest)) (string, error) {
	var result map[string]interface{}

	if so.Environment == true {
		if m.config.Server.Environment == "development" {
			m.config.Server.Environment = "staging"
		}

		so.Index = m.config.Server.Environment + "-" + so.Index
	}

	resp, err := m.elastic.Delete(
		so.Index,
		id,
	)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			log.Errorf("Error parsing the response body: %v", err)
		} else {
			log.Error(e)
		}

		return "", errors.New("Error")
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error(err)
		return "", err
	}

	return result["result"].(string), err
}
