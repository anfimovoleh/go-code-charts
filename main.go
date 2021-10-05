package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"

	"go.uber.org/zap"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Points map[string][]time.Time

type App struct {
	log *zap.Logger
}

func NewApp(log *zap.Logger) *App {
	return &App{log: log}
}

func (a App) FetchRepo(url string) (*git.Repository, error) {
	return git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})
}

func (a App) ChangeCharts(repo *git.Repository) (Points, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, errors.Wrap(err, "get head")
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, errors.Wrap(err, "retrieve commit history")
	}

	points := make(Points)

	err = cIter.ForEach(func(c *object.Commit) error {
		files, err := c.Files()
		if err != nil {
			return errors.Wrapf(err, "get files from commit %s", c.Hash)
		}

		err = files.ForEach(func(file *object.File) error {
			points[file.Name] = append(points[file.Name], c.Author.When)
			return nil
		})
		if err != nil {
			return errors.Wrapf(err, "iterate over files in commit %s", c.Hash)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "process commit history")
	}

	// sort timestamps for each file
	for file := range points {
		sort.Slice(points[file], func(i, j int) bool {
			return points[file][i].Before(points[file][j])
		})
	}

	return points, nil
}

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
		app = NewApp(log)
	)

	repo, err := app.FetchRepo("https://github.com/go-git/go-billy")
	if err != nil {
		log.With(zap.Error(err)).Fatal("failed to fetch git repository")
	}

	points, err := app.ChangeCharts(repo)
	if err != nil {
		log.With(zap.Error(err)).Fatal("failed to get changeCharts")
	}

	for key, vals := range points {
		fmt.Println("file:", key)
		fmt.Println("values:", vals)
		fmt.Println("________________________________________________________________________")
	}

	//log.With(zap.Any("points", points)).Info("result")
}
