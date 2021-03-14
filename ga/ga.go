package ga

const POPULATION = 1
const NUMBER_GENES = 100

type CpuPlayer struct {
	gene []int
	stepCnt int
	death bool
	idx int
}

func (player *CpuPlayer) ShouldJump() bool {
	return true
}

func (player *CpuPlayer) NextStep() {
	if !player.death{
		player.stepCnt++
		player.idx++

		if 100 == player.idx{
			player.idx = 0
		}
	}
}

type GA struct{
	CpuPlayers [] CpuPlayer
	population int
}

func NewGA() *GA{
	ga := &GA{}
	ga.init()

	return ga
}

func (g *GA) init() {
	cnt := 0
	cpuPlayers := [] CpuPlayer{}

	gene := [] int{1}
	for cnt < POPULATION {
		player := CpuPlayer{ gene: gene, stepCnt: 0, death: false, idx: 0 }
		cpuPlayers = append(cpuPlayers, player)

		cnt++
	}

	g.CpuPlayers = cpuPlayers
	g.population = POPULATION
}
