package pronwords

import (
        "fmt"
        "os"
        "bufio"
        "strings"
)

type Pronouncable struct {
    unigram map[string]int // Mapping of unigrams to the number of occurrences.
    bigram map[string]int // Mapping of bigrams to the number of occurrences.
    trigram map[string]int // Mapping of trigrams to the number of occurrences.

    unisum int // Sum of all occurences of unigrams.
    bisum int // Sum of all occurences of bigrams.
    trisum int // Sum of all occurences of trigrams.

    sumDirty bool // Flag whether the sums are dirty and need updating.

    uniweight float64 // Weighting of matched unigram occurences.
    biweight float64 // Weighting of matched bigram occurences.
    triweight float64 // Weighting of matched trigram occurences.
}

const (
    UniWeightDefault = 1.0
    BiWeightDefault = 3.0
    TriWeightDefault = 5.0
)

func NewPronouncable() *Pronouncable {
    return &Pronouncable{
        unigram: make(map[string]int),
        bigram: make(map[string]int),
        trigram: make(map[string]int),

        uniweight: UniWeightDefault,
        biweight: BiWeightDefault,
        triweight: TriWeightDefault,
    }
}

// SetWeights updates the weight distribution of the n-grams.
func (p *Pronouncable) SetWeights(uni, bi, tri float64) {
    p.uniweight = uni
    p.biweight = bi
    p.triweight = tri
}

// AddWordListFile takes a file path and scans the words, recording their n-grams.
func (p *Pronouncable) AddWordListFile(path string) {
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)

    for scanner.Scan() {
        word := scanner.Text()
        p.AddWord(word)
    }

    if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading file %s:", path, err)
	}

    p.sumDirty = true
}

// AddWord records n-grams of the given word.
func (p *Pronouncable) AddWord(word string) {
    word = strings.ToUpper(word)

    for i := 0; i < len(word); i += 1 {
        if i < len(word) - 2 {
            p.trigram[word[i:i+3]] += 1
        }
        if i < len(word) - 1 {
            p.bigram[word[i:i+2]] += 1
        }
        p.unigram[word[i:i+1]] += 1
    }
}

// WordScore calculates a score that attempts to express how easy it is to
// pronounce the word.
func (p *Pronouncable) WordScore(word string) float64 {
    word = strings.ToUpper(word)
    if p.sumDirty {
        p.calculateSum()
    }

    score := 0.0

    for i := 0; i < len(word); i += 1 {
        if i <= len(word) - 3 {
            score += 5.0 * float64(p.trigram[word[i:i+3]]) / float64(p.trisum)
        }
        if i <= len(word) - 2 {
            score += float64(p.bigram[word[i:i+2]]) / float64(p.bisum)
        }
        score += float64(p.unigram[word[i:i+1]]) / float64(p.unisum)
    }

    return score
}

// CalculateSum updates the unisum, bisum, and trisum aggregators.
func (p *Pronouncable) calculateSum() {
    if !p.sumDirty {
        return
    }

    p.unisum = 0
    p.bisum = 0
    p.trisum = 0

    for _,o := range p.unigram { p.unisum += o }
    for _,o := range p.bigram { p.bisum += o }
    for _,o := range p.trigram { p.trisum += o }

    p.sumDirty = false
}

