package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/antzucaro/matchr"
	. "github.com/emxqm/lex-analysis/pkg"
)

func main() {
	//start timer
	start := time.Now()

	fmt.Println("Beginning Levenshtein evaluation...")

	//load data as list of strings
	misspell, err := ScanLines("./data/misspell.txt")
	if err != nil {
		panic(err)
	}

	dict, err := ScanLines("./data/dict.txt")
	if err != nil {
		panic(err)
	}

	correct, err := ScanLines("./data/correct.txt")
	if err != nil {
		panic(err)
	}
	//get unique list of mispelled words
	uniqueMispell := UniqueSlice(misspell)

	// mapping of mispelled words to correct canonical form
	m := make(map[string]string)
	for index, misspelledWord := range misspell {
		m[misspelledWord] = correct[index]
	}

	corrects := 0
	totalPredictions := 0
	maxEditDistance := 2
	maxResults := 100
	arbitraryHighValue := 9999
	var results []BestMatch
	var performance []Metrics

	for i := 0; i < len(uniqueMispell); i++ {
		var l []Word

		//loop over dictionary calculating Levenshtein distance on each pairing
		for _, word := range dict {
			l = append(l, Word{word, matchr.Levenshtein(uniqueMispell[i], word)})
		}
		//sort by lowest distance
		sort.Sort(ByDistance(l))

		dis := arbitraryHighValue
		if len(l) > 0 {
			dis = l[0].Distance
		}
		// fmt.Println("lwoest dist is ", dis)
		var matches []string

		if dis <= maxEditDistance {
			for _, w := range l[:maxResults] {
				if w.Distance == dis {
					matches = append(matches, w.Word)
					totalPredictions++
				} else {
					break
				}
			}
		}
		// fmt.Println(matches)

		stringExists := StringInSlice(m[uniqueMispell[i]], matches)

		if stringExists == true {
			corrects++
		}

		results = append(results, BestMatch{uniqueMispell[i], dis, matches})
	}
	// save all results to a file
	WriteToFile(results, "./output/results-Levenshtein.json")
	// metrics calculations
	recall := float64(corrects) / float64(len(uniqueMispell))
	precision := float64(corrects) / float64(totalPredictions)
	DisplayPerformanceStats(precision, recall)

	performance = append(performance,
		Metrics{Precision: precision, Recall: recall, Algorithm: "Levenshtein"})

	WriteToFile(performance, "./output/metrics-Levenshtein.json")
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Total Execution Time: ", elapsed.String())
}
