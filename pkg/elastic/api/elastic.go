package api

import (
	"context"
	"sync"

	"github.com/tokopedia/tdk/go/log"

	"github.com/tokopedia/sauron/src/elastic"
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
