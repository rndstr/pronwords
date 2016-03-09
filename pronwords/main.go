package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/rndstr/pronwords"
	"math"
	"os"
	"sort"
)

const (
	optTopDefault = 0
	optPercentileDefault = 0
	optThresholdDefault = -1.0
)

func createGenerator(c *cli.Context) (*pronwords.Generator, error) {
	g := pronwords.NewGenerator(c.Int("length"))
	if c.String("characters") != "" {
		g.SetCharacters(c.String("characters"))
	}

	return g, nil
}

func createPronouncable(c *cli.Context) (*pronwords.Pronounceable, error) {
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

func determineThreshold(percentile, top int, g *pronwords.Generator, p *pronwords.Pronounceable) float64 {
	count := g.Count()
	size := 0
	if percentile > 0 {
		fmt.Fprint(os.Stderr, "        // determining threshold for ", percentile, "th percentile… ")
		size = int(math.Ceil(float64(100.0 - percentile) / 100.0 * float64(count)))
	} else if top > 0 {
		fmt.Fprint(os.Stderr, "        // determining threshold for top ", top, "… ")
		size = top
		if top > count {
			fmt.Fprintln(os.Stderr, "none (--top larger than total count)")
			return optThresholdDefault
		}
	}

	scores := make([]float64, size)
	minScoreIndex := -1
	for word := g.Next(); word != ""; word = g.Next() {
		s := p.WordScore(word)
		if minScoreIndex == -1 {
			minScoreIndex = 0
		}
		if s > scores[minScoreIndex] {
			scores[minScoreIndex] = s
			minScoreIndex = 0
			for i := range scores {
				if scores[i] < scores[minScoreIndex] {
					minScoreIndex = i
				}
			}
		}
	}

	sort.Sort(sort.Float64Slice(scores))

	fmt.Fprintln(os.Stderr, scores[0])
	g.Reset() // start over

	return scores[0]
}

func main() {
	app := cli.NewApp()
	app.Name = "pronwords"
	app.Version = "0.1.1"
	app.Usage = "find pronounceable words"
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

		threshold := c.Float64("threshold")
		percentile := c.Int("percentile")
		top := c.Int("top")
		if (percentile != optPercentileDefault && top != optTopDefault) || (percentile != optPercentileDefault && threshold != optThresholdDefault) || (threshold != optThresholdDefault && top != optTopDefault) {
			fmt.Fprintln(os.Stderr, "error: --top --percentile --threshold are mutually exclusive, you cannot combine them")
			return
		}

		// percentile or top? -> threshold
		if percentile < 0 || percentile >= 100 {
			fmt.Fprintln(os.Stderr, "error: --percentile, please provide a number between 0 and 99")
			return
		}
		if top < 0 {
			fmt.Fprintln(os.Stderr, "error: --top, please provide a number greater than 0")
			return
		}

		if percentile > 0 || top > 0 {
			threshold = determineThreshold(percentile, top, g, p)
		}

		// print scores
		var match int
		max := -1.0
		min := math.MaxFloat64
		hideScores := c.Bool("hide-scores")
		for word := g.Next(); word != ""; word = g.Next() {
			s := p.WordScore(word)
			if s >= threshold {
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
		fmt.Fprint(os.Stderr, "        // Matched ", match, "/", count, " words (", int(100.0 * float64(match) / float64(count)), "%)")
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
		cli.Float64Flag{
			Name:  "threshold, t",
			Value: optThresholdDefault,
			Usage: "only show words with score >=threshold",
		},
		cli.IntFlag{
			Name:  "percentile, p",
			Value: optPercentileDefault,
			Usage: "only show scores in given percentile",
		},
		cli.IntFlag{
			Name: "top",
			Value: optTopDefault,
			Usage: "only print top N words",
		},
		cli.BoolFlag{
			Name: "hide-scores",
			Usage: "do not display scores",
		},
	}

	app.Run(os.Args)
}
