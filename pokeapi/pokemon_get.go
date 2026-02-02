package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

// fetchPokemonFromAPI fetches a Pokemon from the API and caches it.
// Returns ErrPokemonNotFound if the Pokemon doesn't exist.
// Follows SRP: only responsible for fetch + cache.
func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	url := BaseURL + PokemonEndpoint + pokemonName
	
	if val, ok := c.cache.Get(url); ok {
		pokemonResp := Pokemon{}
		err := json.Unmarshal(val, &pokemonResp)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemonResp, nil
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}
	
	pokemonResp := Pokemon{}
	err = json.Unmarshal(data, &pokemonResp)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, data)
	
	return pokemonResp, nil
}