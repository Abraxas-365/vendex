package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Abraxas-365/hada-commerce/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration from environment variables.
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Connect to PostgreSQL.
	db, err := openDB(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer db.Close()

	// Build the DI container — wires all domain containers.
	ctr, err := NewContainer(db, cfg)
	if err != nil {
		return fmt.Errorf("building container: %w", err)
	}

	// Build Fiber app with all routes registered.
	app := newApp(ctr)

	// Listen for OS shutdown signals.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("received %v, initiating graceful shutdown...", sig)
		cancel()
	}()

	return runServer(ctx, app, ":"+cfg.Port)
}

// openDB opens a PostgreSQL connection pool and verifies connectivity.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Connection pool tuning.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}
