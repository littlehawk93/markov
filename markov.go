package markov

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"math/rand"
	"strings"
	"time"
)

const (
	delimTypeNone   = 0
	delimTypeIgnore = 1
	delimTypeLine   = 2
	delimTypeToken  = 3
)

// Chain a markov chain. Generates randomized sequences of text based on training data
type Chain struct {
	ignoreCase            bool
	maxDepth              int
	ran                   *rand.Rand
	wordTreeRoot          *node
	sentenceStartTreeRoot *node
}

type node struct {
	weightSum    uint64
	children     map[string]*node
	childWeights map[string]uint64
}

func (me *node) addChildren(words []string, weights []uint64, index int, ignoreCase bool) {

	if index < 0 || index >= len(words) {
		return
	}

	var word string

	if ignoreCase {
		word = strings.ToLower(words[index])
	} else {
		word = words[index]
	}

	me.addChild(word, weights[index], ignoreCase)

	node, _ := me.children[word]

	node.addChildren(words, weights, index+1, ignoreCase)
}

func (me *node) addChild(word string, weight uint64, ignoreCase bool) {

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

func (me *Chain) trainLine(words []string) {

	weights := make([]uint64, len(words))

	for i := 0; i < len(weights); i++ {
		weights[i] = 1
	}

	for i := 0; i < len(words); i++ {

		word := words[i]

		if i == 0 {
			me.sentenceStartTreeRoot.addChild(word, 1, me.ignoreCase)
			me.wordTreeRoot.addChild(word, 1, me.ignoreCase)
		} else if i < me.maxDepth {
			me.wordTreeRoot.addChildren(words[0:i+1], weights[0:i+1], i, me.ignoreCase)
		} else {
			me.wordTreeRoot.addChildren(words[i-me.maxDepth+1:i+1], weights[i-me.maxDepth+1:i+1], i, me.ignoreCase)
		}
	}
}

// Train trains this markov chain by reading text from the provided reader
func (me *Chain) Train(reader io.Reader, lineDelim []rune, tokenDelim []rune, ignore []rune) error {

	runeLookup, err := constructRuneLookupMap(lineDelim, tokenDelim, ignore)

	if err != nil {
		return err
	}

	bReader := bufio.NewReader(reader)

	lineBuffer := make([]string, 0)
	var wordBuffer bytes.Buffer
	var r rune

	for {
		if r, _, err = bReader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		runeType, ok := runeLookup[r]

		if ok {
			if runeType == delimTypeLine {

				if wordBuffer.Len() > 0 {
					lineBuffer = append(lineBuffer, wordBuffer.String())
				}

				me.trainLine(lineBuffer)

				wordBuffer.Reset()
				lineBuffer = make([]string, 0)

			} else if runeType == delimTypeToken {

				lineBuffer = append(lineBuffer, wordBuffer.String())
				wordBuffer.Reset()

			} else if runeType == delimTypeIgnore {
				continue
			} else {
				wordBuffer.WriteRune(r)
			}
		} else {
			wordBuffer.WriteRune(r)
		}
	}

	return nil
}

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

func constructRuneLookupMap(lineDelim, tokenDelim, ignore []rune) (map[rune]int, error) {

	var runeLookup map[rune]int

	for _, r := range lineDelim {

		if _, ok := runeLookup[r]; ok {
			return nil, errors.New("Duplicate rune")
		}

		runeLookup[r] = delimTypeLine
	}

	for _, r := range tokenDelim {

		if _, ok := runeLookup[r]; ok {
			return nil, errors.New("Duplicate rune")
		}

		runeLookup[r] = delimTypeToken
	}

	for _, r := range ignore {

		if _, ok := runeLookup[r]; ok {
			return nil, errors.New("Duplicate rune")
		}

		runeLookup[r] = delimTypeIgnore
	}

	return runeLookup, nil
}
