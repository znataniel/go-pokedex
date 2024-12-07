package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/znataniel/go-pokedex/internal/pokecache"
)

const BaseUrl string = "https://pokeapi.co/api/v2/"

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

type pokeEncounter struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

var cache *pokecache.Cache = pokecache.NewCache(time.Minute * 5)

func makeRequest[T any](url string, dest *T) error {
	if data, exists := cache.Get(url); exists {
		if err := json.Unmarshal(data, dest); err != nil {
			return err
		}
		return nil
	}

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http error: %v", err)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(dest); err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	pokeResBytes, err := json.Marshal(dest)
	if err != nil {
		return fmt.Errorf("error storing cache: %v", err)
	}

	if err := cache.Add(url, pokeResBytes); err != nil {
		return err
	}

	return nil
}

func GetLocations(url string) (pokeMap, error) {
	pokeRes := pokeMap{}
	err := makeRequest(url, &pokeRes)
	if err != nil {
		return pokeMap{}, err
	}
	return pokeRes, nil
}

func GetLocationData(url string) (pokeEncounter, error) {
	pokeRes := pokeEncounter{}
	err := makeRequest(url, &pokeRes)
	if err != nil {
		return pokeEncounter{}, err
	}
	return pokeRes, nil
}

func GetPokemon(url string) (Pokemon, error) {
	pokeRes := Pokemon{}
	err := makeRequest(url, &pokeRes)
	if err != nil {
		return Pokemon{}, err
	}
	return pokeRes, nil
}
