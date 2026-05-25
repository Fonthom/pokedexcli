package battle

import (
	"testing"

	"github.com/Fonthom/pokedexcli/internal/party"
	"github.com/Fonthom/pokedexcli/internal/pokeapi"
)

func makeMember(name string, hp, attack, baseExp, level int) *party.PartyMember {
	return &party.PartyMember{
		Nickname: name,
		Level:    level,
		Pokemon: pokeapi.Pokemon{
			Name:           name,
			BaseExperience: baseExp,
			Stats: []struct {
				BaseStat int `json:"base_stat"`
				Stat     struct {
					Name string `json:"name"`
				} `json:"stat"`
			}{
				{BaseStat: hp, Stat: struct {
					Name string `json:"name"`
				}{Name: "hp"}},
				{BaseStat: attack, Stat: struct {
					Name string `json:"name"`
				}{Name: "attack"}},
			},
		},
	}
}

func TestSimulateReturnsWinner(t *testing.T) {
	a := makeMember("bulbasaur", 100, 50, 64, 10)
	b := makeMember("caterpie", 10, 5, 39, 1)
	result := Simulate(a, b)
	if result.Winner == nil {
		t.Error("expected a winner, got nil")
	}
	if result.Rounds == 0 {
		t.Error("expected at least one round")
	}
}

func TestSimulateXPAwarded(t *testing.T) {
	a := makeMember("strong", 200, 100, 200, 20)
	b := makeMember("weak", 10, 1, 50, 1)
	before := a.XP
	Simulate(a, b)
	if a.XP == before && a.Level == 20 {
		t.Error("expected winner to gain XP")
	}
}