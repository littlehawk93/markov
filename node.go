package markov

import (
	"math/rand"
	"strings"
)

type node struct {
	weightSum    uint64
	children     map[string]*node
	childWeights map[string]uint64
}

func (me *node) addChildren(words []string, index int, ignoreCase bool) {

	weights := make([]uint64, len(words))

	for i := 0; i < len(words); i++ {
		weights[i] = 1
	}

	me.addWeightedChildren(words, weights, index, ignoreCase)
}

func (me *node) addWeightedChildren(words []string, weights []uint64, index int, ignoreCase bool) {

	if index < 0 || index >= len(words) {
		return
	}

	var word string

	if ignoreCase {
		word = strings.ToLower(words[index])
	} else {
		word = words[index]
	}

	me.addWeightedChild(word, weights[index], ignoreCase)

	node, _ := me.children[word]

	node.addWeightedChildren(words, weights, index+1, ignoreCase)
}

func (me *node) addChild(word string, ignoreCase bool) {

	me.addWeightedChild(word, 1, ignoreCase)
}

func (me *node) addWeightedChild(word string, weight uint64, ignoreCase bool) {

	if ignoreCase {
		word = strings.ToLower(word)
	}

	me.weightSum += weight

	if _, ok := me.children[word]; !ok {
		me.children[word] = new(node)
	}

	wordWeight, ok := me.childWeights[word]

	if ok {
		me.childWeights[word] = weight + wordWeight
	} else {
		me.childWeights[word] = weight
	}
}

func (me *node) seek(words []string, index int, ignoreCase bool) *node {

	if index < 0 || index >= len(words) {
		return me
	}

	word := words[index]

	if ignoreCase {
		word = strings.ToLower(word)
	}

	for k, v := range me.children {

		if word == k {
			return v.seek(words, index+1, ignoreCase)
		}
	}

	return nil
}

func (me *node) nextWord(ran *rand.Rand) string {

	decision := uint64(ran.Int63n(int64(me.weightSum)))

	var currWeight uint64
	currWeight = 0

	for key := range me.children {

		childWeight, _ := me.childWeights[key]

		currWeight += childWeight

		if decision < currWeight {
			return key
		}
	}

	return ""
}
