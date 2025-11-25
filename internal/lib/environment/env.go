package environment

import (
	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func LoadEnvironment[T any](dotenvPaths ...string) (T, error) {
	out := new(T)

	godotenv.Load(dotenvPaths...)

	err := env.Parse(out)

	if err != nil {
		return *out, err
	}

	err = validator.New(validator.WithRequiredStructEnabled()).Struct(out)

	if err != nil {
		return *out, err
	}

	return *out, nil
}
