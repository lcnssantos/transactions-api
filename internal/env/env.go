package env

type Environment struct {
	Environment         string `env:"ENVIRONMENT" validate:"required"`
	DatabaseHost        string `env:"DB_HOST" validate:"required"`
	DatabaseUser        string `env:"DB_USER" validate:"required"`
	DatabasePass        string `env:"DB_PASS" validate:"required"`
	DatabasePort        string `env:"DB_PORT" validate:"required"`
	DatabaseSSLMode     string `env:"DB_SSL_MODE" validate:"required"`
	DatabaseName        string `env:"DB_NAME" validate:"required"`
	DatabasePoolMinimum int    `env:"DB_POOL_MINIMUM" validate:"required"`
	DatabasePoolMaximum int    `env:"DB_POOL_MAXIMUM" validate:"required"`
	Port                string `env:"PORT" validate:"required"`
}
