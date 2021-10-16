package main

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"go.uber.org/zap"
)

func logger() *zap.Logger {
	log, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to setup production logger %v", err))
	}

	return log
}

func main() {
	var (
		log = logger()
		app = NewGitFetcher(log)
	)

	repo, err := app.FetchRepo("https://github.com/go-git/go-billy")
	if err != nil {
		log.With(zap.Error(err)).Fatal("failed to fetch git repository")
	}

	fileChanges, err := app.FileChanges(repo)
	if err != nil {
		log.With(zap.Error(err)).Fatal("failed to get changeCharts")
	}

	filesChart := chartsFromFileChanges(fileChanges)

	renderCharts(filesChart)
}

func renderCharts(chs []Chart) {
	page := components.NewPage()
	for _, c := range chs {
		line := charts.NewLine()
		// set some global options like Title/Legend/ToolTip or anything else
		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
			charts.WithTitleOpts(opts.Title{
				Title: fmt.Sprintf("File: %s", c.Filename),
			}))

		var (
			xAxis     []string
			lineItems []opts.LineData
		)
		for _, point := range c.TimelinePoints {
			xAxis = append(xAxis, point.Key)
			lineItems = append(lineItems, opts.LineData{
				Value: point.AmountOfChanges,
			})
		}

		// Put data into instance
		line.SetXAxis(xAxis).
			AddSeries("changes", lineItems).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
		page.AddCharts(line)
	}

	// Where the magic happens
	f, _ := os.Create("charts.html")
	_ = page.Render(f)
	_ = f.Close()
}
