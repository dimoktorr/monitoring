package repository

import "time"

type Config struct {
	DSN             string        `env:"DSN"`
	MaxConnIdleTime time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"30m"`
}
