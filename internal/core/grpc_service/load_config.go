package grpc_service

import (
	"os"

	"github.com/Ulqiora/Route256Project/internal/config"
)

func LoadConfig(configFilePath string) *config.Config {
	file, err := os.Open(configFilePath)
	if err != nil {
		panic("file not founded")
		return nil
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	return config.New(file)
}
