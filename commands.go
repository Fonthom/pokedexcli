package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Fonthom/pokedexcli/internal/battle"
	"github.com/Fonthom/pokedexcli/internal/party"
	"github.com/Fonthom/pokedexcli/internal/pokeapi"
	"github.com/Fonthom/pokedexcli/internal/pokecache"
)

type config struct {
	Next     *string
	Previous *string
	cache    *pokecache.Cache
	client   *pokeapi.Client
	pokedex  map[string]pokeapi.Pokemon
	party    *party.Party
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help":    {name: "help", description: "Displays a help message", callback: commandHelp},
		"exit":    {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"map":     {name: "map", description: "Display the next 20 location areas", callback: commandMap},
		"mapb":    {name: "mapb", description: "Display the previous 20 location areas", callback: commandMapb},
		"explore": {name: "explore", description: "List all Pokemon in a location area", callback: commandExplore},
		"catch":   {name: "catch", description: "Try to catch a Pokemon", callback: commandCatch},
		"inspect": {name: "inspect", description: "Inspect a caught Pokemon", callback: commandInspect},
		"pokedex": {name: "pokedex", description: "List all caught Pokemon", callback: commandPokedex},
		"party":   {name: "party", description: "Show your party", callback: commandParty},
		"battle":  {name: "battle", description: "Battle two party members: battle <a> <b>", callback: commandBattle},
		"evolve":  {name: "evolve", description: "Evolve a caught Pokemon if ready", callback: commandEvolve},
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
	resp, err := cfg.client.GetLocationAreas(url)
	if err != nil {
		return err
	}
	cfg.Next = resp.Next
	cfg.Previous = resp.Previous
	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	resp, err := cfg.client.GetLocationAreas(*cfg.Previous)
	if err != nil {
		return err
	}
	cfg.Next = resp.Next
	cfg.Previous = resp.Previous
	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <location-area>")
	}
	resp, err := cfg.client.ExploreArea(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	for _, e := range resp.PokemonEncounters {
		fmt.Printf(" - %s\n", e.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon>")
	}
	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := cfg.client.GetPokemon(name)
	if err != nil {
		return err
	}

	threshold := pokemon.BaseExperience / 2
	if rand.Intn(pokemon.BaseExperience) > threshold {
		fmt.Printf("%s escaped!\n", name)
		return nil
	}

	cfg.pokedex[name] = pokemon
	fmt.Printf("%s was caught!\n", name)
	fmt.Println("You may now inspect it with the inspect command.")

	if err := cfg.party.Add(pokemon); err != nil {
		fmt.Printf("(Party full — %s added to Pokedex only)\n", name)
	} else {
		fmt.Printf("%s was added to your party!\n", name)
	}
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon>")
	}
	pokemon, ok := cfg.pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", pokemon.Name, pokemon.Height, pokemon.Weight)
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

func commandParty(cfg *config, args ...string) error {
	members := cfg.party.Members()
	if len(members) == 0 {
		fmt.Println("Your party is empty")
		return nil
	}
	fmt.Println("Your Party:")
	for _, m := range members {
		fmt.Printf(" - %s (Lv.%d) HP:%d ATK:%d XP:%d\n",
			m.Nickname, m.Level, m.HP(), m.Attack(), m.XP)
	}
	return nil
}

func commandBattle(cfg *config, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: battle <pokemon-a> <pokemon-b>")
	}
	a, okA := cfg.party.Get(args[0])
	b, okB := cfg.party.Get(args[1])
	if !okA {
		return fmt.Errorf("%s not found in party", args[0])
	}
	if !okB {
		return fmt.Errorf("%s not found in party", args[1])
	}
	battle.Simulate(a, b)
	return nil
}

func commandEvolve(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: evolve <pokemon>")
	}
	name := args[0]
	member, inParty := cfg.party.Get(name)
	if !inParty {
		fmt.Println("that pokemon is not in your party")
		return nil
	}

	// require 2 minutes since caught before evolving
	if time.Since(member.CaughtAt) < 2*time.Minute {
		remaining := 2*time.Minute - time.Since(member.CaughtAt)
		fmt.Printf("%s is not ready to evolve yet. Try again in %s\n", name, remaining.Round(time.Second))
		return nil
	}

	species, err := cfg.client.GetSpecies(member.Pokemon.Species.Name)
	if err != nil {
		return err
	}
	chain, err := cfg.client.GetEvolutionChain(species.EvolutionChain.URL)
	if err != nil {
		return err
	}

	nextEvolution := findNextEvolution(chain.Chain, member.Pokemon.Name)
	if nextEvolution == "" {
		fmt.Printf("%s cannot evolve any further\n", name)
		return nil
	}

	evolved, err := cfg.client.GetPokemon(nextEvolution)
	if err != nil {
		return err
	}

	fmt.Printf("What?! %s is evolving!\n", member.Nickname)
	fmt.Printf("%s evolved into %s!\n", member.Nickname, evolved.Name)

	cfg.pokedex[evolved.Name] = evolved
	member.Pokemon = evolved
	member.Nickname = evolved.Name
	member.CaughtAt = time.Now()
	return nil
}

func findNextEvolution(link pokeapi.EvolutionLink, currentName string) string {
	if link.Species.Name == currentName {
		if len(link.EvolvesTo) > 0 {
			return link.EvolvesTo[0].Species.Name
		}
		return ""
	}
	for _, next := range link.EvolvesTo {
		if result := findNextEvolution(next, currentName); result != "" {
			return result
		}
	}
	return ""
}