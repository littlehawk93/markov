package markov

import (
	"math/rand"
)

// Chain a markov chain. Generates randomized sequences of text based on training data
type Chain struct {
	ignoreCase            bool
	maxDepth              int
	ran                   *rand.Rand
	wordTreeRoot          *node
	sentenceStartTreeRoot *node
}

// NextSentence generates a new random sentence based on trained data
func (me *Chain) NextSentence() []string {

	words := make([]string, 0)

	words = append(words, me.sentenceStartTreeRoot.nextWord(me.ran))

	for true {

		var node *node

		if len(words) > me.maxDepth {
			node = me.wordTreeRoot.seek(words[len(words)-me.maxDepth:], 0, me.ignoreCase)
		} else {
			node = me.wordTreeRoot.seek(words, 0, me.ignoreCase)
		}

		if node == nil {
			break
		}

		word := node.nextWord(me.ran)

		if word == "" {
			break
		}

		words = append(words, word)
	}

	return words
}

// TrainLine trains the markov chain with a single line of text
func (me *Chain) TrainLine(line []string) {

	for i := 0; i < len(line); i++ {

		word := line[i]

		if i == 0 {
			me.sentenceStartTreeRoot.addChild(word, me.ignoreCase)
			me.wordTreeRoot.addChild(word, me.ignoreCase)
		} else if i < me.maxDepth {
			me.wordTreeRoot.addChildren(line[0:i+1], i, me.ignoreCase)
		} else {
			me.wordTreeRoot.addChildren(line[i-me.maxDepth+1:i+1], i, me.ignoreCase)
		}
	}
}

// Train trains this markov chain by reading arrays of words
func (me *Chain) Train(lines [][]string) {

	for _, line := range lines {

		me.TrainLine(line)
	}
}
