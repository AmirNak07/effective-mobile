package config

import pkgconfig "effective-mobile/pkg/config"

type Config struct {
	Env      string
	Postgres pkgconfig.PostgresConfig
}

func MustLoad() *Config {
	return pkgconfig.MustLoadFromEnv[Config]()
}
