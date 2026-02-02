package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

func commandCatch(cfg *config, args ...string) error {
	// Parameter validation
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	
	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	// Centralized catch logic
	if attemptCatch(pokemon.BaseExperience) {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		cfg.caughtPokemon[pokemon.Name] = pokemon	
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	
	return nil
}



// attemptCatch determines if a Pokemon was caught.
// Catch probability is inversely proportional to BaseExperience.
// Pokemon with higher experience are harder to catch.
func attemptCatch(pokemonBaseExperience int) bool {
	// Avoid division by zero
	if pokemonBaseExperience <= 0 {
		return true
	}
	oportunity := int(math.Round(float64(pokemonBaseExperience) / 10))

	return rand.Intn(oportunity) == 0
}