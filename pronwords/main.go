package main

import (
    "fmt"
    "github.com/rndstr/pronwords"
)

func main() {
    g := pronwords.NewGenerator("abc", 3)
    for word := g.Next() ; word != ""; word = g.Next()  {
        fmt.Println(word)
    }

    p := pronwords.NewPronouncable()
    p.AddWordList("words5.list")
    fmt.Println("score = ", p.WordScore("kzbd"))
}
