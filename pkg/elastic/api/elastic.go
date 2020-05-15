package api

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/tokopedia/tdk/go/log"

	"github.com/tokopedia/sauron/src/elastic"
	"github.com/tokopedia/sauron/src/httpclient"
	"github.com/tokopedia/sauron/src/utils"
)

var (
	once sync.Once
	m    Module
)

func New(c Config) Method {
	once.Do(func() {
		elastic := elastic.New(c.Datadog, &utils.GConfig{
			Server: utils.ServerConfig{
				Environment: c.Config.Server.Environment,
			},
			ElasticSearch: utils.ElasticSearchConfig{
				Sauron: c.Config.ElasticSearch.URL,
			},
		}, c.Location)
		if elastic == nil {
			log.Fatal("Elastic API client nil")
			return
		}

		m = Module{
			config:  c.Config,
			elastic: elastic,
		}
	})

	return m
}

func (m Module) Search(ctx context.Context, so *elastic.SearchOption) error {
	return m.elastic.Search(so)
}

func (m Module) Count(ctx context.Context, so *elastic.SearchOption) (int, error) {
	return m.elastic.Count(so)
}

func (m Module) Insert(ctx context.Context, io *elastic.InsertOption) error {
	return m.elastic.Insert(io)
}

func (m Module) Update(ctx context.Context, io *elastic.InsertOption) error {
	return m.elastic.Update(io)
}

func (m Module) Delete(ctx context.Context, do *elastic.DeleteOption) (elastic.ElasticSearchDeleteResponse, error) {
	return m.elastic.Delete(do)
}

func (m Module) Bulk(ctx context.Context, url, input string) (bool, error) {
	type response struct {
		Created bool `json:"created"`
	}

	var resp response

	agent := httpclient.NewHTTPRequest()
	agent.Url = url
	agent.Path = "/_bulk"
	agent.Method = "POST"
	agent.IsJson = true
	agent.Json = input

	body, err := agent.DoReq()
	if err != nil {
		log.Error(err)
		agent.Debug()
		return resp.Created, err
	}

	if err := json.Unmarshal(*body, &resp); err != nil {
		log.Error(err)
		return resp.Created, err
	}

	agent.Debug()
	return resp.Created, err
}
