package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
    "sync"
	"github.com/antzucaro/matchr"
)

type Word struct {
	word     string
	distance int
}

func main() {

	misspell, err := scanLines("./data/misspell.txt")
	if err != nil {
		panic(err)
	}

	// for _, line := range misspell {
	// 	fmt.Println(line)
	// }

	dict, err := scanLines("./data/dict.txt")
	if err != nil {
		panic(err)
	}

	// for _, line := range dict {
	// 	fmt.Println(line)
	// }

	correct, err := scanLines("./data/correct.txt")
	if err != nil {
		panic(err)
	}
	// for _, line := range correct {
	// 	fmt.Println(line)
    // }
    
    // TESTING 
    // for i := 0; i < len(misspell); i++ {
    //     if misspell[i] != correct[i]{
    //         fmt.Printf("%d: missp:%s correct:%s\n", i, misspell[i], correct[i])
    //     }
    // }
    var wg sync.WaitGroup
    queue := make(chan Word, 1)
        
    corrects := 0
	for i := 0; i < len(misspell); i++ {
        
        var l []Word
        wg.Add(len(dict))
        for _, word := range dict {
            // go calculate_levenstein(&wg, index, word, l, misspell)
            go func(i int, word string, l []Word, misspell []string) {
            // defer wg.Done()  <- will result in the last int to be missed in the receiving channel
            queue <- Word{word, matchr.Levenshtein(misspell[i], word)}  
        }(i, word, l, misspell)
        }

        go func() {
        // defer wg.Done() <- Never gets called since the 100 `Done()` calls are made above, resulting in the `Wait()` to continue on before this is executed
        for t := range queue {

            // slice = append(slice, t)

            if len(t.word)>0{
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
        // fmt.Println(l)
    // //descending
    //     sort.Sort(sort.Reverse(ByDistance(l)))
        var dis int
        if len(l) > 0 {
        dis = l[0].distance
        }
       
        // fmt.Println("lwoest dist is ", dis)
        var predicts []string

        for _, w := range l{
            if w.distance == dis{
            predicts = append(predicts, w.word)
            } else {
                break
            }
        }
        // fmt.Println(predicts)
        stringExists := stringInSlice(correct[i], predicts)

        if stringExists == true {
            corrects = corrects + 1
        }

        // fmt.Println( countInSlice(correct[i], predicts))
        // accuracy := float64(corrects/len(predicts))
        // fmt.Printf("accuracy: %f\n", accuracy)
            // precision = float(is_correct.token.count()) / float(candidates.candidates_len.sum())
        // print("{}/{} precision: {}".format(i, len(test_word), predicts.count(correct_word[i]) / len(predicts)))
        // precision :=  float64(countInSlice(correct[i], predicts)/len(predicts))
        
        // fmt.Printf("precision: %f\n", precision)
    }
    // fmt.Printf("recall: %f/n", corrects/len(misspell))
    fmt.Println("DONE")
     
}

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
            count ++
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
