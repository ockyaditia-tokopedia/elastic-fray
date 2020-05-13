package monitor

import (
	"time"

	"github.com/tokopedia/tdk/go/log"

	"github.com/tokopedia/sauron/src/logging/monitor"
)

func New(c Config) Method {
	return Module{
		datadog: monitor.New(monitor.Config{
			Client: c.Datadog,
		}),
	}
}

func (m Module) SetHistogram(start time.Time, name string, tags []string) {
	if m.datadog == nil {
		log.Error("Empty monitor datadog for metric: ", name)
		return
	}

	m.datadog.SetHistogram(start, name, tags)
}

func (m Module) SetCount(name string, tags []string) {
	if m.datadog == nil {
		log.Error("Empty monitor datadog for metric: ", name)
		return
	}

	m.datadog.SetCount(name, tags)
}
