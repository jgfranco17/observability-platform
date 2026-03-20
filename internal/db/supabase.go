package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/jgfranco17/observability-platform/internal/environment"
	supabase "github.com/supabase-community/supabase-go"
)

var (
	loadClientOnce = sync.OnceValues(loadSupabaseClient)
)

// GetSupabaseClient returns a singleton Supabase client
func GetSupabaseClient() (*supabase.Client, error) {
	return loadClientOnce()
}

// loadSupabaseClient initializes and returns a Supabase client.
// Serves as a wrapper for the Once sync to ensure singleton behavior.
func loadSupabaseClient() (*supabase.Client, error) {
	databaseURL := os.Getenv(environment.ENV_KEY_DB_URL)
	databaseKey := os.Getenv(environment.ENV_KEY_DB_KEY)
	if databaseURL == "" || databaseKey == "" {
		return nil, fmt.Errorf("both database URL and private key must be set in environment.")
	}

	client, err := supabase.NewClient(databaseURL, databaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Supabase client: %v", err)
	}
	return client, nil
}
