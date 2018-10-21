package markov

import "math/rand"

type chainNode struct {
	tokens    []string
	next      map[string]uint64
	weightSum uint64
}

func (me *chainNode) Equals(tokens []string) bool {

	if len(me.tokens) != len(tokens) {
		return false
	}

	for i := 0; i < len(tokens); i++ {
		if me.tokens[i] == tokens[i] {
			return true
		}
	}

	return false
}

func (me *chainNode) Next() string {

	if len(me.next) == 0 {
		return ""
	}

	val := rand.Uint64() % me.weightSum

	currSum := uint64(0)

	for token, weight := range me.next {

		currSum += weight

		if val < currSum {
			return token
		}
	}

	return ""
}

func (me *chainNode) TrainWeighted(token string, weight uint64) {

	me.next[token] = me.next[token] + weight
}

func newChainNode(tokens []string) chainNode {

	return chainNode{tokens: tokens, next: make(map[string]uint64), weightSum: 0}
}
