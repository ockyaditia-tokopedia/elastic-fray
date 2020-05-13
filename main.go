package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ooyala/go-dogstatsd"

	"github.com/elastic-fray/entity/elastic"
	"github.com/elastic-fray/pkg/monitor"
	"github.com/elastic-fray/pkg/utils"
	"github.com/elastic-fray/usecase/elastic/api"
	"github.com/elastic-fray/usecase/elastic/officialclient"

	"github.com/tokopedia/tdk/go/log"
)

var (
	Context       context.Context
	Config        utils.Config
	DatadogClient *dogstatsd.Client
	Location      *time.Location
	Monitor       monitor.Method
	err           error
)

func init() {
	Context = context.Background()

	if err = log.SetConfig(&log.Config{
		Level:   "info",
		AppName: "elastic-fray",
		Caller:  true,
	}); err != nil {
		log.Fatal(err)
	}

	Config = utils.Config{
		Server: utils.ServerConfig{
			Environment: "", // TODO: please fill
		},
		Datadog: utils.DatadogConfig{
			Connection: "", // TODO: please fill
		},
		ElasticSearch: utils.ElasticSearchConfig{
			URL: "", // TODO: please fill
		},
	}

	DatadogClient, err = dogstatsd.New(Config.Datadog.Connection)
	if err != nil {
		log.Fatal(err)
	}
	DatadogClient.Namespace = "elastic-fray."
	DatadogClient.Tags = append(DatadogClient.Tags, "env:"+Config.Server.Environment)

	Location, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal(err)
	}

	Monitor = monitor.New(monitor.Config{
		Datadog: DatadogClient,
	})
}

func main() {
	// Elastic API
	processElasticAPI()

	// Elastic Official Client
	processElasticOfficialClient()
}

func processElasticAPI() {
	defer Monitor.SetHistogram(time.Now(), "handler.elastic.api.get.promo.order.usage", nil)

	elasticAPI := api.New(api.Config{
		Config:   Config,
		Datadog:  DatadogClient,
		Location: Location,
		Monitor:  Monitor,
	})

	resp, err := elasticAPI.GetPromoOrderUsage(Context, elastic.ElasticSearchParameter{
		QueryString: "source:marketplace",
		Source:      "api.benchmark",
	})
	if err != nil {
		log.Error(err)
	}

	fmt.Println("Total Result: ", len(resp))
}

func processElasticOfficialClient() {
	defer Monitor.SetHistogram(time.Now(), "handler.elastic.official.client.get.promo.order.usage", nil)

	elasticOfficial, err := officialclient.New(officialclient.Config{
		Config:  Config,
		Monitor: Monitor,
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := elasticOfficial.GetPromoOrderUsage(Context, elastic.ElasticSearchParameter{
		QueryString: "source:marketplace",
		Source:      "officialclient.benchmark",
	})
	if err != nil {
		log.Error(err)
	}

	fmt.Println("Total Result: ", len(resp))
}
