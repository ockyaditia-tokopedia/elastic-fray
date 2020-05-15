package api

import (
	"context"
	"strconv"
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
		Environment: true,
		Label:       "promo.order.usage",
		Index:       elastic.ConstElasticSearchIndexPromoOrderUsage,
		Input:       req,
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

func (m Module) CountPromoOrderUsage(ctx context.Context, query string) (int, error) {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.api.count.promo.order.usage", nil)

	total, err := m.usecase.elastic.Count(ctx, &elastic.SearchOption{
		URL:         m.config.ElasticSearch.URL,
		Environment: true,
		Label:       "promo.order.usage",
		Index:       elastic.ConstElasticSearchIndexPromoOrderUsage,
		Input: elastic.Query{
			QueryString: map[string]interface{}{
				"query": query,
			},
		},
		PreferNode: elastic.ConstPreferNodeTypeDefault,
	})
	if err != nil {
		log.Error(err)
	}

	return total, err
}

func (m Module) InsertPromoOrderUsage(ctx context.Context, req marketplace.Promo) error {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.api.insert.promo.order.usage", nil)

	err := m.usecase.elastic.Insert(ctx, &elastic.InsertOption{
		URL:         m.config.ElasticSearch.URL,
		Environment: true,
		Index:       elastic.ConstElasticSearchIndexPromoOrderUsage,
		Type:        "order",
		ID:          strconv.FormatInt(req.OrderID, 10),
		Data:        req,
	})
	if err != nil {
		log.Error(err)
	}

	return err
}

func (m Module) UpdatePromoOrderUsage(ctx context.Context, req marketplace.Promo) error {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.api.update.promo.order.usage", nil)

	err := m.usecase.elastic.Update(ctx, &elastic.InsertOption{
		URL:         m.config.ElasticSearch.URL,
		Environment: true,
		Index:       elastic.ConstElasticSearchIndexPromoOrderUsage,
		Type:        "order",
		ID:          strconv.FormatInt(req.OrderID, 10),
		Data:        req,
	})
	if err != nil {
		log.Error(err)
	}

	return err
}

func (m Module) DeletePromoOrderUsage(ctx context.Context, query string) (int, error) {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.api.delete.promo.order.usage", nil)

	resp, err := m.usecase.elastic.Delete(ctx, &elastic.DeleteOption{
		URL:         m.config.ElasticSearch.URL,
		Environment: true,
		Index:       elastic.ConstElasticSearchIndexPromoOrderUsage,
		Type:        "order",
		Query: elastic.Query{
			Bool: &elastic.Bool{
				Should: []elastic.Should{
					{
						Bool: &elastic.Bool{
							Must: []elastic.Must{
								{
									QueryString: map[string]interface{}{
										"query": query,
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Error(err)
	}

	return resp.Deleted, err
}

func (m Module) BulkPromoOrderUsage(ctx context.Context, url, input string) (bool, error) {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.api.bulk.promo.order.usage", nil)

	resp, err := m.usecase.elastic.Bulk(ctx, url, input)
	if err != nil {
		log.Error(err)
	}

	return resp, err
}
