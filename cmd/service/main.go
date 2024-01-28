package main

import (
	"WB_Tech_level_0/internal/app"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	configPath = "configs/app.yaml"
)

func main() {
	config := app.NewConfig()
	file, err := os.ReadFile(configPath)

	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}

	s := app.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
