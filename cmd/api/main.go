package main

import (
	"log"
	"net/http"
	"time"

	"beach/internal/adapters/cache"
	"beach/internal/adapters/catalog"
	adapter "beach/internal/adapters/http"
	"beach/internal/adapters/ipma"
	"beach/internal/adapters/stream"
	"beach/internal/application"
	"beach/internal/config"
)

func main() {
	cfg := config.Load()

	// Static catalog — cameras + stream slugs embedded in binary
	repo, err := catalog.New()
	if err != nil {
		log.Fatalf("loading camera catalog: %v", err)
	}
	log.Printf("loaded %d cameras from static catalog", countCameras(repo))

	// Stream resolver — builds HLS URL directly from embedded slugs
	resolver := stream.NewResolver()

	// Conditions resolver — IPMA public APIs
	httpClient := &http.Client{Timeout: cfg.HTTPClientTimeout}
	ipmaCache := cache.New(cfg.CatalogCacheTTL)
	ipmaClient := ipma.NewClient(httpClient, ipmaCache)
	conditions := ipma.NewConditionsResolver(ipmaClient)

	service := application.NewCameraService(repo, resolver, conditions)

	handler := adapter.NewHandler(service)
	router := adapter.NewRouter(handler)

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("%s starting on port %s (env=%s)", cfg.AppName, cfg.AppPort, cfg.AppEnv)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func countCameras(repo *catalog.StaticCatalog) int {
	cameras, _ := repo.ListCameras()
	return len(cameras)
}
