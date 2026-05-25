package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Fonthom/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

type Client struct {
	cache *pokecache.Cache
}

func NewClient(cache *pokecache.Cache) *Client {
	return &Client{cache: cache}
}

func (c *Client) fetch(url string) ([]byte, error) {
	if cached, ok := c.cache.Get(url); ok {
		return cached, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching %s: %w", url, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	c.cache.Add(url, body)
	return body, nil
}

func (c *Client) GetLocationAreas(pageURL string) (LocationAreaResponse, error) {
	body, err := c.fetch(pageURL)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	var resp LocationAreaResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error unmarshalling: %w", err)
	}
	return resp, nil
}

func (c *Client) ExploreArea(area string) (ExploreResponse, error) {
	url := baseURL + "/location-area/" + area
	body, err := c.fetch(url)
	if err != nil {
		return ExploreResponse{}, err
	}
	var resp ExploreResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return ExploreResponse{}, fmt.Errorf("error unmarshalling: %w", err)
	}
	return resp, nil
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name
	body, err := c.fetch(url)
	if err != nil {
		return Pokemon{}, err
	}
	var p Pokemon
	if err := json.Unmarshal(body, &p); err != nil {
		return Pokemon{}, fmt.Errorf("error unmarshalling: %w", err)
	}
	return p, nil
}

func (c *Client) GetSpecies(name string) (SpeciesResponse, error) {
	url := baseURL + "/pokemon-species/" + name
	body, err := c.fetch(url)
	if err != nil {
		return SpeciesResponse{}, err
	}
	var s SpeciesResponse
	if err := json.Unmarshal(body, &s); err != nil {
		return SpeciesResponse{}, fmt.Errorf("error unmarshalling: %w", err)
	}
	return s, nil
}

func (c *Client) GetEvolutionChain(url string) (EvolutionChain, error) {
	body, err := c.fetch(url)
	if err != nil {
		return EvolutionChain{}, err
	}
	var ec EvolutionChain
	if err := json.Unmarshal(body, &ec); err != nil {
		return EvolutionChain{}, fmt.Errorf("error unmarshalling: %w", err)
	}
	return ec, nil
}