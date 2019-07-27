package main

import (
	"fmt"
	"image/color"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func histogram() {

	df := getAllData()

	// Create a histogram for each of the columns in the dataset.
	for _, colName := range df.Names() {

		plotVals := make(plotter.Values, df.Nrow())
		for i, floatVal := range df.Col(colName).Float() {
			plotVals[i] = floatVal
		}

		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)

		h, err := plotter.NewHist(plotVals, 16)
		if err != nil {
			log.Fatal(err)
		}

		h.Normalize(1)

		p.Add(h)

		if err := p.Save(4*vg.Inch, 4*vg.Inch, "./plots/"+colName+"_hist.png"); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Histograms saved in plots folder")
}

func scatterplot() {

	df := getAllData()

	yVals := df.Col("Sales").Float()

	// Create a scatter plot for each of the features in the dataset.
	for _, colName := range df.Names() {

		pts := make(plotter.XYs, df.Nrow())

		for i, floatVal := range df.Col(colName).Float() {
			pts[i].X = floatVal
			pts[i].Y = yVals[i]
		}

		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.X.Label.Text = colName
		p.Y.Label.Text = "y"
		p.Add(plotter.NewGrid())

		s, err := plotter.NewScatter(pts)
		if err != nil {
			log.Fatal(err)
		}
		s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
		s.GlyphStyle.Radius = vg.Points(3)

		p.Add(s)
		if err := p.Save(4*vg.Inch, 4*vg.Inch, "./plots/"+colName+"_scatter.png"); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Scatterplots saved in plots folder")
}
