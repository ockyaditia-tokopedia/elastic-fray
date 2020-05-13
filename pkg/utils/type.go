package utils

type (
	Config struct {
		Server        ServerConfig
		Datadog       DatadogConfig
		ElasticSearch ElasticSearchConfig
	}

	ServerConfig struct {
		Environment string
	}

	DatadogConfig struct {
		Connection string
	}

	ElasticSearchConfig struct {
		URL string
	}
)
