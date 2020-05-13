package elastic

import (
	"time"

	"github.com/elastic-fray/entity/promo/marketplace"
)

type (
	ElasticSearchParameter struct {
		QueryString string
		Size        int64
		IsUsingTime bool
		GTE         time.Time
		LTE         time.Time
		Source      string
		Sort        map[string]interface{}
		PreferNode  string
	}

	PromoOrderUsage struct {
		Took     int  `json:"took"`
		TimedOut bool `json:"timed_out"`
		Shards   struct {
			Total      int `json:"total"`
			Successful int `json:"successful"`
			Failed     int `json:"failed"`
		} `json:"_shards"`
		Hits struct {
			Total struct {
				Value    int64  `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
			MaxScore float64 `json:"max_score"`
			Hits     []struct {
				Index  string            `json:"_index"`
				Type   string            `json:"_type"`
				ID     string            `json:"_id"`
				Score  float64           `json:"_score"`
				Source marketplace.Promo `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
)
