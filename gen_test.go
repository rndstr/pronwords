package pronwords

import (
        "testing"
        "strings"

        "github.com/stretchr/testify/assert"
)

func TestCharSetSameLengthAsWordLength(t *testing.T) {
    assertGeneratedWords(t, "ab", "aa,ab,ba,bb", 2)
}

func TestCharSetLongerThanWordLength(t *testing.T) {
    assertGeneratedWords(
        t,
        "abc", // charset
        "aa,ab,ac,ba,bb,bc,ca,cb,cc",
        2, // word length
    )
}

func TestCharSetShorterThanWordLength(t *testing.T) {
    assertGeneratedWords(
        t,
        "ab", // charset
        "aaa,aab,aba,abb,baa,bab,bba,bbb",
        3, // word length
    )
}


func assertGeneratedWords(t *testing.T, characters, expectedCommaSeparated string, length int) {
    var (
        expectedList = strings.Split(strings.ToUpper(expectedCommaSeparated), ",")
        words = make(map[string]bool, len(expectedList))
    )

    for _,word := range expectedList {
        words[word] = true
    }

    g := NewGenerator(characters, length)
    count := g.Count()
    assert.Equal(t, len(words), count)

    for ; count > 0; count-- {
        word := g.Next()
        assert.Contains(t, words, word)

        delete(words, word)
    }

    assert.Equal(t, "", g.Next())
    assert.Equal(t, 0, len(words))
}
