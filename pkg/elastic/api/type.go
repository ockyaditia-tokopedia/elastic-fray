package api

import (
	"context"
	"time"

	"github.com/ooyala/go-dogstatsd"

	"github.com/elastic-fray/pkg/utils"

	"github.com/tokopedia/sauron/src/elastic"
)

type (
	Method interface { // TODO: should using own param, avoid external param
		Search(ctx context.Context, so *elastic.SearchOption) error
		Count(ctx context.Context, so *elastic.SearchOption) (int, error)
		Insert(ctx context.Context, io *elastic.InsertOption) error
		Update(ctx context.Context, io *elastic.InsertOption) error
		Delete(ctx context.Context, do *elastic.DeleteOption) (elastic.ElasticSearchDeleteResponse, error)
	}

	ElasticMethod interface {
		Search(so *elastic.SearchOption) error
		Count(so *elastic.SearchOption) (int, error)
		Insert(io *elastic.InsertOption) error
		Update(io *elastic.InsertOption) error
		Delete(do *elastic.DeleteOption) (elastic.ElasticSearchDeleteResponse, error)
	}
)

type (
	Config struct {
		Config   utils.Config
		Datadog  *dogstatsd.Client
		Location *time.Location
	}

	Module struct {
		config  utils.Config
		elastic ElasticMethod
	}
)
