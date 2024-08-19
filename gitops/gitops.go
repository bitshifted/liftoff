// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package gitops

import (
	"os"

	"github.com/bitshifted/liftoff/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitHandler struct {
	URL         string
	Version     string
	Destination string
}

func (gh *GitHandler) Fetch() error {
	log.Logger.Info().Msgf("Cloning Git repository %s", gh.URL)
	repo, err := git.PlainClone(gh.Destination, false, &git.CloneOptions{
		URL:      gh.URL,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Failed to clone git repository %s", gh.URL)
		return nil
	}
	if gh.Version == "" {
		log.Logger.Info().Msg("Version is not specified, defaulting to main branch")
		return nil
	}
	log.Logger.Debug().Msgf("Looking up tag %s", gh.Version)
	commitHash, err := gh.getCommitHashForTagName(repo)
	if err != nil {
		return err
	}
	if commitHash == "" {
		commitHash = gh.Version
	}
	wt, err := repo.Worktree()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get repository work tree")
		return err
	}
	err = wt.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
	})
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Failed to checkout commit hash %s", commitHash)
	} else {
		log.Logger.Info().Msgf("Checked out repository version %s", commitHash)
	}
	return err
}

func (gh *GitHandler) getCommitHashForTagName(repo *git.Repository) (string, error) {
	iter, err := repo.Tags()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get tags")
		return "", err
	}
	commitHash := ""
	ierr := iter.ForEach(func(ref *plumbing.Reference) error {
		to, err := repo.TagObject(ref.Hash())
		if err != nil {
			log.Logger.Error().Err(err).Msgf("Failed to get tag object for hash %s", ref.Hash().String())
		}
		if gh.Version == to.Name {
			log.Logger.Debug().Msgf("Found tag with name %s. Commit hash: %s", gh.Version, to.Hash.String())
			commitHash = to.Hash.String()
		}
		return nil
	})
	if ierr != nil {
		return "", ierr
	}

	return commitHash, nil
}
