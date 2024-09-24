package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PokeMapConfig struct {
	Prev *string
	Next *string
}

type pokeMap struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func PkmnMap(cfg *PokeMapConfig) error {
	url := "https://pokeapi.co/api/v2/location-area/"

	if cfg.Next != nil {
		url = *cfg.Next
	}

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http error: %v", err)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	pokeRes := pokeMap{}

	if err = decoder.Decode(&pokeRes); err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	for _, loc := range pokeRes.Results {
		fmt.Println(loc.Name)
	}

	cfg.Prev = pokeRes.Previous
	cfg.Next = pokeRes.Next

	return nil
}

func PkmnMapb(cfg *PokeMapConfig) error {
	if cfg.Prev == nil {
		return fmt.Errorf("error: no previous locations available")
	}

	res, err := http.Get(*cfg.Prev)
	if err != nil {
		return fmt.Errorf("http error: %v", err)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	pokeRes := pokeMap{}

	if err = decoder.Decode(&pokeRes); err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	for _, loc := range pokeRes.Results {
		fmt.Println(loc.Name)
	}

	cfg.Prev = pokeRes.Previous
	cfg.Next = pokeRes.Next

	return nil
}
