package party

import (
	"fmt"
	"time"

	"github.com/Fonthom/pokedexcli/internal/pokeapi"
)

const maxPartySize = 6

type PartyMember struct {
	Pokemon    pokeapi.Pokemon
	Level      int
	XP         int
	CaughtAt   time.Time
	Nickname   string
}

type Party struct {
	members []*PartyMember
}

func NewParty() *Party {
	return &Party{}
}

func (p *Party) Add(pokemon pokeapi.Pokemon) error {
	if len(p.members) >= maxPartySize {
		return fmt.Errorf("party is full (max %d)", maxPartySize)
	}
	p.members = append(p.members, &PartyMember{
		Pokemon:  pokemon,
		Level:    1,
		XP:       0,
		CaughtAt: time.Now(),
		Nickname: pokemon.Name,
	})
	return nil
}

func (p *Party) Members() []*PartyMember {
	return p.members
}

func (p *Party) Get(name string) (*PartyMember, bool) {
	for _, m := range p.members {
		if m.Nickname == name || m.Pokemon.Name == name {
			return m, true
		}
	}
	return nil, false
}

func (m *PartyMember) GainXP(amount int) {
	m.XP += amount
	threshold := m.Level * 100
	for m.XP >= threshold {
		m.XP -= threshold
		m.Level++
		fmt.Printf("  %s leveled up to level %d!\n", m.Nickname, m.Level)
		threshold = m.Level * 100
	}
}

func (m *PartyMember) HP() int {
	base := 0
	for _, s := range m.Pokemon.Stats {
		if s.Stat.Name == "hp" {
			base = s.BaseStat
		}
	}
	return base + m.Level*5
}

func (m *PartyMember) Attack() int {
	base := 0
	for _, s := range m.Pokemon.Stats {
		if s.Stat.Name == "attack" {
			base = s.BaseStat
		}
	}
	return base + m.Level*2
}