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
	}

	ElasticMethod interface {
		Search(so *elastic.SearchOption) error
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
