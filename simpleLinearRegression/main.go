// https://godoc.org/gonum.org/v1/gonum/stat

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
)

func main() {
	f, err := os.Open("../data/Advertising_ISLR.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var xs, ys []float64

	reader := csv.NewReader(f)

	reader.FieldsPerRecord = 5
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, record := range trainingData {

		// Skip the header.
		if i == 0 {
			continue
		}

		// TV column is our X
		xv, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Salary column is our y
		yv, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Fatal(err)
		}

		xs = append(xs, xv)
		ys = append(ys, yv)
	}

	// sort.Float64s(xs)
	// sort.Float64s(ys)

	n := len(xs)

	xmean := stat.Mean(xs, nil)
	ymean := stat.Mean(ys, nil)
	uxvariance := stat.Variance(xs, nil)

	fmt.Println("\nStats needed for Simple Linear Regression")
	fmt.Printf("\nN: %v\n", n)
	fmt.Printf("X mean, y mean: %v, %v\n", xmean, ymean)
	fmt.Printf("unbiased x variance:  %v\n\n", uxvariance)

	covar := covariance(xs, xmean, ys, ymean, float64(n-1))
	fmt.Println("Covariance function allows N - 1")
	fmt.Printf("unbiased covariance: %v\n\n", covar)
	b1 := covar / uxvariance
	fmt.Println("b1 = unbiased covar / unbiased x variance")
	fmt.Printf("b1: %v\n", b1)
	fmt.Println("b0 = ymean - (b1*xmean)")
	fmt.Printf("b0: %v\n", ymean-(b1*xmean))
}

// Ability to be unbiased
func covariance(xs []float64, xmean float64, ys []float64, ymean float64, n float64) float64 {
	var sum float64
	for i := 0; i < len(xs); i++ {
		sum += ((xs[i] - xmean) * (ys[i] - ymean))
	}
	return sum / n
}
