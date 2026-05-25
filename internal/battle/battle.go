package battle

import (
	"fmt"
	"math/rand"

	"github.com/Fonthom/pokedexcli/internal/party"
)

type BattleResult struct {
	Winner *party.PartyMember
	Loser  *party.PartyMember
	Rounds int
}

func Simulate(a, b *party.PartyMember) BattleResult {
	hpA := a.HP()
	hpB := b.HP()
	rounds := 0

	fmt.Printf("Battle: %s (Lv.%d) vs %s (Lv.%d)\n", a.Nickname, a.Level, b.Nickname, b.Level)

	for hpA > 0 && hpB > 0 {
		rounds++
		dmgToB := rand.Intn(a.Attack()) + 1
		dmgToA := rand.Intn(b.Attack()) + 1
		hpB -= dmgToB
		hpA -= dmgToA
		fmt.Printf("  Round %d: %s deals %d dmg | %s deals %d dmg\n",
			rounds, a.Nickname, dmgToB, b.Nickname, dmgToA)
	}

	if hpA > hpB {
		fmt.Printf("  %s wins after %d rounds!\n", a.Nickname, rounds)
		a.GainXP(b.Pokemon.BaseExperience)
		return BattleResult{Winner: a, Loser: b, Rounds: rounds}
	}
	fmt.Printf("  %s wins after %d rounds!\n", b.Nickname, rounds)
	b.GainXP(a.Pokemon.BaseExperience)
	return BattleResult{Winner: b, Loser: a, Rounds: rounds}
}