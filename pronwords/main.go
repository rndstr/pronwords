package main

import (
    "fmt"
    "os"
    "github.com/rndstr/pronwords"
    "github.com/codegangsta/cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "pronwords"
    app.Version = "0.1.0"
    app.Usage = "find pronouncable words"
    app.Action = func(c *cli.Context) {
        // Validate
        if c.String("corpus") == "" {
            fmt.Fprintln(os.Stderr, "missing required option --corpus (see `pronwords help` for usage)")
            return
        }
        if c.Int("length") == 0 {
            fmt.Fprintln(os.Stderr, "missing required option --length (see `pronwords help` for usage)")
            return
        }

        g := pronwords.NewGenerator(c.Int("length"))
        if c.String("characters") != "" {
            g.SetCharacters(c.String("characters"))
        }

        file, err := os.Open(c.String("corpus"))
        if err != nil {
            fmt.Fprintln(os.Stderr, "error opening corpus: ", err)
            return
        }
        p := pronwords.NewPronouncable()
        p.AddWordList(file)

        for word := g.Next() ; word != ""; word = g.Next()  {
            s := p.WordScore(word)
            fmt.Printf("%v %v\n", word, s)
        }
    }
    app.Flags = []cli.Flag {
        cli.IntFlag{
            Name: "length, n",
            Usage: "generated word length",
        },
        cli.StringFlag{
            Name: "characters, c",
            Value: "abcdefghijklmnopqrstuvwxyz",
            Usage: "generated words character set",
        },
        cli.StringFlag{
            Name: "corpus, i",
            Usage: "path to the corpus containing text to learn from",
        },
    }

    app.Run(os.Args)
}

