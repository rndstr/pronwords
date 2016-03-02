# pronword [![GoDoc](https://godoc.org/github.com/rndstr/pronwords?status.svg)](https://godoc.org/github.com/rndstr/pronwords) [![Build Status](https://travis-ci.org/rndstr/pronwords.svg?branch=master)](https://travis-ci.org/rndstr/pronwords) [![Coverage](http://gocover.io/_badge/github.com/rndstr/pronwords?0)](http://gocover.io/github.com/rndstr/pronwords)

`pronwords` is a Go library and CLI about pronounceable words.

## Introduction
`pronwords` comes with a word generator and a scoring mechanism to compare the pronounceableness between words. It acts
upon a corpus from which it will deduct the score.

The library does not come with a corpus but any list of words will suffice. A good starting point is the
[5 character word list](http://www.poslarchive.com/math/scrabble/lists/common-5.html) from
http://www.poslarchive.com/math/scrabble/lists/index.html

The package also comes with a CLI to use the library directly to generate pronounceable words. See Usage

## `pronwords` CLI

### Usage

    $ pronwords help
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
       

