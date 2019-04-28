package stringmatching

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type ByDistance []Word

type Word struct {
	Word     string
	Distance int
}

type BestMatch struct {
	Word     string   `json:"word"`
	Distance int      `json:"distance"`
	Matches  []string `json:"matches"`
}

type Metrics struct {
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	Algorithm string  `json:"algorithm"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func DisplayPerformanceStats(precision float64, recall float64) {
	fmt.Println()
	fmt.Println("============ EVALUATION METRICS ============")
	fmt.Printf("precision: %f\n", precision)
	fmt.Printf("recall: %f\n", recall)
}

func WriteToFile(results interface{}, filePath string) {
	jsonData, _ := json.MarshalIndent(results, "", "\t")

	// write to JSON file
	jsonFile, err := os.Create(filePath)
	check(err)
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("Output JSON data written to ", jsonFile.Name())
}

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func CountInSlice(a string, list []string) int {
	count := 0
	for _, b := range list {
		if b == a {
			count++
		}
	}
	return count
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ScanLines(path string) ([]string, error) {

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

func UniqueSlice(stringSlice []string) []string {
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
