package config

type IConfig interface {
	BusConfig() BusConfig
}

type BusConfig struct {
	RmqConnection            string
	PropagateContextMetadata bool
	RetryCount               int
}
