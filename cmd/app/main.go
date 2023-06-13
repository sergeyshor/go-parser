package main

import (
	"fmt"
	"go-parser/config"
	"log"
	"go-parser/internal/app"
	"time"
)

func main() {
	start := time.Now()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)

	duration := time.Since(start)
	fmt.Println(duration)
}