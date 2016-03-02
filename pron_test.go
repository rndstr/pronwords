package pronwords

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddWord(t *testing.T) {
	p := NewPronounceable()
	p.AddWord("fool")

	assert.Equal(t, 1, p.unigram["f"])
	assert.Equal(t, 2, p.unigram["o"])
	assert.Equal(t, 1, p.unigram["l"])

	assert.Contains(t, p.bigram, "fo")
	assert.Contains(t, p.bigram, "oo")
	assert.Contains(t, p.bigram, "ol")

	assert.Contains(t, p.trigram, "foo")
	assert.Contains(t, p.trigram, "ool")
}

func TestWordScore(t *testing.T) {
	p := NewPronounceable()
	p.AddWord("woo")

	assert.Equal(t, 1.5, p.WordScore("wo"))
	assert.Equal(t, 0.5, p.WordScore("w"))
	assert.Equal(t, 0.0, p.WordScore("woz"), "unknown character should give score 0")
}

func TestIsPronouncable(t *testing.T) {
	p := NewPronounceable()
	p.AddWord("wooh")

	assert.Equal(t, 2.25, p.WordScore("woo"))
	assert.True(t, p.IsPronounceable("woo", 2.25))
	assert.True(t, p.IsPronounceable("woo", 2.2))
	assert.False(t, p.IsPronounceable("woo", 2.3))
}

func TestAddWordList(t *testing.T) {
	p := NewPronounceable()
	p.AddWordList(bytes.NewBufferString("WOO MOO"))

	assert.Equal(t, 1, p.unigram["w"])
	assert.Equal(t, 4, p.unigram["o"])
	assert.Equal(t, 1, p.unigram["m"])
}

func TestSetWeights(t *testing.T) {
	p := NewPronounceable()
	p.AddWord("wools bools")

	p.SetWeights(1, 0, 0)
	assert.Equal(t, 0.375, p.WordScore("woo"))
	p.SetWeights(0, 1, 0)
	assert.Equal(t, 0.25, p.WordScore("woo"))
	p.SetWeights(0, 0, 1)
	assert.Equal(t, 0.08333333333333333, p.WordScore("woo"))
}
