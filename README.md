# Twitter Lexical Normalisation

COMP90049 Knowledge Technologies, Semester 1 2019 Project 1 Assignment - Lexical Normalisation of Twitter Data

## Overview

This project is a golang program that utilized several methods to find a correct form of correct words from a mispelled word in the dictionary.

In this project, the author want to compare token that comes from twitter text (contain some mispelled word) against a list of English words to find a correct form.

The author use four methods for approximate string matching in the script.

- Only use Levenshtein Distance
- Levenshtein Distance + Soundex algorithm
- Levenshtein Distance + Phonex algorithm
- Levenshtein Distance + NYSIIS algorithm

## Getting Started

Make sure you have go installed in your system.
https://golang.org/doc/install

## RUN ALL Programs

VERY important to clone this repository into your \$GOPATH, otherwise it will not run.

'dep ensure' will populate external libraries in a directory /vendor from existing Gopkg.toml and Gopkg.lock

'make' will run the default command to run all four methods from the Makefile

```
$ dep ensure -v
$ make
```

## RUN a particular program

```
make levenshtein
make soundex
make nysiis
make phonex
```

External packages used :
github.com/antzucaro/matchr

## /cmd

A directory to contain all programs.

## /data

Data set sourced from

Baldwin, Timothy, Marie Catherine de Marneffe, Bo Han, Young-Bum Kim, Alan Ritter,and Wei Xu (2015) Shared Tasks of the 2015 Workshop on Noisy User-generated Text:Twitter Lexical Normalization and Named Entity Recognition. InProceedings of the ACL2015 Workshop on Noisy User-generated Text, Beijing, China, pp. 126–135.

## /output

Evaluaton metrics output from running the programs in /cmd

## TODO:

levenshtein-concurrency.go is currently a WIP to improve the execution time calculating levenshtein distance by using go-routines. Not used in the report.
