package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func splitTestTrain() {

	df := getAllData()

	// Set the split proportions (4/5 = 80 20 split)
	trainingNum := (4 * df.Nrow()) / 5
	testNum := df.Nrow() / 5
	if trainingNum+testNum < df.Nrow() {
		trainingNum++
	}

	trainingIdx := make([]int, trainingNum)
	testIdx := make([]int, testNum)

	for i := 0; i < trainingNum; i++ {
		trainingIdx[i] = i
	}

	for i := 0; i < testNum; i++ {
		testIdx[i] = trainingNum + i
	}

	trainingDF := df.Subset(trainingIdx)
	testDF := df.Subset(testIdx)
	fmt.Println("DF Dims")
	fmt.Println(df.Dims())
	fmt.Println("trainingDF Dims")
	fmt.Println(trainingDF.Dims())
	fmt.Println("testDF Dims")
	fmt.Println(testDF.Dims())

	setMap := map[int]dataframe.DataFrame{
		0: trainingDF,
		1: testDF,
	}

	// Create training and test files
	for idx, setName := range []string{"../data/training.csv", "../data/test.csv"} {

		f, err := os.Create(setName)
		if err != nil {
			log.Fatal(err)
		}

		w := bufio.NewWriter(f)

		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Fatal(err)
		}
	}
}
