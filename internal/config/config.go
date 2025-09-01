package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	AIProvider  string
	OpenAIKey   string
	OpenAIModel string
	AIFakeMode  bool
}

func mustGet(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if fallback != "" {
		return fallback
	}
	log.Fatalf("missing required env var: %s", key)
	return ""
}

func FromEnv() Config {
	return Config{
		Port:        mustGet("PORT", "8080"),
		DatabaseURL: mustGet("DATABASE_URL", "appuser:apppass@tcp(db:3306)/ai_triage?parseTime=true"),
		AIProvider:  mustGet("AI_PROVIDER", "gemini"),
		OpenAIKey:   os.Getenv("GEMINI_API_KEY"),
		OpenAIModel: mustGet("GEMINI_MODEL", "gemini-1.5-flash"),
		AIFakeMode:  os.Getenv("AI_FAKE_MODE") == "false",
	}
}
