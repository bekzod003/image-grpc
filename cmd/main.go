package main

import (
	"github.com/bekzod003/image-grpc/config"
	"github.com/bekzod003/image-grpc/internal/app"
)

func main() {
	// Run the server
	cfg := config.GetConfig()
	app.Run(cfg)
}
