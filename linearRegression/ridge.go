package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/berkmancenter/ridge"
	"github.com/gonum/matrix/mat64"
)

func ridgegression() {

	f, err := os.Open("../data/training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 4

	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	xData := make([]float64, 4*len(rawCSVData))
	yData := make([]float64, len(rawCSVData))

	var xIndex int
	var yIndex int

	for idx, record := range rawCSVData {

		// Skip header
		if idx == 0 {
			continue
		}

		for i, val := range record {

			valParsed, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal("Could not parse float value")
			}

			if i < reader.FieldsPerRecord-1 {

				// Add bias
				if i == 0 {
					xData[xIndex] = 1
					xIndex++
				}

				xData[xIndex] = valParsed
				xIndex++
			}

			if i == reader.FieldsPerRecord-1 {
				yData[yIndex] = valParsed
				yIndex++
			}

		}
	}

	xTrain := mat64.NewDense(len(rawCSVData), reader.FieldsPerRecord, xData)
	yTrain := mat64.NewVector(len(rawCSVData), yData)

	r := ridge.New(xTrain, yTrain, 1.0)

	// Train
	r.Regress()

	// b0 := r.Coefficients.At(0, 0)
	// b1 := r.Coefficients.At(1, 0)
	// b2 := r.Coefficients.At(2, 0)
	// b3 := r.Coefficients.At(3, 0)
	fmt.Print("\nCoefficients:\n")
	for i := 0; i < reader.FieldsPerRecord; i++ {
		fmt.Printf("b%v: %v\n", i, r.Coefficients.At(i, 0))
	}

	ftest, err := os.Open("../data/test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	readtest := csv.NewReader(ftest)

	readtest.FieldsPerRecord = 4
	testData, err := readtest.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var mAE, mSE float64
	xVals := make([]float64, readtest.FieldsPerRecord-1)

	for i, record := range testData {

		// Skip the header.
		if i == 0 {
			continue
		}

		yObserved, err := strconv.ParseFloat(record[readtest.FieldsPerRecord-1], 64)
		if err != nil {
			log.Fatal(err)
		}

		for f := 0; f < len(xVals); f++ {

			// Skip the header.
			if i == 0 {
				continue
			}

			xv, err := strconv.ParseFloat(record[f], 64)
			if err != nil {
				log.Fatal(err)
			}
			xVals[f] = xv
		}

		yPredicted := predictRidge(xVals, r.Coefficients)
		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
		mSE += math.Pow(yObserved-yPredicted, 2) / float64(len(testData))
	}

	fmt.Printf("MAE = %0.2f\n", mAE)
	fmt.Printf("MSE = %0.2f\n\n", mSE)
}

func predictRidge(xVals []float64, coeffs *mat64.Vector) float64 {
	var prediction float64
	prediction = coeffs.At(0, 0)
	for i := 0; i < len(xVals); i++ {
		prediction += coeffs.At(i+1, 0) * xVals[i]
	}
	return prediction
}
