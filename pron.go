package pronwords

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

type Pronouncable struct {
    unigram map[string]int
    bigram map[string]int
    trigram map[string]int

    unisum int
    bisum int
    trisum int

    sumDirty bool

    uniweight float64
    biweight float64
    triweight float64
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

func (p *Pronouncable) AddWordList(path string) {
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)

    for scanner.Scan() {
        word := strings.ToUpper(scanner.Text())
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

    if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading file %s:", path, err)
	}

    p.sumDirty = true
}

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

