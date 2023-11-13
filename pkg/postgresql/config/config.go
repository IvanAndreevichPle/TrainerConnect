package config

import (
	"encoding/json"
	"os"
)

// ReadConfig читает конфигурацию из файла и возвращает объект DBConfig
func ReadConfig(filename string) (DBConfig, error) {
	var config DBConfig

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
