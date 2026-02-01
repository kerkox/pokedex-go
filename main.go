package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

type Config struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type CachedResponse struct {
	Next     string                   `json:"next"`
	Previous string                   `json:"previous"`
	Results  []map[string]interface{} `json:"results"`
}

// LocationAreaResponse - Estructura específica para el endpoint location-area/{name}
// Solo definimos los campos que necesitamos (Interface Segregation Principle)
type LocationAreaResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
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
		// ✅ Actualizar config desde cache
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

	// Actualizar config
	config.Next = cached.Next
	config.Previous = cached.Previous

	// Cachear respuesta completa con navegación
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
		// ✅ Actualizar config desde cache
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

	// Actualizar config
	config.Next = cached.Next
	config.Previous = cached.Previous

	// Cachear respuesta completa con navegación
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
	// Validación de parámetros
	if len(params) == 0 || params[0] == "" {
		fmt.Println("Please provide a location to explore.")
		return nil
	}
	location := params[0]

	PokedexApiURLLocationArea := POKEDEX_API_URL + LOCATION_AREA_ENDPOINT + location + "/"
	fmt.Printf("Exploring %s...\n", location)

	// Intentar obtener del cache primero
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

	// Hacer la petición HTTP
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

	// Leer el body completo para poder cachearlo
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return err
	}

	// Deserializar a nuestra estructura específica
	var locationResp LocationAreaResponse
	if err := json.Unmarshal(body, &locationResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return err
	}

	// Cachear los bytes crudos (más eficiente y reutilizable)
	cache.Add(PokedexApiURLLocationArea, body)

	// Imprimir los nombres de los pokemon
	fmt.Println("Found Pokemon:")
	printPokemonEncounters(locationResp)

	return nil
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


