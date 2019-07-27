package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gota/gota/dataframe"
)

var (
	xv []string
)

func main() {
	for {
		// Describe = mean, std, etc; Plot = Histogram, Scatter Plot; Model = Regression Using Variables
		fmt.Println("Choose: Describe, Plot, Split, Model, Ridge")

		// Read in data from user input
		reader := bufio.NewReader(os.Stdin)
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Choose function to call from user input
		switch cmd := strings.TrimSpace(strings.ToLower(cmdString)); cmd {
		case "ridge":
			ridgegression()
		case "describe":
			descriptiveStats()
		case "plot":
			histogram()
			scatterplot()
		case "split":
			splitTestTrain()
		case "model":
			for {
				fmt.Println("Choose variables (use space between variables): TV, Radio, Newspaper")
				xReader := bufio.NewReader(os.Stdin)
				xString, err := xReader.ReadString('\n')
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				xSlice := strings.SplitAfter(xString, " ")
				ok := true

				// Input control
				for _, s := range xSlice {
					switch strings.TrimSpace(strings.ToLower(s)) {
					case "tv":
						xv = append(xv, strings.TrimSpace(strings.ToLower(s)))
					case "radio":
						xv = append(xv, strings.TrimSpace(strings.ToLower(s)))
					case "newspaper":
						xv = append(xv, strings.TrimSpace(strings.ToLower(s)))
					default:
						// Check for typos
						ok = false
					}
				}

				if ok == true {
					// If no typos break and run regression
					break
				} else {
					xv = nil
				}
			}

			xv = unique(xv)
			// fmt.Printf("%v", xv)
			model()
		default:
			fmt.Println("Hello World")
		}
	}
}

// Check slice for redundant values
func unique(stringSlice []string) []string {
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

func getAllData() dataframe.DataFrame {
	// Open the CSV file from the ISLR book
	data, err := os.Open("../data/Advertising_ISLR.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	// Create dataframe from the CSV file.
	df := dataframe.ReadCSV(data)

	// Return columns without the index
	return df.Select([]int{1, 2, 3, 4})
}
