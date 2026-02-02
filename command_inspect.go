package main

import (
	"errors"
	"fmt"

	"github.com/kerkox/pokedex-cli-go/pokeapi"
)



func commandInspect(cfg *config, args ...string) error {
	// Parameter validation
	if len(args) != 1 {
		return errors.New("You must provide a pokemon name")
	}
	
	pokemonName := args[0]

	pokemon, ok := cfg.caughtPokemon[pokemonName]
	if !ok {
		return errors.New("You have not caught that pokemon")
	}

	PrintPokemonDetails(&pokemon)
	return nil
}



func PrintPokemonDetails(pokemon *pokeapi.Pokemon) {
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pType := range pokemon.Types {
		fmt.Printf("  - %s\n", pType.TypeInfo.Name)
	}
}