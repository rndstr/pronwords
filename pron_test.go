package pronwords

import (
	"testing"
	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestAddWord(t *testing.T) {
	p := NewPronouncable()
	p.AddWord("fool")

	assert.Equal(t, 1, p.unigram["F"])
	assert.Equal(t, 2, p.unigram["O"])
	assert.Equal(t, 1, p.unigram["L"])

	assert.Contains(t, p.bigram, "FO")
	assert.Contains(t, p.bigram, "OO")
	assert.Contains(t, p.bigram, "OL")

	assert.Contains(t, p.trigram, "FOO")
	assert.Contains(t, p.trigram, "OOL")
}

func TestWordScore(t *testing.T) {
	p := NewPronouncable()
	p.AddWord("woo")

	assert.Equal(t, 1.5, p.WordScore("wo"))
	assert.Equal(t, 0.5, p.WordScore("w"))
	assert.Equal(t, 0.0, p.WordScore("woz"), "unknown character should give score 0")
}

func TestIsPronouncable(t *testing.T) {
	p := NewPronouncable()
	p.AddWord("wooh")

	assert.Equal(t, 2.25, p.WordScore("woo"))
	assert.True(t, p.IsPronouncable("woo", 2.25))
	assert.True(t, p.IsPronouncable("woo", 2.2))
	assert.False(t, p.IsPronouncable("woo", 2.3))
}

func TestAddWordList(t *testing.T) {
	p := NewPronouncable()
	p.AddWordList(bytes.NewBufferString("WOO MOO"))

	assert.Equal(t, 1, p.unigram["W"])
	assert.Equal(t, 4, p.unigram["O"])
	assert.Equal(t, 1, p.unigram["M"])
}

func TestSetWeights(t *testing.T) {
	p := NewPronouncable()
	p.AddWord("wools bools")

	p.SetWeights(1, 0, 0)
	assert.Equal(t, 0.375, p.WordScore("woo"))
	p.SetWeights(0, 1, 0)
	assert.Equal(t, 0.25, p.WordScore("woo"))
	p.SetWeights(0, 0, 1)
	assert.Equal(t, 0.08333333333333333, p.WordScore("woo"))
}

