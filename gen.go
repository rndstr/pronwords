package pronwords

import (
    "strings"
)

type Generator struct {
    characters []byte
    length int
    current []byte
    positions []int
}

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

func NewGenerator(characters string, length int) *Generator {
    bytes := deduplicate(characters)

    return &Generator{
        characters: bytes,
        length: length,
        positions: make([]int, len(bytes)),
    }
}

func (g *Generator) Next() string {
    if g.current == nil {
        // first word is first character at each position
        g.current = make([]byte, g.length)
        for i := 0; i < g.length; i += 1 {
            g.current[i] = g.characters[0]
        }
    } else {
        inc := g.length - 1
        for inc >= 0 && g.positions[inc] == len(g.characters) - 1 {
            inc--
        }

        // done?
        if inc == -1 {
            return ""
        }

        g.positions[inc] += 1
        g.current[inc] = g.characters[g.positions[inc]]

        if inc != g.length - 1 {
            g.positions[inc+1] = 0
            g.current[inc+1] = g.characters[0]
        }
    }

    return string(g.current)
}

