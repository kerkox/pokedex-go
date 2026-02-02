package pokeapi

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height		 int    `json:"height"`
	Weight		 int    `json:"weight"`
	Stats		 []PokemonStat `json:"stats"`
	Types		 []PokemonType `json:"types"`
		
}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Effort int `json:"effort"`
	Stat StatInfo `json:"stat"`
}

type StatInfo struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type PokemonType struct {
	Slot int `json:"slot"`
	TypeInfo TypeInfo `json:"type"`
}

type TypeInfo struct {
	Name string `json:"name"`
	URL string `json:"url"`
}