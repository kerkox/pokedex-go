package pokeapi

import "os"

// Default values as private constants (compile-time safe)
const (
	defaultBaseURL   = "https://pokeapi.co/api/v2/"
	PokemonEndpoint  = "pokemon/"
	LocationEndpoint = "location-area/"
)

// BaseURL is the configured API base URL.
// Reads from POKEDEX_API_URL environment variable, falls back to default.
var BaseURL string

func init() {
	BaseURL = os.Getenv("POKEDEX_API_URL")
	if BaseURL == "" {
		BaseURL = defaultBaseURL
	}
}