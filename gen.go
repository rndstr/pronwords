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
    bytes := deduplicateAndUpperToBytes(characters)

    return &Generator{
        charset: bytes,
        length: length,
        levels: make([]int, length),
    }
}

// Next generates word. It returns the generated word in all uppercase as string.
func (g *Generator) Next() string {
    if g.word == nil {
        // Populate first word
        g.word = make([]byte, g.length)
        for i := 0; i < g.length; i += 1 {
            g.word[i] = g.charset[0]
        }
    } else {
        // Start at the very end of the levels
        col := g.length - 1
        // Move towards the front while skipping all the columns with highest level
        for col >= 0 && g.levels[col] == len(g.charset) - 1 {
            col--
        }

        // Are all columns at the top level?
        if col == -1 {
            return ""
        }

        // Increase level for the first level from the back not at the top.
        g.levels[col]++
        g.word[col] = g.charset[g.levels[col]]

        // All levels behind col are at the top. Since we modified the level at col,
        // we can reset all the levels behind col.
        for c := col + 1; c < g.length; c++ {
            g.levels[c] = 0
            g.word[c] = g.charset[0]
        }
    }

    return string(g.word)
}

// Count returns the number of words it will generate.
func (g *Generator) Count() int {
    return int(math.Pow(float64(len(g.charset)), float64(g.length)))
}


// Deduplicate removes duplicate characters in string. It also converts the string
// into a byte array.
func deduplicateAndUpperToBytes(str string) []byte {
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

