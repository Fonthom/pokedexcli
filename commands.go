package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"math/rand"

	"github.com/Fonthom/pokedexcli/internal/pokecache"
)

type config struct {
	Next     *string
	Previous *string
	cache    *pokecache.Cache
	pokedex  map[string]Pokemon
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
	callback    func(*config, ...string) error
}

type exploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
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
		"explore": {
    		name:        "explore",
    		description: "List all Pokemon in a location area",
    		callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
    		description: "Try to catch a Pokemon",
    		callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
    		name:        "pokedex",
    		description: "List all caught Pokemon",
    		callback:    commandPokedex,
		},
	}
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	if cfg.Next != nil {
		url = *cfg.Next
	}
	return fetchAndPrintLocations(cfg, url)
}

func commandMapb(cfg *config, args ...string) error {
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

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <location-area>")
	}
	area := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + area

	var body []byte
	if cached, ok := cfg.cache.Get(url); ok {
		body = cached
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error fetching area: %w", err)
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %w", err)
		}
		cfg.cache.Add(url, body)
	}

	var exploreResp exploreResponse
	if err := json.Unmarshal(body, &exploreResp); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	fmt.Printf("Exploring %s...\n", area)
	fmt.Println("Found Pokemon:")
	for _, e := range exploreResp.PokemonEncounters {
		fmt.Printf(" - %s\n", e.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon>")
	}
	name := args[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	var body []byte
	if cached, ok := cfg.cache.Get(url); ok {
		body = cached
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error fetching pokemon: %w", err)
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %w", err)
		}
		cfg.cache.Add(url, body)
	}

	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	// higher base experience = harder to catch
	// e.g. base_experience 100 -> ~66% catch chance
	//      base_experience 300 -> ~25% catch chance
	threshold := pokemon.BaseExperience / 2
	roll := rand.Intn(pokemon.BaseExperience)
	if roll > threshold {
		fmt.Printf("%s escaped!\n", name)
		return nil
	}

	cfg.pokedex[name] = pokemon
	fmt.Printf("%s was caught!\n", name)
	fmt.Println("You may now inspect it with the inspect command.")
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon>")
	}
	name := args[0]

	pokemon, ok := cfg.pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.pokedex) == 0 {
		fmt.Println("You have not caught any Pokemon yet")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name := range cfg.pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}