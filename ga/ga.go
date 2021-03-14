package ga

import (
	"math/rand"
	"time"
)

const (
	POPULATION = 1
	NUMBER_GENES = 100
)

const (
	NONE = iota
	JUMP
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type CpuPlayer struct {
	gene []int
	evaluation int
	death bool
	idx int
}

func (player *CpuPlayer) ShouldJump() bool {
	// 落ちたら死ぬので、とにかく飛んだら評価する
	player.evaluation += 1
	return JUMP == player.gene[player.idx]
}

func (player *CpuPlayer) NextStep() {

	if !player.death{
		// 長く生き残った個体は評価する
		player.evaluation += 1
		player.idx += 1

		if 100 == player.idx{
			player.idx = 0
		}
	}
}

type GA struct{
	CpuPlayers [] *CpuPlayer
	population int
}

func NewGA() *GA{
	ga := &GA{}
	ga.init()

	return ga
}

func createInitalGenes(number int) []int{
	cnt:= 0
	gene := make([]int, 0)

	for cnt < number{
		gene = append(gene, rand.Intn(2))
		cnt++
	}

	return gene
}

func (g *GA) init() {
	cnt := 0
	cpuPlayers := [] *CpuPlayer{}

	for cnt < POPULATION {
		player := &CpuPlayer{ gene: createInitalGenes(NUMBER_GENES), evaluation: 0, death: false, idx: 0 }
		cpuPlayers = append(cpuPlayers, player)

		cnt++
	}

	g.CpuPlayers = cpuPlayers
	g.population = POPULATION
}
