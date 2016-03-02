package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/rndstr/pronwords"
	"math"
	"os"
	"sort"
	"strconv"
)

func createGenerator(c *cli.Context) (*pronwords.Generator, error) {
	g := pronwords.NewGenerator(c.Int("length"))
	if c.String("characters") != "" {
		g.SetCharacters(c.String("characters"))
	}

	return g, nil
}

func createPronouncable(c *cli.Context) (*pronwords.Pronounceable, error) {
	// Validate
	if c.String("corpus") == "" {
		return nil, errors.New("missing required option --corpus (see `pronwords help` for usage)")
	}
	if c.Int("length") == 0 {
		return nil, errors.New("missing required option --length (see `pronwords help` for usage)")
	}

	p := pronwords.NewPronounceable()

	file, err := os.Open(c.String("corpus"))
	if err != nil {
		return p, errors.New("error opening corpus: " + err.Error())
	}

	p.AddWordList(file)

	return p, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "pronwords"
	app.Version = "0.2.0"
	app.Usage = "find pronouncable words"
	app.Action = func(c *cli.Context) {
		g, err := createGenerator(c)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		p, err := createPronouncable(c)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		hasThreshold := c.String("threshold") != ""
		threshold := 0.0
		if hasThreshold {
			threshold, err = strconv.ParseFloat(c.String("threshold"), 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: --threshold, please provide a decimal number with '.' as decimal point")
				return
			}
		}

		// percentile or top? -> threshold
		percentile := c.Int("percentile")
		if percentile < 0 || percentile >= 100 {
			fmt.Fprintln(os.Stderr, "error: --percentile, please provide a number between 0 and 99")
			return
		}
		top := c.Int("top")
		if top < 0 {
			fmt.Fprintln(os.Stderr, "error: --top, please provide a number greater than 0")
			return
		}

		if percentile != 0 && top != 0 {
			fmt.Fprintln(os.Stderr, "error: --top -- percentile, you cannot use both")
			return

		}

		if percentile > 0 || top > 0 {
			if percentile > 0 {
				fmt.Fprint(os.Stderr, "        // determining threshold for ", percentile, "th percentile… ")
			} else if top > 0 {
				fmt.Fprint(os.Stderr, "        // determining threshold for top ", top, "… ")

			}
			scores := make([]float64, 0)
			for word := g.Next(); word != ""; word = g.Next() {
				if s := p.WordScore(word); !hasThreshold || s >= threshold {
					scores = append(scores, p.WordScore(word)) // let's hope append() increases capacity not one by one?
				}
			}

			sort.Sort(sort.Float64Slice(scores))
			if percentile > 0 {
				threshold = scores[int(math.Floor(float64(percentile)/100.0*float64(len(scores))))]
				hasThreshold = true
			} else if top > 0 && top < len(scores) {
				threshold = scores[len(scores) - top]
				hasThreshold = true
			}
			fmt.Fprintln(os.Stderr, threshold)

			g.Reset() // start over
		}

		// print scores
		var match int
		var max, min float64
		hideScores := c.Bool("hide-scores")
		for word := g.Next(); word != ""; word = g.Next() {
			if s := p.WordScore(word); !hasThreshold || s >= threshold {
				if s > max {
					max = s
				}
				if s < min {
					min = s
				}
				if hideScores {
					fmt.Println(word)
				} else {
					fmt.Printf("%v %v\n", word, s)
				}
				match++
			}
		}

		count := g.Count()
		fmt.Fprintln(os.Stderr)
		fmt.Fprint(os.Stderr, "        // Matched ", match, "/", count, " words (", int(100.0 * float64(match)/float64(count)), "%)")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "        // Scores min =", min, "max =", max)
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "length, n",
			Usage: "generated word length",
		},
		cli.StringFlag{
			Name:  "characters, c",
			Value: "abcdefghijklmnopqrstuvwxyz",
			Usage: "generated words character set",
		},
		cli.StringFlag{
			Name:  "corpus, i",
			Usage: "path to the corpus containing text to learn from",
		},
		cli.StringFlag{
			Name:  "threshold, t",
			Value: "",
			Usage: "only show words with score >=threshold",
		},
		cli.IntFlag{
			Name:  "percentile, p",
			Value: 0,
			Usage: "only show scores in given percentile",
		},
		cli.IntFlag{
			Name: "top",
			Value: 0,
			Usage: "only print top N words",
		},
		cli.BoolFlag{
			Name: "hide-scores",
			Usage: "do not display scores",
		},
	}

	app.Run(os.Args)
}
