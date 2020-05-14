package api

import (
	"context"
	"time"

	"github.com/ooyala/go-dogstatsd"

	elasticEntity "github.com/elastic-fray/entity/elastic"
	"github.com/elastic-fray/entity/promo/marketplace"
	"github.com/elastic-fray/pkg/monitor"
	"github.com/elastic-fray/pkg/utils"

	"github.com/tokopedia/sauron/src/elastic"
)

type (
	Method interface {
		GetPromoOrderUsage(ctx context.Context, parameter elasticEntity.ElasticSearchParameter) ([]marketplace.Promo, error)
		CountPromoOrderUsage(ctx context.Context, query string) (int, error)
		InsertPromoOrderUsage(ctx context.Context, req marketplace.Promo) error
		UpdatePromoOrderUsage(ctx context.Context, req marketplace.Promo) error
		DeletePromoOrderUsage(ctx context.Context, query string) (int, error)
	}

	ElasticMethod interface { // TODO: should using own param, avoid external param
		Search(ctx context.Context, so *elastic.SearchOption) error
		Count(ctx context.Context, so *elastic.SearchOption) (int, error)
		Insert(ctx context.Context, io *elastic.InsertOption) error
		Update(ctx context.Context, io *elastic.InsertOption) error
		Delete(ctx context.Context, do *elastic.DeleteOption) (elastic.ElasticSearchDeleteResponse, error)
	}
)

type (
	Config struct {
		Config   utils.Config
		Datadog  *dogstatsd.Client
		Location *time.Location
		Monitor  monitor.Method
	}

	Usecase struct {
		elastic ElasticMethod
	}

	Module struct {
		config  utils.Config
		monitor monitor.Method
		usecase Usecase
	}
)
