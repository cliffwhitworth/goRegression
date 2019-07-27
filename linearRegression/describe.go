package main

import (
	"fmt"
)

func descriptiveStats() {

	df := getAllData()

	summary := df.Describe()

	fmt.Println(summary)
}
