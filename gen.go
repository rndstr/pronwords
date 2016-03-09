package pronwords

import (
	"math"
	"strings"
)

type Generator struct {
	charset []byte // The characters to choose from.
	word    []byte // The current word that has been returned by Next().
	length  int    // Length of the generated word (not necessarily the same as len(charset)).
	levels  []int  // Mapping of each character in `word` to the n-th character in `charset`.
}

func NewGenerator(length int) *Generator {

	return &Generator{
		charset: []byte("abcdefghijklmnopqrstuvwxyz"),
		length:  length,
		levels:  make([]int, length),
	}
}

// SetCharacters changes the default character set of available chars.
func (g *Generator) SetCharacters(characters string) {
	if g.word != nil {
		panic("SetCharacters() called after generation started. Either call Reset() or create a new instance.")
	}
	g.charset = deduplicateAndLowerToBytes(characters)
}

// Reset restores the initial state and makes the generator to start over at word 1.
func (g *Generator) Reset() {
	g.word = nil

	for i := 0; i < g.length; i++ {
		g.levels[i] = 0
	}
}

// Next generates a word. It returns the generated word in all lowercase as string, and an empty word denotes the end
// of the generation.
func (g *Generator) Next() string {
	if g.word == nil {
		// Populate first word
		g.word = make([]byte, g.length)
		for i := 0; i < g.length; i++ {
			g.word[i] = g.charset[0]
		}
	} else {
		// Start at the very end of the levels
		col := g.length - 1
		// Move towards the front while skipping all the columns with highest level
		for col >= 0 && g.levels[col] == len(g.charset)-1 {
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

// DeduplicateAndLowerToByts removes duplicate characters in a string. It also converts the string
// into a byte array with all lower case characters.
func deduplicateAndLowerToBytes(str string) []byte {
	seen := map[byte]bool{}
	bytes := []byte{}
	characters := strings.ToLower(str)
	for i := range characters {
		value := characters[i]
		if !seen[value] {
			bytes = append(bytes, value)
			seen[value] = true
		}
	}
	return bytes
}
