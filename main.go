package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name		string
	description	string
	callback    func(config *Config) error
}

type Config struct {
	Next string `json:"next"`
	Previous string `json:"previous"`
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
} 

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, cmd := range registry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMapBack(config *Config) error {
	PokedexApiURL := config.Previous
	if PokedexApiURL == "" {
		fmt.Printf("you're on the first page\n")
		return nil
	}
	res, err := http.Get(PokedexApiURL)	
	if err != nil {
		fmt.Printf("Error fetching previous map data: %v\n", err)
		return err
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	
	if res.StatusCode > 299 {
		fmt.Printf("Error: received status code %d\n", res.StatusCode)
		return fmt.Errorf("received status code %d", res.StatusCode)
	}
	var result map[string]interface{}
	err = decoder.Decode(&result)
	if err != nil && err != io.EOF {
		fmt.Printf("Error decoding previous map data: %v\n", err)	
		return err
	}
	
	if next, ok := result["next"].(string); ok {
		config.Next = next
	} else {
		config.Next = ""
	}
	
	if previous, ok := result["previous"].(string); ok {
		config.Previous = previous
	} else {
		config.Previous = ""
	}
	
	for _, value := range result["results"].([]interface{}) {
		item := value.(map[string]interface{})
		fmt.Printf("%s\n", item["name"])
	}
	// Placeholder for actual map display logic
	return nil
}

func commandMap(config *Config) error {

	PokedexApiURL := config.Next
	if PokedexApiURL == "" {
		fmt.Printf("No more map data to fetch.\n")
	}
	fmt.Printf("Fetching map data from %s\n", PokedexApiURL)
	res, err := http.Get(PokedexApiURL)	
	if err != nil {
		fmt.Printf("Error fetching map data: %v\n", err)
		return err
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	
	if res.StatusCode > 299 {
		fmt.Printf("Error: received status code %d\n", res.StatusCode)
		return fmt.Errorf("received status code %d", res.StatusCode)
	}
	var result map[string]interface{}
	err = decoder.Decode(&result)
	if err != nil && err != io.EOF {
		fmt.Printf("Error decoding map data: %v\n", err)	
		return err
	}
	
	if next, ok := result["next"].(string); ok {
		config.Next = next
	} else {
		config.Next = ""
	}
	
	if previous, ok := result["previous"].(string); ok {
		config.Previous = previous
	} else {
		config.Previous = ""
	}
	
	fmt.Println("Map Data:")
	for _, value := range result["results"].([]interface{}) {
		item := value.(map[string]interface{})
		fmt.Printf("%s\n", item["name"])
	}
	// Placeholder for actual map display logic
	return nil
}

var registry map[string]cliCommand
var POKEDEX_API_URL = os.Getenv("POKEDEX_API_URL")


var config Config

func init() {
	if POKEDEX_API_URL == "" {
		POKEDEX_API_URL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}
	config = Config{
		Next: POKEDEX_API_URL,
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
			err := command.callback(&config)
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


