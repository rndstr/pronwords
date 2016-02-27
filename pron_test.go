package pronwords

import (
        "testing"

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

