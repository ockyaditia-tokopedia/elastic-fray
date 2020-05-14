package officialclient

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"

	elasticEntity "github.com/elastic-fray/entity/elastic"
	"github.com/elastic-fray/entity/promo/marketplace"
	"github.com/elastic-fray/pkg/elastic/officialclient"

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
		elastic, err := officialclient.New(officialclient.Config{
			Config: c.Config,
		})
		if err != nil {
			log.Error(err)
			return
		}

		m = Module{
			config:  c.Config,
			monitor: c.Monitor,
			usecase: Usecase{
				elastic: elastic,
			},
		}
	})

	return m, err
}

func (m Module) GetInfo(ctx context.Context, o ...func(*esapi.InfoRequest)) (*esapi.Response, error) {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.officialclient.get.info", nil)

	return m.usecase.elastic.GetInfo(ctx)
}

func (m Module) GetPromoOrderUsage(ctx context.Context, parameter elasticEntity.ElasticSearchParameter, o ...func(*esapi.SearchRequest)) ([]marketplace.Promo, error) {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.officialclient.get.promo.order.usage", nil)

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

	if err := m.usecase.elastic.ProcessSearch(ctx, &elastic.SearchOption{
		URL:         m.config.ElasticSearch.URL,
		Label:       "promo.order.usage",
		Index:       elastic.ConstElasticSearchIndexPromoOrderUsage,
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

func (m Module) CountPromoOrderUsage(ctx context.Context, parameter elasticEntity.ElasticSearchParameter, o ...func(*esapi.SearchRequest)) (int, error) {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.officialclient.count.promo.order.usage", nil)

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

	total, err := m.usecase.elastic.ProcessCount(ctx, &elastic.SearchOption{
		URL:         m.config.ElasticSearch.URL,
		Label:       "promo.order.usage",
		Index:       elastic.ConstElasticSearchIndexPromoOrderUsage,
		Input:       req,
		Environment: true,
		PreferNode:  parameter.PreferNode,
	})
	if err != nil {
		log.Error(err)
	}

	return total, err
}

func (m Module) InsertPromoOrderUsage(ctx context.Context, req marketplace.Promo) error {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.officialclient.insert.promo.order.usage", nil)

	err := m.usecase.elastic.ProcessInsert(ctx, &elastic.InsertOption{
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
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.officialclient.update.promo.order.usage", nil)

	err := m.usecase.elastic.ProcessUpdate(ctx, &elastic.InsertOption{
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

func (m Module) DeletePromoOrderUsage(ctx context.Context, id string) (string, error) {
	defer m.monitor.SetHistogram(time.Now(), "usecase.elastic.officialclient.delete.promo.order.usage", nil)

	resp, err := m.usecase.elastic.ProcessDelete(ctx, id, &elastic.DeleteOption{
		URL:         m.config.ElasticSearch.URL,
		Environment: true,
		Index:       elastic.ConstElasticSearchIndexPromoOrderUsage,
	})
	if err != nil {
		log.Error(err)
	}

	return resp, err
}
