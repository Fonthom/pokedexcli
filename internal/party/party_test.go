package party

import (
	"testing"

	"github.com/Fonthom/pokedexcli/internal/pokeapi"
)

func TestAddAndGet(t *testing.T) {
	p := NewParty()
	pokemon := pokeapi.Pokemon{Name: "pikachu", BaseExperience: 112}
	if err := p.Add(pokemon); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := p.Get("pikachu")
	if !ok {
		t.Fatal("expected to find pikachu in party")
	}
	if m.Level != 1 {
		t.Errorf("expected level 1, got %d", m.Level)
	}
}

func TestPartyFull(t *testing.T) {
	p := NewParty()
	for i := 0; i < maxPartySize; i++ {
		p.Add(pokeapi.Pokemon{Name: "mon"})
	}
	err := p.Add(pokeapi.Pokemon{Name: "one-too-many"})
	if err == nil {
		t.Error("expected error when adding to full party")
	}
}

func TestLevelUp(t *testing.T) {
	p := NewParty()
	p.Add(pokeapi.Pokemon{Name: "charmander"})
	m, _ := p.Get("charmander")
	m.GainXP(100)
	if m.Level != 2 {
		t.Errorf("expected level 2 after 100 XP, got %d", m.Level)
	}
}