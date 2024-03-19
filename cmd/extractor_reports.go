package main

import (
	"github.com/CesarDelgadoM/extractor-reports/config"
	"github.com/CesarDelgadoM/extractor-reports/server"
)

func main() {
	// Config
	load := config.LoadConfig("config-local.yml")
	config := config.ParseConfig(load)

	// Server
	server := server.NewServer(config)

	server.Run()
}
