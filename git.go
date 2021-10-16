package main

import (
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type GitFetcher struct {
	log *zap.Logger
}

func NewGitFetcher(log *zap.Logger) *GitFetcher {
	return &GitFetcher{log: log}
}

func (a GitFetcher) FetchRepo(url string) (*git.Repository, error) {
	return git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})
}

func (a GitFetcher) FileChanges(repo *git.Repository) (FileChanges, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, errors.Wrap(err, "get head")
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, errors.Wrap(err, "retrieve commit history")
	}

	points := make(FileChanges)

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
