package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Profile struct {
	Dotfiles struct {
		Root  string
		Links map[string]string
	}
}

func loadProfile(path string) (Profile, error) {
	var profile Profile

	_, err := toml.DecodeFile(path, &profile)
	if err != nil {
		return Profile{}, fmt.Errorf("cannot load profile %s: %w", path, err)
	}
	return profile, nil
}

func validateProfile(profile Profile) error {
	if profile.Dotfiles.Root == "" {
		return fmt.Errorf("dotfiles.root cannot be empty")
	}

	if len(profile.Dotfiles.Links) == 0 {
		return fmt.Errorf("dotfiles.links cannot be empty")
	}

	for name, path := range profile.Dotfiles.Links {
		if path == "" {
			return fmt.Errorf("dotfiles.links.%s cannot be empty", name)
		}
	}
	return nil
}
