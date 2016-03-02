# pronword [![GoDoc](https://godoc.org/github.com/rndstr/pronwords?status.svg)](https://godoc.org/github.com/rndstr/pronwords) [![Build Status](https://travis-ci.org/rndstr/pronwords.svg?branch=master)](https://travis-ci.org/rndstr/pronwords) [![Coverage](http://gocover.io/_badge/github.com/rndstr/pronwords?0)](http://gocover.io/github.com/rndstr/pronwords)

`pronwords` is a Go library and CLI about pronounceable words.

## Introduction
`pronwords` comes with a word generator and a scoring mechanism to compare the pronounceableness between words. It acts
upon a corpus from which it will attempt to deduct a score.

The library does not come with a corpus but any list of words will suffice. A good starting point is the
[5 character word list](http://www.poslarchive.com/math/scrabble/lists/common-5.html) from
http://www.poslarchive.com/math/scrabble/lists/index.html

The package also comes with a CLI to use the library directly to generate pronounceable words. See `CLI Examples`

## Library

### Installation

    go get github.com/rndstr/pronwords

### Usage
See [documentation](https://godoc.org/github.com/rndstr/pronwords)

## CLI

### Installation
…

### Usage

    $ go run main.go help
    NAME:
       pronwords - find pronounceable words

    USAGE:
       main [global options] command [command options] [arguments...]
       
    VERSION:
       0.2.0
       
    COMMANDS:
       help, h      Shows a list of commands or help for one command
       
    GLOBAL OPTIONS:
       --length, -n "0"                                     generated word length
       --characters, -c "abcdefghijklmnopqrstuvwxyz"        generated words character set
       --corpus, -i                                         path to the corpus containing text to learn from
       --threshold, -t                                      only show words with score >=threshold
       --percentile, -p "0"                                 only show scores in given percentile
       --top "0"                                            only print top N words
       --hide-scores                                        do not display scores
       --help, -h                                           show help
       --version, -v                                        print the version
       
### Examples

```
$ go run main.go --corpus common5.words --length 5 --characters pronwords --top 20
        // determining threshold for top 20… 0.9447127049619598
proon 0.9502984259577109
prons 1.0502507547377256
pross 0.9876569715273623
poros 0.9447127049619598
poons 0.9492079539463733
rrons 0.9483483687726579
roros 0.9702004445676686
roops 0.9492949680898136
roors 0.9548376449101342
roons 1.0888509446936347
ronds 0.9774030670997887
ronso 0.954709606223315
ronss 0.9635767466471651
rosor 0.9521906700242896
orons 1.0419937814516256
oross 0.9793999982412626
ooros 0.9619801323607914
drons 0.9933544797710215
spros 0.9625463836278191
soros 0.9889266106841851

        // Matched 20/16807 words (0%)
        // Scores min = 0 max = 1.0888509446936347
```


```
$ go run main.go --corpus common5.words --length 3 --characters pronwords --percentile 95
        // determining threshold for 95th percentile… 0.8070014351051092
pro 0.8360481040136291
ror 0.8170963571976495
roo 0.9931571597245344
ron 1.016804606308029
row 0.858306126618284
ros 1.1708372168279426
rns 0.8370574373840399
rds 0.8090073862977848
ops 0.8070014351051092
oro 0.8195341574414293
ors 0.9322420398873033
oon 0.8862087307527223
oos 0.8461774143319966
ons 1.0861133883127514
oss 0.8467705707504732
nds 0.8510672335816807
sor 0.8545486894306823
son 0.8259464362579569

        // Matched 18/343 words (5%)
        // Scores min = 0 max = 1.1708372168279426
```
