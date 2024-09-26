package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetLocations(url string) (pokeMap, error) {
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

	return pokeRes, nil
}
