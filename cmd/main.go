package main

import (
	"balancer-wrote-myself/internal/balancer"
	"balancer-wrote-myself/internal/config"
	"net/http"
	"os"

	"github.com/hashicorp/go-hclog"
)

func main() {
	log := hclog.Default()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Error("Error to load config", "config", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	lb := balancer.NewLoadBalacner(log, cfg.Backends)

	mux.Handle("GET /api/{path...}", lb)

	log.Info("The server was opened")
	if err := http.ListenAndServe(cfg.Server.Port, mux); err != nil {
		log.Error("Error to start server", "error", err)
		os.Exit(1)
	}
}
