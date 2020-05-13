package api

import (
	"context"
	"sync"
	"time"

	elasticEntity "github.com/elastic-fray/entity/elastic"
	"github.com/elastic-fray/entity/promo/marketplace"
	"github.com/elastic-fray/pkg/elastic/api"

	"github.com/tokopedia/tdk/go/log"

	"github.com/tokopedia/sauron/src/elastic"
)

var (
	once sync.Once
	m    Module
)

func New(c Config) Method {
	once.Do(func() {
		m = Module{
			config:  c.Config,
			monitor: c.Monitor,
			usecase: Usecase{
				elastic: api.New(api.Config{
					Config:   c.Config,
					Datadog:  c.Datadog,
					Location: c.Location,
				}),
			},
		}
	})

	return m
}

func (m Module) GetPromoOrderUsage(ctx context.Context, parameter elasticEntity.ElasticSearchParameter) ([]marketplace.Promo, error) {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.api.get.promo.order.usage", nil)

	var (
		promos []marketplace.Promo
		resp   elasticEntity.PromoOrderUsage
	)

	req := elastic.Query{
		Bool: &elastic.Bool{
			Must: []elastic.Must{
				elastic.Must{
					QueryString: map[string]interface{}{
						"query": parameter.QueryString,
					},
				},
			},
		},
	}

	if parameter.IsUsingTime {
		req.Bool.Must = append(req.Bool.Must, elastic.Must{
			Range: map[string]interface{}{
				"create_time": map[string]interface{}{
					"gte":       parameter.GTE.Format("2006-01-02"),
					"lte":       parameter.LTE.Format("2006-01-02"),
					"format":    "yyyy-MM-dd",
					"time_zone": "+07:00",
				},
			},
		})
	}

	if err := m.usecase.elastic.Search(ctx, &elastic.SearchOption{
		URL:         m.config.ElasticSearch.URL,
		Label:       "promo.order.usage",
		Index:       elastic.ConstElasticSearchIndexPromoOrderUsage,
		Type:        "",
		Input:       req,
		Environment: true,
		Output:      &resp,
		Size:        parameter.Size,
		Sort:        parameter.Sort,
		PreferNode:  parameter.PreferNode,
	}); err != nil {
		log.Error(err)
		return promos, err
	}

	for _, hit := range resp.Hits.Hits {
		promos = append(promos, hit.Source)
	}

	return promos, nil
}
