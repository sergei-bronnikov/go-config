package config

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

/*
	type Config struct {
		Mode              string `env:"APP_MODE" default:"production"`
		Postgres          PostgresConfig
		AdminServerPort   int    `env:"ADMIN_SERVER_PORT" default:"8080"`
	}
*/
func Load(conf *interface{}) error {
	err := prepareConfig(reflect.ValueOf(conf).Elem())
	if err != nil {
		return err
	}
	return nil
}

func prepareConfig(obj reflect.Value) error {
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		if field.Kind() == reflect.Struct {
			err := prepareConfig(field)
			if err != nil {
				return err
			}
			continue
		}
		val := getConfigValue(obj.Type().Field(i))
		switch field.Kind() {
		case reflect.String:
			field.SetString(val)
		case reflect.Int:
			intVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return errors.New("config parse process. invalid int value")
			}
			field.SetInt(intVal)
		default:
			return errors.New("unsupported type in config")
		}
	}
	return nil
}

func getConfigValue(field reflect.StructField) string {
	envKey, ok := field.Tag.Lookup("env")
	if ok {
		envVal := os.Getenv(envKey)
		if envVal != "" {
			return envVal
		}
	}
	return field.Tag.Get("default")
}
