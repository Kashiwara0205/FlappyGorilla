package ga

import (
	"math/rand"
	"time"
	"sort"
	"fmt"
)

const (
	POPULATION_SIZE = 1300
	GENE_SIZE = 1200
	ACTION_SPAN = 1
	MUTATION_RATE = 50
)

const (
	NONE = iota
	JUMP
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Individual struct {
	gene []int
	Score int
	death bool
	locus int
}

func (individual *Individual) ShouldJump() bool {
	return JUMP == individual.gene[individual.locus]
}

func (individual *Individual) NextStep() {
	if !individual.death{
		individual.locus += 1

		if GENE_SIZE == individual.locus{
			individual.locus = 0
		}
	}
}

func (individual *Individual) Dead() {
	individual.death = true
}

func (individual *Individual) CheckDead() bool{
	return individual.death

}

type GA struct{
	Individuals [] *Individual
}

func (ga *GA) CheckAllDead() bool {
	for _, individual := range ga.Individuals{
		if !individual.death{
			return false
		}
	}

	return true
}

func mergeGene(values []int, values2 []int) []int{
	result := []int{}
	for _, r := range values{
		result = append(result, r)
	}

	for _, r := range values2{
		result = append(result, r)
	}

	return result
}

func mutation(gene []int) []int{
	if 1 == rand.Intn(MUTATION_RATE){
		cnt := 0
		for cnt < 200 {
			locus := rand.Intn(GENE_SIZE)
			if 1== gene[locus]{
				gene[locus] = 0
			}else{
				gene[locus] = 1
			}
			cnt ++;
		}
	}

	return gene
}

func copyIndividuals(Individuals [] *Individual, rangeA, rangeB int) [] *Individual{
	copyIndividuals := [] *Individual{}
	for _, individual := range Individuals[rangeA:rangeB]{
		copyIndividuals = append(copyIndividuals, &Individual{ gene: individual.gene , Score: 0, death: false, locus: 0 })
	}

	return copyIndividuals
}

func createParentsA(Individuals [] *Individual)[] *Individual{
	parentsA  := [] *Individual{}
	for i:= 0; i < POPULATION_SIZE / 3; i++{
		locus := rand.Intn(GENE_SIZE / 2)
		parentsA = append(parentsA, Individuals[locus])
	}

	return parentsA
}

func appendNewIndividual(newIndividuals [] *Individual, individualA *Individual, individualB *Individual) [] *Individual{
	separationPoint := rand.Intn(GENE_SIZE)
	childGenA := individualA.gene[0:separationPoint]
	childGenB := individualB.gene[separationPoint:GENE_SIZE]

	gene := mergeGene(childGenA, childGenB)
	gene = mutation(gene)
	child := &Individual{ gene: gene, Score: 0, death: false, locus: 0 }

	newIndividuals = append(newIndividuals, child)

	return newIndividuals
}

func (ga *GA) Update() {
	Individuals := ga.Individuals
	sort.Slice(Individuals, func(i, j int) bool { return Individuals[i].Score > Individuals[j].Score })

	fmt.Printf("[ Ranking ]\n")
	fmt.Printf("------------------------------------\n")
	for i, individual := range Individuals[0:30]{
		fmt.Printf("Rank: %v, Score: %v\n",  i + 1, individual.Score)
	}
	fmt.Printf("------------------------------------\n")

	topIndividuals := copyIndividuals(Individuals, 0, 100)
	middleIndividuals := copyIndividuals(Individuals, 100, 900)
	lowIndividuals := copyIndividuals(Individuals, 900, 1000)

	newIndividuals := [] *Individual{}
	newIndividuals = append(newIndividuals, copyIndividuals(Individuals, 0, 300)...)

	for i := 0; i < 10; i++ {
		individualA := topIndividuals[rand.Intn(100)]
		individualB := middleIndividuals[rand.Intn(800)]

		newIndividuals = appendNewIndividual(newIndividuals, individualA, individualB)
	}
	for j := 0; j < 980; j++ {
		individualA := middleIndividuals[rand.Intn(800)]
		individualB := middleIndividuals[rand.Intn(800)]

		newIndividuals = appendNewIndividual(newIndividuals, individualA, individualB)
	}
	for j := 0; j < 10; j++ {
		individualA := lowIndividuals[rand.Intn(100)]
		individualB := middleIndividuals[rand.Intn(800)]

		newIndividuals = appendNewIndividual(newIndividuals, individualA, individualB)
	}

	fmt.Printf("=> GORIRA SIZE%v\n", len(newIndividuals))
	ga.Individuals = newIndividuals[0:POPULATION_SIZE]
}

func NewGA() *GA{
	ga := &GA{}
	ga.init()

	return ga
}

func createInitalGenes(number int, randNumber int) []int{
	cnt:= 0
	gene := make([]int, 0)

	for cnt < number{
		action := 0
		if 0 == rand.Intn(randNumber){
			action = 1
		}

		gene = append(gene, action)
		cnt++
	}

	return gene
}

func (g *GA) init() {
	Individuals := [] *Individual{}

	for i := 0; i < POPULATION_SIZE; i++ {
		individual := &Individual{ gene: createInitalGenes(GENE_SIZE, rand.Intn(20) + 2), Score: 0, death: false, locus: 0 }
		Individuals = append(Individuals, individual)
	}

	g.Individuals = Individuals
}
