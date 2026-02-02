# Pokedex CLI

A command-line Pokedex application built in Go as part of the [Boot.dev](https://boot.dev) backend development course.

## About

This project is a hands-on practice exercise to learn and apply core Go concepts including:

- **HTTP Clients** - Making requests to external REST APIs (PokeAPI)
- **JSON Parsing** - Encoding/decoding JSON responses with typed structs
- **Caching** - Implementing an in-memory cache with TTL expiration
- **REPL** - Building an interactive command-line interface
- **Concurrency** - Using goroutines and mutexes for thread-safe caching
- **Clean Architecture** - Separating concerns with dedicated packages
- **Environment Variables** - Configurable API endpoints with fallback defaults

## Features

| Command | Description |
|---------|-------------|
| `help` | Display available commands |
| `map` | Display the next 20 location areas |
| `mapb` | Display the previous 20 location areas |
| `explore <location>` | List all Pokemon in a location area |
| `catch <pokemon>` | Attempt to catch a Pokemon |
| `inspect <pokemon>` | View details of a caught Pokemon |
| `exit` | Exit the Pokedex |

## Installation

```bash
# Clone the repository
git clone https://github.com/kerkox/pokedex-cli-go.git
cd pokedex-cli-go

# Build the project
go build -o pokedex

# Run the application
./pokedex
```

## Usage

```bash
$ ./pokedex

Pokedex > help
Welcome to the Pokedex!
Usage:
exit: Exit the Pokedex
help: Display a help message
map: Display the map
mapb: Display the previous map
explore: Explore the Pokedex
catch: Catch a Pokemon
inspect: Inspect a caught Pokemon

Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...

Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 - magikarp
 - gyarados
 ...

Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!

Pokedex > inspect pikachu
Name: pikachu
Height: 4
Weight: 60
Stats:
  - hp: 35
  - attack: 55
  - defense: 40
  - special-attack: 50
  - special-defense: 50
  - speed: 90
Types:
  - electric

Pokedex > exit
Closing the Pokedex... Goodbye!
```

## Project Structure

```
pokedex/
├── main.go                 # Application entry point
├── repl.go                 # REPL and command routing
├── repl_test.go            # Unit tests for REPL
├── command_help.go         # Help command
├── command_exit.go         # Exit command
├── command_map.go          # Map navigation commands
├── command_explore.go      # Location exploration command
├── command_catch.go        # Pokemon catching command
├── command_inspect.go      # Pokemon inspection command
├── go.mod                  # Go module definition
├── .env.example            # Environment variables example
├── internal/
│   ├── pokecache.go        # Cache implementation with TTL
│   └── pokecache_test.go   # Cache tests
└── pokeapi/
    ├── pokeapi.go          # API configuration and constants
    ├── client.go           # HTTP client with caching
    ├── location_list.go    # List locations endpoint
    ├── location_get.go     # Get location details endpoint
    ├── pokemon_get.go      # Get Pokemon endpoint
    ├── types_locations.go  # Location response types
    └── types_pokemon.go    # Pokemon response types
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────────────────────────────────────────────────┤
│  │ REPL (repl.go) + Commands (command_*.go)                │
│  └─────────────────────────────────────────────────────────┤
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Client Layer                         │
│  ┌─────────────────────────────────────────────────────────┤
│  │ pokeapi/ - HTTP Client with typed requests/responses    │
│  └─────────────────────────────────────────────────────────┤
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                     │
│  ┌─────────────────────────────────────────────────────────┤
│  │ internal/pokecache - Thread-safe cache with TTL         │
│  └─────────────────────────────────────────────────────────┤
└─────────────────────────────────────────────────────────────┘
```

## API

This project uses the [PokeAPI](https://pokeapi.co/) - a free RESTful API for Pokemon data.

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `POKEDEX_API_URL` | Base URL for the Pokemon API | `https://pokeapi.co/api/v2/` |

### Setup

```bash
# Copy the example environment file
cp .env.example .env

# Edit as needed
vim .env
```

## Key Concepts Practiced

### API Client Design
- Dedicated `Client` struct with HTTP client and cache
- Typed request/response structs for each endpoint
- Separation of concerns between different API resources

### Caching Strategy
- In-memory cache stores API responses as raw bytes
- TTL-based expiration with background cleanup goroutine
- Thread-safe operations using `sync.Mutex`

### Clean Code Principles
- **SRP**: Each command in its own file
- **DRY**: Shared API client logic in `pokeapi` package
- **ISP**: Structs only define fields they need from the API

### Error Handling
- Descriptive error messages for user feedback
- Graceful handling of API failures and edge cases

## License

This project is for educational purposes as part of the Boot.dev curriculum.

---

Built with ❤️ while learning Go at [Boot.dev](https://boot.dev)
