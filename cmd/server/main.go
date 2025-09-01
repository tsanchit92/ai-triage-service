package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/yourname/ai-triage/docs"
	"github.com/yourname/ai-triage/internal/ai"
	"github.com/yourname/ai-triage/internal/config"
	"github.com/yourname/ai-triage/internal/db"
	"github.com/yourname/ai-triage/internal/httpx"
	"github.com/yourname/ai-triage/internal/incidents"
)

const maxRetries = 10
const retryInterval = 60 * time.Second // try every 10 seconds
var sqlxDB *sqlx.DB
var err error

// @title Incident Management API
// @version 1.0
// @description API for managing incidents with AI classification
// @BasePath /
func main() {

	cfg := config.FromEnv()
	doMigrate := flag.Bool("migrate", true, "run migrations and exit")
	flag.Parse()

	ctx := context.Background()
	for i := 1; i <= maxRetries; i++ {
		sqlxDB, err = db.Connect(ctx, cfg.DatabaseURL)
		if err == nil {
			break
		}
		log.Printf("db connect attempt %d/%d failed: %v, retrying in %v...", i, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
	}
	if err != nil {
		log.Fatalf("could not connect to db after %d attempts: %v", maxRetries, err)
	}
	defer sqlxDB.Close()

	if *doMigrate {
		if err := db.Migrate(ctx, sqlxDB.DB, "internal/migrations"); err != nil {
			log.Fatalf("migrate: %v", err)
		}
		log.Println("migrations applied")
	}

	// AI client selection
	var aiClient ai.Client
	if cfg.AIFakeMode {
		aiClient = ai.FakeClient{Next: ai.Classification{Severity: "Medium", Category: "Software"}}
	} else {
		switch cfg.AIProvider {
		case "openai":
			aiClient = ai.NewGeminiClient(cfg.OpenAIKey, cfg.OpenAIModel)
		case "gemini":
			aiClient = ai.NewGeminiClient(cfg.OpenAIKey, cfg.OpenAIModel)
		default:
			log.Fatalf("unknown AI_PROVIDER: %s", cfg.AIProvider)
		}
	}

	repo := incidents.NewRepository(sqlxDB)
	svc := incidents.NewService(aiClient, repo)
	h := &incidents.Handler{Svc: svc}

	r := httpx.New()
	incidents.RegisterRoutes(r, h)

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
