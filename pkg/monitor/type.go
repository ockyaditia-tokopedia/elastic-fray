package monitor

import (
	"time"

	"github.com/ooyala/go-dogstatsd"
)

type (
	Method interface {
		SetHistogram(start time.Time, name string, tags []string)
		SetCount(name string, tags []string)
	}

	DatadogMethod interface {
		SetHistogram(start time.Time, name string, tags []string)
		SetCount(name string, tags []string)
	}
)

type (
	Config struct {
		Datadog *dogstatsd.Client
	}

	Module struct {
		datadog DatadogMethod
	}
)
