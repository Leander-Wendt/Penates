package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Leander-Wendt/Penates/backend/internal/config"
	"github.com/Leander-Wendt/Penates/backend/internal/database"
	"github.com/Leander-Wendt/Penates/backend/internal/handler"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-healthcheck" {
		runHealthcheck()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "-seed" {
		runSeed()
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := database.BootstrapAdmin(db, cfg.SeedAdminEmail, cfg.SeedAdminPassword, cfg.SeedAdminOrg); err != nil {
		log.Fatalf("failed to bootstrap admin user: %v", err)
	}

	router := handler.NewRouter(cfg, db)

	if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
		log.Printf("Penates backend listening on port %s (HTTPS)", cfg.Port)
		err = router.RunTLS(":"+cfg.Port, cfg.TLSCertFile, cfg.TLSKeyFile)
	} else {
		log.Printf("Penates backend listening on port %s (HTTP)", cfg.Port)
		err = router.Run(":" + cfg.Port)
	}
	if err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

// runHealthcheck performs a local GET /healthz request and exits 0 on a 200
// response or 1 otherwise. It exists so a container healthcheck can probe the
// server without a shell or curl, which the distroless runtime image lacks.
func runHealthcheck() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	scheme := "http"
	client := http.Client{Timeout: 3 * time.Second}
	if os.Getenv("TLS_CERT_FILE") != "" && os.Getenv("TLS_KEY_FILE") != "" {
		scheme = "https"
		// The healthcheck talks to the server's own self-signed cert over
		// loopback, so there's no third party to be fooled by skipping
		// verification here.
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} //nolint:gosec
	}

	resp, err := client.Get(fmt.Sprintf("%s://127.0.0.1:%s/healthz", scheme, port))
	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
	os.Exit(0)
}

// runSeed connects to the configured database, migrates it, and populates it
// with dummy Organisations, Users, Items, and LoanRequests for local
// development and demos. It exits 1 on failure.
func runSeed() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := database.BootstrapAdmin(db, cfg.SeedAdminEmail, cfg.SeedAdminPassword, cfg.SeedAdminOrg); err != nil {
		log.Fatalf("failed to bootstrap admin user: %v", err)
	}

	if err := database.SeedDummyData(db); err != nil {
		log.Fatalf("failed to seed dummy data: %v", err)
	}

	log.Println("database seeded with dummy data")
}
