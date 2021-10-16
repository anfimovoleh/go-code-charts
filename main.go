package main

import (
	"fmt"

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

	charts := chartsFromFileChanges(fileChanges)
	fmt.Println(charts)
}
