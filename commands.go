package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Fonthom/pokedexcli/internal/pokecache"
)

type config struct {
	Next     *string
	Previous *string
	cache    *pokecache.Cache
}

type locationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapb,
		},
	}
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	if cfg.Next != nil {
		url = *cfg.Next
	}
	return fetchAndPrintLocations(cfg, url)
}

func commandMapb(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return fetchAndPrintLocations(cfg, *cfg.Previous)
}

func fetchAndPrintLocations(cfg *config, url string) error {
	var body []byte

	if cached, ok := cfg.cache.Get(url); ok {
		fmt.Println("(cache hit)")
		body = cached
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error fetching locations: %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %w", err)
		}
		cfg.cache.Add(url, body)
	}

	var locationResp locationAreaResponse
	if err := json.Unmarshal(body, &locationResp); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	cfg.Next = locationResp.Next
	cfg.Previous = locationResp.Previous

	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	return nil
}