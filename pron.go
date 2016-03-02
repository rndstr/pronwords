package pronwords

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"strings"
)

type Pronounceable struct {
	unigram   map[string]int // Mapping of unigrams to the number of occurrences.
	bigram    map[string]int // Mapping of bigrams to the number of occurrences.
	trigram   map[string]int // Mapping of trigrams to the number of occurrences.

	uninorm   float64        // Normalization for the unigram score.
	binorm    float64        //  Normalization for the bigram score.
	trinorm   float64        // Normalization for the trigram score.
	normdirty bool           // Flag whether the norms are dirty and need updating.

	uniweight float64        // Weighting of matched unigram occurrences.
	biweight  float64        // Weighting of matched bigram occurrences.
	triweight float64        // Weighting of matched trigram occurrences.
}

const (
	UniWeightDefault = 1.0
	BiWeightDefault = 3.0
	TriWeightDefault = 5.0
)

func NewPronounceable() *Pronounceable {
	return &Pronounceable{
		unigram: make(map[string]int),
		bigram: make(map[string]int),
		trigram: make(map[string]int),

		uniweight: UniWeightDefault,
		biweight: BiWeightDefault,
		triweight: TriWeightDefault,
	}
}

// SetWeights updates the weight distribution of the n-grams.
func (p *Pronounceable) SetWeights(uni, bi, tri float64) {
	p.uniweight = uni
	p.biweight = bi
	p.triweight = tri
}

// AddWordList takes a reader and scans it word by word, recording their n-grams.
func (p *Pronounceable) AddWordList(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		p.AddWord(word)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading word list:", err)
	}
}

// AddWord records n-grams of the given word.
func (p *Pronounceable) AddWord(word string) {
	word = strings.ToLower(word)

	for i := 0; i < len(word); i += 1 {
		if i < len(word) - 2 {
			p.trigram[word[i:i + 3]] += 1
		}
		if i < len(word) - 1 {
			p.bigram[word[i:i + 2]] += 1
		}
		p.unigram[word[i:i + 1]] += 1
	}

	p.normdirty = true
}

// WordScore calculates a score that attempts to express how easy it is to
// pronounce the word.
func (p *Pronounceable) WordScore(word string) float64 {
	word = strings.ToLower(word)
	if p.normdirty {
		p.calculateNormalization()
	}

	score := 0.0
	for i := 0; i < len(word); i += 1 {
		// If a character does not exist in unigram, it is unclassifiable and
		// any statement about pronounceability
		// cannot say anything about its pronounceability anymore. Just bail.
		if p.unigram[word[i:i + 1]] == 0 {
			return 0.0
		}

		if i < len(word) - 2 {
			score += p.triweight * float64(p.trigram[word[i:i + 3]]) / p.trinorm
		}
		if i < len(word) - 1 {
			score += p.biweight * float64(p.bigram[word[i:i + 2]]) / p.binorm
		}
		score += p.uniweight * float64(p.unigram[word[i:i + 1]]) / p.uninorm
	}

	// Normalize by how many scores have been computed
	lengthnorm := 1.0
	if len(word) > 1 {
		lengthnorm = float64(len(word) - 1) * 3.0
	}

	return score / lengthnorm
}

// IsPronounceable determines whether a word is pronounceable by comparing the word
// score to the given threshold.
func (p *Pronounceable) IsPronounceable(word string, threshold float64) bool {
	return p.WordScore(word) >= threshold
}

// CalculateNormalization updates the uninorm, binorm, and trinorm aggregators.
// They influences score calculation by dividing the n-gram count.
func (p *Pronounceable) calculateNormalization() {
	if !p.normdirty {
		return
	}

	p.uninorm = 0
	p.binorm = 0
	p.trinorm = 0

	p.uninorm = maxInGram(p.unigram)
	p.binorm = maxInGram(p.bigram)
	p.trinorm = maxInGram(p.trigram)

	p.normdirty = false
}

func maxInGram(gram map[string]int) float64 {
	max := 0
	for _, o := range gram {
		if o > max {
			max = o
		}
	}

	return float64(max)
}

