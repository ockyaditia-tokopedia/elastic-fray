package officialclient

import (
	"bytes"
	"context"
	"encoding/json"
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
		m.elastic.Search.WithContext(ctx),
		m.elastic.Search.WithIndex(so.Index),
		m.elastic.Search.WithBody(&buffer),
		m.elastic.Search.WithTrackTotalHits(true),
		m.elastic.Search.WithPretty(),
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
		return err
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
