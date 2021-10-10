package model

type Config struct {
	Concurrency ConcurrencyConfig
	Tls         TlsConfig
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

type TlsConfig struct {
	Cert string
	Key  string
}
