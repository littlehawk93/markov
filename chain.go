package markov

import (
	"math/rand"

	"github.com/chobie/go-gaussian"
	"github.com/littlehawk93/rollstat"
)

// Chain a markov chain
type Chain struct {
	keyDepth  uint
	nodes     []chainNode
	wordStats *rollstat.FloatStat
}

// Generate create a new line of text based on trained data for this Markov chain
func (me *Chain) Generate() []string {

	dist := gaussian.NewGaussian(me.wordStats.Mean(), me.wordStats.Var())

	result := make([]string, 0)

	key := me.newKey()

	for rand.Float64() > dist.Cdf(float64(len(result))) {

		index := me.findNode(key)

		if index == -1 {
			break
		}

		nextToken := me.nodes[index].Next()

		if nextToken == "" {
			break
		}

		result = append(result, nextToken)
		key = me.pushKey(key, nextToken)
	}

	return result
}

// Train train this Markov chain with a set of tokens of equal weight
func (me *Chain) Train(tokens []string) {

	me.TrainWeighted(tokens, 1)
}

// TrainWeighted train this Markov chain with a set of tokens at as specified weight
func (me *Chain) TrainWeighted(tokens []string, weight uint64) {

	me.wordStats.Add(float64(len(tokens)))

	key := me.newKey()

	for _, token := range tokens {

		index := me.matchNode(key)

		if index == -1 {
			newNode := newChainNode(key)
			me.nodes = append(me.nodes, newNode)
			index = len(me.nodes) - 1
		}

		me.nodes[index].TrainWeighted(token, weight)

		key = me.pushKey(key, token)
	}

	index := me.matchNode(key)

	if index == -1 {
		newNode := newChainNode(key)
		me.nodes = append(me.nodes, newNode)
		index = len(me.nodes) - 1
	}

	me.nodes[index].TrainWeighted("", weight)
}

func (me *Chain) findNode(key []string) int {

	for i := 0; i < len(key); i++ {

		if index := me.matchNode(key); index != -1 {
			return index
		}

		key[i] = ""
	}

	return -1
}

func (me *Chain) matchNode(key []string) int {

	for i, node := range me.nodes {
		if node.Equals(key) {
			return i
		}
	}

	return len(me.nodes) - 1
}

func (me *Chain) newKey() []string {

	key := make([]string, me.keyDepth)

	for i := uint(0); i < me.keyDepth; i++ {
		key[i] = ""
	}

	return key
}

func (me *Chain) pushKey(key []string, value string) []string {

	key = append(key, value)
	return key[1:]
}

// NewChain initialize a new Markov Chain instance
func NewChain(depth uint) *Chain {

	return &Chain{keyDepth: depth, nodes: make([]chainNode, 0), wordStats: &rollstat.FloatStat{}}
}
