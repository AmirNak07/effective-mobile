package config

import pkgconfig "effective-mobile/pkg/config"

type Config struct {
	Postgres pkgconfig.PostgresConfig
}

func MustLoad() *Config {
	return pkgconfig.MustLoad[Config](".env")
}
