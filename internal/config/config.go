package config

import (
	"fmt"
	"reflect"
	"strings"
)

type Validator interface {
	Validate() error
}

type Config struct {
	DB  DBConfig
	JWT JWTConfig
}
type DBConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

type JWTConfig struct {
	Secret string
}

func (db DBConfig) Validate() error {
	return validateStruct(db)
}

func (jwt JWTConfig) Validate() error {
	return validateStruct(jwt)
}

func validateStruct(s any) error {
	v := reflect.ValueOf(s)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.String {
			if strings.TrimSpace(f.String()) == "" {
				return fmt.Errorf("%s is required", t.Field(i).Name)
			}
		} else if f.IsZero() {
			return fmt.Errorf("%s is required", t.Field(i).Name)
		}
	}

	return nil
}
