package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ooyala/go-dogstatsd"

	"github.com/elastic-fray/entity/elastic"
	"github.com/elastic-fray/entity/promo/marketplace"
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
	// currentTime := time.Now()

	// for i := currentTime.Hour(); i < 13; {
	// Elastic API
	processElasticAPI()

	// Elastic Official Client
	processElasticOfficialClient()
	// }
}

func processElasticAPI() {
	defer Monitor.SetHistogram(time.Now(), "handler.elastic.api.get.promo.order.usage", nil)

	elasticAPI := api.New(api.Config{
		Config:   Config,
		Datadog:  DatadogClient,
		Location: Location,
		Monitor:  Monitor,
	})

	searchResp, err := elasticAPI.GetPromoOrderUsage(Context, elastic.ElasticSearchParameter{
		QueryString: "source:marketplace",
		Source:      "api.benchmark",
	})
	if err != nil {
		log.Error(err)
	}

	fmt.Println("API Search - Total Result: ", len(searchResp))

	countResp, err := elasticAPI.CountPromoOrderUsage(Context, "source:marketplace")
	if err != nil {
		log.Error(err)
	}

	fmt.Println("API Count - Total Result: ", countResp)

	if err = elasticAPI.InsertPromoOrderUsage(Context, marketplace.Promo{
		OrderID: 69696969,
	}); err != nil {
		log.Error(err)
	}

	if err = elasticAPI.UpdatePromoOrderUsage(Context, marketplace.Promo{
		OrderID: 69696969,
	}); err != nil {
		log.Error(err)
	}

	time.Sleep(1000000000) // 1s, let give it time

	deleteResp, err := elasticAPI.DeletePromoOrderUsage(Context, "order_id:69696969")
	if err != nil {
		log.Error(err)
	}

	fmt.Println("API Delete - Status: ", deleteResp)
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

	searchResp, err := elasticOfficial.GetPromoOrderUsage(Context, elastic.ElasticSearchParameter{
		QueryString: "source:marketplace",
		Source:      "officialclient.benchmark",
	})
	if err != nil {
		log.Error(err)
	}

	fmt.Println("Official Client Search - Total Result: ", len(searchResp))

	countResp, err := elasticOfficial.CountPromoOrderUsage(Context, elastic.ElasticSearchParameter{
		QueryString: "source:marketplace",
		Source:      "officialclient.benchmark",
	})
	if err != nil {
		log.Error(err)
	}

	fmt.Println("Official Client Count - Total Result: ", countResp)

	if err = elasticOfficial.InsertPromoOrderUsage(Context, marketplace.Promo{
		OrderID: 96969696,
	}); err != nil {
		log.Error(err)
	}

	if err = elasticOfficial.UpdatePromoOrderUsage(Context, marketplace.Promo{
		OrderID: 96969696,
	}); err != nil {
		log.Error(err)
	}

	time.Sleep(1000000000) // 1s, let give it time

	deleteResp, err := elasticOfficial.DeletePromoOrderUsage(Context, "96969696")
	if err != nil {
		log.Error(err)
	}

	fmt.Println("Official Client Delete - Status:", deleteResp)
}
