package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
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

	reader := csv.NewReader(f)

	reader.FieldsPerRecord = 5
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	x := make([][]float64, reader.FieldsPerRecord-1)
	var y []float64

	for f := 0; f < reader.FieldsPerRecord; f++ {
		for i, record := range trainingData {

			// Skip header
			if i == 0 {
				continue
			}

			if f == 0 {
				// Fill the bias with ones
				x[f] = append(x[f], 1)
			} else if f < reader.FieldsPerRecord-1 {
				// Populate Features
				xv, err := strconv.ParseFloat(record[f], 64)
				if err != nil {
					log.Fatal(err)
				}
				x[f] = append(x[f], xv)
			} else {
				// Populate Labels
				yv, err := strconv.ParseFloat(record[f], 64)
				if err != nil {
					log.Fatal(err)
				}
				y = append(y, yv)
			}
		}
	}

	// Standardize X to avoid under / over flow
	x = standardize(x)
	iterations := 1000
	learningRate := 0.1

	coefficients := train(x, y, iterations, learningRate)
	fmt.Println("TV, Radio, Newspaper features have been standardized")
	fmt.Println("Coefficients : ", coefficients)
}

func standardize(x [][]float64) [][]float64 {
	for i := 1; i < len(x); i++ {
		xmean := stat.Mean(x[i], nil)
		// Unbiased variance
		uxvariance := stat.Variance(x[i], nil)
		for j := 0; j < len(x[i]); j++ {
			// (x - mean) / std
			x[i][j] = (x[i][j] - xmean) / math.Pow(uxvariance, 0.5)
		}
	}
	return x
}

func train(x [][]float64, y []float64, iterations int, learningRate float64) []float64 {

	coefficients := make([]float64, len(x))

	for i := 0; i < iterations; i++ {
		gradientDescent(x, y, coefficients, learningRate)
	}

	return coefficients
}

func gradientDescent(x [][]float64, y []float64, coefficients []float64, learningRate float64) []float64 {

	residuals := make([]float64, len(y))
	for i := 0; i < len(y); i++ {
		prediction := 0.0
		for j := 0; j < len(coefficients); j++ {
			prediction += coefficients[j] * x[j][i]
		}
		residuals[i] = y[i] - prediction
	}

	for i := 0; i < len(residuals); i++ {
		// https://mccormickml.com/2014/03/04/gradient-descent-derivation/
		// Andrew Ng
		// We multiply our MSE cost function by 1/2 so that when we take the derivative, the 2s cancel out
		// Multiplying the cost function by a scalar does not affect the location of its minimum
		coefficients[0] += learningRate * residuals[i] / float64(len(residuals))
		for j := 1; j < len(x); j++ {
			coefficients[j] += learningRate * (residuals[i] * x[j][i]) / float64(len(residuals))
		}
	}

	return coefficients
}
