package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/sajari/regression"
)

func model() {
	// Open the training dataset file.
	f, err := os.Open("../data/training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new CSV reader reading from the opened file.
	reader := csv.NewReader(f)

	// Read in all of the CSV records
	reader.FieldsPerRecord = 4
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// In this case we are going to try and model our Sales
	// by the TV and Radio features plus an intercept.
	var r regression.Regression
	r.SetObserved("Sales")

	for i := range xv {
		r.SetVar(i, strings.Title(xv[i]))
	}

	var (
		tv        float64
		radio     float64
		newspaper float64
	)

	// Loop over the CSV records adding the training data.
	for i, record := range trainingData {

		// Skip the header.
		if i == 0 {
			continue
		}

		// Parse the Sales.
		yVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		for j := range xv {
			switch xv[j] {
			case "tv":
				tv = parseRecord(record, 0)
			case "radio":
				radio = parseRecord(record, 1)
			case "newspaper":
				newspaper = parseRecord(record, 2)
			}
		}

		xtrain := make([]float64, len(xv), len(xv))
		for i := range xv {
			switch xv[i] {
			case "tv":
				xtrain[i] = tv
			case "radio":
				xtrain[i] = radio
			case "newspaper":
				xtrain[i] = newspaper
			}
		}

		// Add these points to the regression value.
		r.Train(regression.DataPoint(yVal, xtrain))
	}

	// Train/fit the regression model.
	r.Run()

	// Output the trained model parameters.
	fmt.Printf("\nRegression Formula:\n%v\n", r.Formula)
	// fmt.Printf("Coefficients:\n%v\n", r.Coeff(3))

	// Open the test dataset file.
	f, err = os.Open("../data/test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a CSV reader reading from the opened file.
	reader = csv.NewReader(f)

	// Read in all of the CSV records
	reader.FieldsPerRecord = 4
	testData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Loop over the test data predicting y and evaluating the prediction
	// with the mean absolute error.
	var mAE, mSE float64
	for i, record := range testData {

		// Skip the header.
		if i == 0 {
			continue
		}

		// Parse the Sales.
		yObserved, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		for j := range xv {
			switch xv[j] {
			case "tv":
				tv = parseRecord(record, 0)
			case "radio":
				radio = parseRecord(record, 1)
			case "newspaper":
				newspaper = parseRecord(record, 2)
			}
		}

		xtest := make([]float64, len(xv), len(xv))
		for i := range xv {
			switch xv[i] {
			case "tv":
				xtest[i] = tv
			case "radio":
				xtest[i] = radio
			case "newspaper":
				xtest[i] = newspaper
			}
		}

		// Predict y with our trained model.
		yPredicted, err := r.Predict(xtest)

		// Add the to the mean absolute error.
		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
		mSE += math.Pow(yObserved-yPredicted, 2) / float64(len(testData))
	}

	// Output the MAE to standard out.
	fmt.Printf("MAE = %0.2f\n", mAE)
	fmt.Printf("MSE = %0.2f\n\n", mSE)

	// Clear xv
	xv = nil
}

func parseRecord(record []string, i int) float64 {
	x, err := strconv.ParseFloat(record[i], 64)
	if err != nil {
		log.Fatal(err)
	}

	return x
}
