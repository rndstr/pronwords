package main

import (
    "fmt"
    "strings"
    "io/ioutil"
)

const CHARACTERS = "rolandschilter"

func readwords(path string) []string {
    dat, err := ioutil.ReadFile(path)
    if err != nil {
        panic(err)
    }

    return strings.Split(strings.TrimSpace(string(dat)), "\n")
}

func main() {
    // remove duplicates in characters
    seen := map[byte]bool{}
    characters := []byte{}
    upperCharacters := strings.ToUpper(CHARACTERS)
    for i := 0; i < len(upperCharacters); i += 1 {
        value := upperCharacters[i]
		if !seen[value] {
			characters = append(characters, value)
			seen[value] = true
		}
	}
    fmt.Printf("%#v\n", string(characters))

    words := readwords("words5.list")

    one := make(map[string]int)
    two := make(map[string]int)
    three := make(map[string]int)


    // precalculate occurences
    for _,c0 := range characters {
        one[string(c0)] = 1
        for _,c1 := range characters {
            two[string(c0) + string(c1)] = 1
            for _,c2 := range characters {
                three[string(c0) + string(c1) + string(c2)] = 1
            }
        }
    }

    for _,word := range words {
        wordlen := len(word)
        for i := 0; i < wordlen; i += 1 {
            // TODO check whether any of the chars is not in characters list. if so, ommit.
            if i <= wordlen - 3 {
                three[word[i:i+3]] += 1
            }
            if i <= wordlen - 2 {
                two[word[i:i+2]] += 1
            }
            one[word[i:i+1]] += 1
        }
    }

    onesum := 0
    twosum := 0
    threesum := 0

    for _,o := range one { onesum += o }
    for _,o := range two { twosum += o }
    for _,o := range three { threesum += o }

    var score float64

    var scores = make(map[string]float64)
    for _,c0 := range characters {
        for _,c1 := range characters {
            for _,c2 := range characters {
                for _,c3 := range characters {
                    word := string(c0) + string(c1) + string(c2) + string(c3)
                    score = wordscore(word, one, two, three, onesum, twosum, threesum)

                    if score > 0.4 {
                        fmt.Printf("%#v = %#v\n", word, score)
                        scores[word] = score
                    }
                }
            }
        }
    }

    //fmt.Printf("%#v\n", scores)
}

func wordscore(word string, one map[string]int, two map[string]int, three map[string]int, onesum int, twosum int, threesum int) float64 {
    score := 0.0

    for i := 0; i < len(word); i += 1 {
        if i <= len(word) - 3 {
            score += 5.0 * float64(three[word[i:i+3]]) / float64(threesum)
        }
        if i <= len(word) - 2 {
            score += float64(two[word[i:i+2]]) / float64(twosum)
        }
        score += float64(one[word[i:i+1]]) / float64(onesum)
    }

    return score
}
