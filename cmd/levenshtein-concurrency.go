package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/antzucaro/matchr"
)

type Word struct {
	word     string
	distance int
}

type BestMatch struct {
	Word     string   `json:"word"`
	Distance int      `json:"distance"`
	Matches  []string `json:"matches"`
}

func main() {

	start := time.Now()

	misspell, err := scanLines("./data/misspell.txt")
	if err != nil {
		panic(err)
	}

	uniqueMispell := uniqueSlice(misspell)
	fmt.Println(len(uniqueMispell))

	dict, err := scanLines("./data/dict.txt")
	if err != nil {
		panic(err)
	}

	correct, err := scanLines("./data/correct.txt")
	if err != nil {
		panic(err)
	}

	uniqueCorrect := uniqueSlice(correct)
	fmt.Println(len(uniqueCorrect))

	m := make(map[string]string)
	for index, misspelledWord := range misspell {
		m[misspelledWord] = correct[index]
	}

	var wg sync.WaitGroup
	queue := make(chan Word, len(dict))

	corrects := 0
	totalPredictions := 0

	for i := 0; i < len(uniqueMispell); i++ {

		var l []Word
		wg.Add(len(dict))
		for _, word := range dict {
			// go calculate_levenstein(&wg, index, word, l, misspell)
			go func(i int, word string, l []Word, uniqueMispell []string) {
				// defer wg.Done()  <- will result in the last int to be missed in the receiving channel
				queue <- Word{word, matchr.Levenshtein(uniqueMispell[i], word)}
			}(i, word, l, uniqueMispell)
		}

		go func() {
			// defer wg.Done() <- Never gets called since the 100 `Done()` calls are made above, resulting in the `Wait()` to continue on before this is executed
			for t := range queue {
				if len(t.word) > 0 {
					l = append(l, t)
				}
				wg.Done()
			}
		}()
		// fmt.Println("Main: Waiting for workers to finish")
		wg.Wait()
		// fmt.Println("Main: Completed")
		//ascending
		sort.Sort(ByDistance(l))
		fmt.Println(l)
		// //descending
		//     sort.Sort(sort.Reverse(ByDistance(l)))
		dis := 9999
		if len(l) > 0 {
			dis = l[0].distance
		}

		// fmt.Println("lwoest dist is ", dis)
		var matches []string
		if dis <= 2 {
			for _, w := range l[:100] {
				if w.distance == dis {
					matches = append(matches, w.word)
					totalPredictions++
				} else {
					break
				}
			}
		}
		// fmt.Println(matches)
		stringExists := stringInSlice(m[uniqueMispell[i]], matches)

		if stringExists == true {
			corrects++
		}

		// results = append(results, BestMatch{uniqueMispell[i], dis, matches})
	}
	// writeToFile(results)
	fmt.Println("TOTAL CORRECTS: ", corrects)
	fmt.Println("TOTAL PREDICTIONS: ", totalPredictions)
	recall := float64(corrects / len(uniqueMispell))
	fmt.Println("RECALL: ", recall)
	precision := float64(corrects / totalPredictions)
	fmt.Println("PRECISION: ", precision)
	displayPerformanceStats(precision, recall)
	fmt.Println("DONE")
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Total Execution Time: ", elapsed.String())
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func displayPerformanceStats(precision float64, recall float64) {
	fmt.Println("**** EVALUATION METRICS ****")
	fmt.Printf("precision: %f\n", precision)
	fmt.Printf("recall: %f/n", recall)
}

// func writeToFile(results []BestMatch) {
// 	jsonData, _ := json.MarshalIndent(results, "", "\t")

// 	// write to JSON file
// 	jsonFile, err := os.Create("./results.json")
// 	check(err)
// 	defer jsonFile.Close()

// 	jsonFile.Write(jsonData)
// 	jsonFile.Close()
// 	fmt.Println("Output JSON data written to ", jsonFile.Name())
// }

func calculate_levenstein(wg *sync.WaitGroup, i int, word string, l []Word, misspell []string) {
	defer wg.Done()

	l = append(l, Word{word, matchr.Levenshtein(misspell[i], word)})
	fmt.Printf("Worker %v: Finished\n", i)
}

type ByDistance []Word

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].distance < a[j].distance }

func countInSlice(a string, list []string) int {
	count := 0
	for _, b := range list {
		if b == a {
			count++
		}
	}
	return count
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func scanLines(path string) ([]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func uniqueSlice(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
