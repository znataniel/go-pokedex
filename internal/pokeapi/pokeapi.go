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

var cache *pokecache.Cache = pokecache.NewCache(time.Minute * 5)

func GetLocations(url string) (pokeMap, error) {

	if data, exists := cache.Get(url); exists {
		pokeRes := pokeMap{}
		if err := json.Unmarshal(data, &pokeRes); err != nil {
			return pokeMap{}, err
		}
		return pokeRes, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return pokeMap{}, fmt.Errorf("http error: %v", err)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	pokeRes := pokeMap{}

	if err = decoder.Decode(&pokeRes); err != nil {
		return pokeMap{}, fmt.Errorf("error reading response body: %v", err)
	}

	pokeResBytes, err := json.Marshal(pokeRes)
	if err != nil {
		return pokeMap{}, fmt.Errorf("error storing cache: %v", err)
	}

	if err := cache.Add(url, pokeResBytes); err != nil {
		return pokeMap{}, err
	}

	return pokeRes, nil
}
