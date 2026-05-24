package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Fonthom/pokedexcli/internal/pokecache"
)

func main() {
	cfg := &config{
		cache: pokecache.NewCache(5 * time.Second),
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
		commandName := words[0]
		if cmd, ok := commands[commandName]; ok {
    		if err := cmd.callback(cfg, words[1:]...); err != nil {
        		fmt.Println("Error:", err)
    		}
		} else {
			fmt.Println("Unknown command")
		}
	}
}