package utils

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// This function is thought over the fact that the most recent commit of a git repository is one of its leafs. The most recent before that is either one of the parents of the most recent one or one of the leafs of the repository, and so on. It will run the callback for all commits posterior to "since", in historical order according to their creation time.
func LogTimeline(repo *git.Repository, since *time.Time, cb func (c *object.Commit) error) error {
    // 1. Define the set of leafs
    var set = map[plumbing.Hash]*object.Commit{}
    branches, _ := repo.References()
    for r, err := branches.Next(); err == nil; r, err = branches.Next() {
        if r.Hash().IsZero() {continue}
        c, c_err := repo.CommitObject(r.Hash())
        if c_err != nil {continue}
        set[c.Hash] = c
    }

    var err error

    // 2. Find the most recent commit by iterating over the set.
    step_two:
    var most_recent *object.Commit
    for _, commit := range set {
        if most_recent == nil {
            most_recent = commit
            continue
        }
        if commit.Committer.When.After(most_recent.Committer.When) {
            most_recent = commit
        }
    }

    // 3. If the commit found in 2. has no parents or its creation time is posterior to "since", return. 
    if most_recent.NumParents() == 0 || (since != nil && since.After(most_recent.Committer.When)) {
        return nil
    }

    // 4. Run the callback for the most recent commit.
    err = cb(most_recent)
    if err != nil {return err}

    // 5. Redefine the set of 1. (remove the commit found in 2., add its parents) and go to step 2.
    delete(set, most_recent.Hash)
    parents_iter := most_recent.Parents()
    for c, err := parents_iter.Next(); err == nil; c, err = parents_iter.Next() {
        set[c.Hash] = c
    }
    goto step_two
}
