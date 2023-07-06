package pkg

import (
	git "github.com/go-git/go-git/v5"
)

var repo *git.Repository

func OpenRepositoryAtRoot() (repo *git.Repository, err error) {
    if repo != nil {
        return repo, err
    }
    repo, err = git.PlainOpen(".")
    return repo, err
}
