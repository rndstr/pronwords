package pronwords

import (
        "strings"
        "math"
)

type Generator struct {
    charset []byte // The characters to choose from.
    word    []byte // The current word that has been returned by Next().
    length  int // Length of the generated word (not necessarily the same as len(charset)).
    levels  []int // Mapping of each character in `word` to the n-th character in `charset`.
}

func NewGenerator(characters string, length int) *Generator {
    bytes := deduplicate(characters)

    return &Generator{
        charset: bytes,
        length: length,
        levels: make([]int, len(bytes)),
    }
}

// Next returns the next generated word.
func (g *Generator) Next() string {
    if g.word == nil {
        // initialize word and first generation
        g.word = make([]byte, g.length)
        for i := 0; i < g.length; i += 1 {
            g.word[i] = g.charset[0]
        }
    } else {
        inc := g.length - 1
        for inc >= 0 && g.levels[inc] == len(g.charset) - 1 {
            inc--
        }

        if inc == -1 {
            return ""
        }

        g.levels[inc] += 1
        g.word[inc] = g.charset[g.levels[inc]]

        if inc != g.length - 1 {
            g.levels[inc+1] = 0
            g.word[inc+1] = g.charset[0]
        }
    }

    return string(g.word)
}

// Count returns the number of words it will generate.
func (g *Generator) Count() int {
    return math.Pow(len(charset), length)
}


// Deduplicate removes duplicate characters in string. It also converts the string
// into a byte array.
func deduplicate(str string) []byte {
    seen := map[byte]bool{}
    characters := []byte{}
    upperCharacters := strings.ToUpper(str)
    for i := 0; i < len(upperCharacters); i += 1 {
        value := upperCharacters[i]
        if !seen[value] {
            characters = append(characters, value)
            seen[value] = true
        }
    }
    return characters
}

