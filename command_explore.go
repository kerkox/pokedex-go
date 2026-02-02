package main

import (
	"errors"
	"fmt"
)




func commandExplore(cfg *config, args ...string) error {
	// Parameter validation
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}

	locationName := args[0]
	location, err := cfg.pokeapiClient.GetLocation(locationName)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found Pokemon: ")
	for _, encounter := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
