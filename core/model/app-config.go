package model

type Config struct {
	Concurrency ConcurrencyConfig
	Listen      ListenConfig
	Proxy       ProxyConfig
}

type ListenConfig struct {
	Addr string
}

type ProxyConfig struct {
	File string
}

type ConcurrencyConfig struct {
	MaxConcurrency int
}
