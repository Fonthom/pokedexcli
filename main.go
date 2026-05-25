package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Fonthom/pokedexcli/internal/party"
	"github.com/Fonthom/pokedexcli/internal/pokeapi"
	"github.com/Fonthom/pokedexcli/internal/pokecache"
)

func main() {
	cache := pokecache.NewCache(5 * time.Second)
	cfg := &config{
		cache:   cache,
		client:  pokeapi.NewClient(cache),
		pokedex: make(map[string]pokeapi.Pokemon),
		party:   party.NewParty(),
	}
	commands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		if cmd, ok := commands[words[0]]; ok {
			if err := cmd.callback(cfg, words[1:]...); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}