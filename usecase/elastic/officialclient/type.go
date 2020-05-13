package officialclient

import (
	"context"

	"github.com/elastic/go-elasticsearch/v7/esapi"

	elasticEntity "github.com/elastic-fray/entity/elastic"
	"github.com/elastic-fray/entity/promo/marketplace"
	"github.com/elastic-fray/pkg/monitor"
	"github.com/elastic-fray/pkg/utils"

	"github.com/tokopedia/sauron/src/elastic"
)

type (
	Method interface { // TODO: should using own param, avoid external param
		GetInfo(ctx context.Context, o ...func(*esapi.InfoRequest)) (*esapi.Response, error)
		GetPromoOrderUsage(ctx context.Context, parameter elasticEntity.ElasticSearchParameter, o ...func(*esapi.SearchRequest)) ([]marketplace.Promo, error)
	}

	ElasticMethod interface { // TODO: should using own param, avoid external param
		GetInfo(ctx context.Context, o ...func(*esapi.InfoRequest)) (*esapi.Response, error)
		ProcessSearch(ctx context.Context, so *elastic.SearchOption, o ...func(*esapi.SearchRequest)) error
	}
)

type (
	Config struct {
		Config  utils.Config
		Monitor monitor.Method
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
