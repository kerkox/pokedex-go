package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	pokecache "github.com/kerkox/pokedex-cli-go/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, params []string) error
}

const LOCATION_AREA_ENDPOINT = "location-area/"
const POKEMON_ENDPOINT = "pokemon/"
var pokemonsCaught = map[string]Pokemon{}

type Config struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type CachedResponse struct {
	Next     string                   `json:"next"`
	Previous string                   `json:"previous"`
	Results  []map[string]interface{} `json:"results"`
}

// LocationAreaResponse - Specific structure for the location-area/{name} endpoint
// We only define the fields we need (Interface Segregation Principle)
type LocationAreaResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

// Pokemon represents both the API response and the domain model.
// In Go, when the structure is identical, there's no need to create separate types.
// If in the future you need additional domain fields (e.g., DateCaught),
// you can create a CaughtPokemon type that embeds Pokemon.
type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

func commandExit(config *Config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	cache.Stop()
	return nil
} 

func commandHelp(config *Config, params []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, cmd := range registry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMapBack(config *Config, params []string) error {
	PokedexApiURL := config.Previous
	if PokedexApiURL == "" {
		fmt.Printf("you're on the first page\n")
		return nil
	}

	if cachedData, found := cache.Get(PokedexApiURL); found {
		fmt.Println("Using cached previous map data:")
		var cached CachedResponse
		err := json.Unmarshal(cachedData, &cached)
		if err != nil {
			fmt.Printf("Error decoding cached previous map data: %v\n", err)
			return err
		}
		// ✅ Update config from cache
		config.Next = cached.Next
		config.Previous = cached.Previous
		
		for _, item := range cached.Results {
			if name, ok := item["name"].(string); ok {
				fmt.Printf("%s\n", name)
			}
		}
		return nil
	}


	res, err := http.Get(PokedexApiURL)
	if err != nil {
		fmt.Printf("Error fetching previous map data: %v\n", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Printf("Error: received status code %d\n", res.StatusCode)
		return fmt.Errorf("received status code %d", res.StatusCode)
	}

	var cached CachedResponse
	err = json.NewDecoder(res.Body).Decode(&cached)
	if err != nil && err != io.EOF {
		fmt.Printf("Error decoding previous map data: %v\n", err)
		return err
	}

	// Update config
	config.Next = cached.Next
	config.Previous = cached.Previous

	// Cache complete response with navigation
	data, err := json.Marshal(cached)
	if err != nil {
		fmt.Printf("Error marshaling for cache: %v\n", err)
	} else {
		cache.Add(PokedexApiURL, data)
	}

	for _, item := range cached.Results {
		if name, ok := item["name"].(string); ok {
			fmt.Printf("%s\n", name)
		}
	}
	return nil
}

func commandMap(config *Config, params []string) error {
	PokedexApiURL := config.Next
	if PokedexApiURL == "" {
		fmt.Printf("No more map data to fetch.\n")
		return nil
	}

	if cachedData, found := cache.Get(PokedexApiURL); found {
		fmt.Println("Using cached map data:")
		var cached CachedResponse
		err := json.Unmarshal(cachedData, &cached)
		if err != nil {
			fmt.Printf("Error decoding cached map data: %v\n", err)
			return err
		}
		// ✅ Update config from cache
		config.Next = cached.Next
		config.Previous = cached.Previous

		for _, item := range cached.Results {
			if name, ok := item["name"].(string); ok {
				fmt.Printf("%s\n", name)
			}
		}
		return nil
	}



	fmt.Printf("Fetching map data from %s\n", PokedexApiURL)
	res, err := http.Get(PokedexApiURL)
	if err != nil {
		fmt.Printf("Error fetching map data: %v\n", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Printf("Error: received status code %d\n", res.StatusCode)
		return fmt.Errorf("received status code %d", res.StatusCode)
	}

	var cached CachedResponse
	err = json.NewDecoder(res.Body).Decode(&cached)
	if err != nil && err != io.EOF {
		fmt.Printf("Error decoding map data: %v\n", err)
		return err
	}

	// Update config
	config.Next = cached.Next
	config.Previous = cached.Previous

	// Cache complete response with navigation
	data, err := json.Marshal(cached)
	if err != nil {
		fmt.Printf("Error marshaling for cache: %v\n", err)
	} else {
		cache.Add(PokedexApiURL, data)
	}

	fmt.Println("Map Data:")
	for _, item := range cached.Results {
		if name, ok := item["name"].(string); ok {
			fmt.Printf("%s\n", name)
		}
	}
	return nil
}

func commandExplore(config *Config, params []string) error {
	// Parameter validation
	if len(params) == 0 || params[0] == "" {
		fmt.Println("Please provide a location to explore.")
		return nil
	}
	location := params[0]

	PokedexApiURLLocationArea := POKEDEX_API_URL + LOCATION_AREA_ENDPOINT + location + "/"
	fmt.Printf("Exploring %s...\n", location)

	// Try to get from cache first
	if cachedData, found := cache.Get(PokedexApiURLLocationArea); found {
		fmt.Println("Using cached data:")
		var locationResp LocationAreaResponse
		if err := json.Unmarshal(cachedData, &locationResp); err != nil {
			fmt.Printf("Error decoding cached data: %v\n", err)
			return err
		}
		printPokemonEncounters(locationResp)
		return nil
	}

	// Make HTTP request
	fmt.Printf("Fetching data from %s\n", PokedexApiURLLocationArea)
	res, err := http.Get(PokedexApiURLLocationArea)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Printf("Error: received status code %d\n", res.StatusCode)
		return fmt.Errorf("received status code %d", res.StatusCode)
	}

	// Read full body to cache it
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return err
	}

	// Deserialize to our specific structure
	var locationResp LocationAreaResponse
	if err := json.Unmarshal(body, &locationResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return err
	}

	// Cache raw bytes (more efficient and reusable)
	cache.Add(PokedexApiURLLocationArea, body)

	// Print pokemon names
	fmt.Println("Found Pokemon:")
	printPokemonEncounters(locationResp)

	return nil
}

// getPokemonFromCache attempts to get a Pokemon from cache.
// Returns nil, nil if not in cache (not an error).
// Follows SRP: only responsible for getting data from cache.
func getPokemonFromCache(url string) (*Pokemon, error) {
	cachedData, found := cache.Get(url)
	if !found {
		return nil, nil
	}
	
	var pokemon Pokemon
	if err := json.Unmarshal(cachedData, &pokemon); err != nil {
		return nil, fmt.Errorf("error decoding cached Pokemon: %w", err)
	}
	return &pokemon, nil
}

// ErrPokemonNotFound indicates that the Pokemon doesn't exist in the API
var ErrPokemonNotFound = fmt.Errorf("pokemon not found")

// fetchPokemonFromAPI fetches a Pokemon from the API and caches it.
// Returns ErrPokemonNotFound if the Pokemon doesn't exist.
// Follows SRP: only responsible for fetch + cache.
func fetchPokemonFromAPI(url string) (*Pokemon, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching Pokemon: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrPokemonNotFound
	}
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("API error: status code %d", res.StatusCode)
	}

	// Read full body to cache it
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return nil, fmt.Errorf("error decoding Pokemon: %w", err)
	}

	// Cache raw bytes
	cache.Add(url, body)

	return &pokemon, nil
}


// GetPokemon fetches a Pokemon, first from cache, then from the API.
// This is the public function that orchestrates data retrieval.
func GetPokemon(pokemonName string) (*Pokemon, error) {
	url := POKEDEX_API_URL + POKEMON_ENDPOINT + pokemonName + "/"

	// Try cache first
	if pokemon, err := getPokemonFromCache(url); err != nil {
		return nil, err
	} else if pokemon != nil {
		return pokemon, nil
	}

	// Cache miss - fetch from API
	return fetchPokemonFromAPI(url)
}

func commandCatch(config *Config, params []string) error {
	// Parameter validation
	if len(params) == 0 || params[0] == "" {
		fmt.Println("Please specify a Pokemon to catch.")
		return nil
	}
	pokemonName := params[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := GetPokemon(pokemonName)
	if err != nil {
		if err == ErrPokemonNotFound {
			fmt.Printf("Pokemon %s not found!\n", pokemonName)
			return nil
		}
		return err
	}

	// Centralized catch logic
	if attemptCatch(pokemon) {
		fmt.Printf("%s was caught!\n", pokemonName)
		pokemonsCaught[pokemonName] = *pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}

// attemptCatch determines if a Pokemon was caught.
// Catch probability is inversely proportional to BaseExperience.
// Pokemon with higher experience are harder to catch.
func attemptCatch(pokemon *Pokemon) bool {
	// Avoid division by zero
	if pokemon.BaseExperience <= 0 {
		return true
	}
	return rand.Intn(pokemon.BaseExperience) == 0
}

// printPokemonEncounters - Helper function siguiendo DRY principle
func printPokemonEncounters(resp LocationAreaResponse) {
	for _, encounter := range resp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
}

var registry map[string]cliCommand
var POKEDEX_API_URL = os.Getenv("POKEDEX_API_URL")


var config Config
var cacheDuration = 10 // seconds
var cache = pokecache.NewCache(time.Duration(cacheDuration) * time.Second)

func init() {
	if POKEDEX_API_URL == "" {
		POKEDEX_API_URL = "https://pokeapi.co/api/v2/"
	}
	config = Config{
		Next: POKEDEX_API_URL+LOCATION_AREA_ENDPOINT+"?offset=0&limit=20",
		Previous: "",
	}
	registry = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous map",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore the Pokedex (alias for map)",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon (not implemented yet)",
			callback:    commandCatch,
		},
	}
}




func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nPokedex > ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		parts := cleanInput(line)
		if len(parts) == 0 {
			continue
		}
		commandName := parts[0]
		if command, exists := registry[commandName]; exists {
			err := command.callback(&config, parts[1:])
			if err != nil {
				fmt.Printf("Error executing command '%s': %v\n", commandName, err)
			}
			if commandName == "exit" {
				break
			}
		} else {
			fmt.Printf("Unknown command: %s\n", commandName)
		}
		// fmt.Printf("Your command was: %s", parts[0])
	}
	// fmt.Printf("%q",cleanInput("   hello  world "))
}


