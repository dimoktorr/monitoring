package storage

import "time"

type Config struct {
	Hosts    string        `env:"HOSTS"`
	TTL      time.Duration `env:"TTL" envDefault:"10s"`
	Password string        `env:"PASSWORD"`
}
