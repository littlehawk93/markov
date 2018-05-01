package markov

import (
	"math/rand"
	"time"
)

// NewChain creates a new empty, untrained markov chain of the specified depth.
// Depth specifies how deep the markov chain's underlying state tree can grow while training.
// If ignoreCase is true, the markov chain will be case insensative when parsing training data.
func NewChain(depth int, ignoreCase bool) *Chain {

	var chain Chain

	if depth < 1 {
		chain.maxDepth = 1
	} else {
		chain.maxDepth = depth
	}

	chain.ignoreCase = ignoreCase

	chain.wordTreeRoot = new(node)
	chain.sentenceStartTreeRoot = new(node)

	chain.ran = new(rand.Rand)
	chain.ran.Seed(time.Now().UTC().UnixNano())

	return &chain
}
