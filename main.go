package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Profile struct {
	Dotfiles struct {
		Root  string
		Links map[string]string
	}
}

func main() {
	log.SetFlags(0)

	path := "/home/myrtle/projects/linum/profile.toml"

	if len(os.Args) < 2 {
		fmt.Println("usage: linum <command>")
		return
	}

	command := os.Args[1]

	profile, err := loadProfile(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := validateProfile(profile); err != nil {
		log.Fatal(err)
	}

	switch command {
	case "check":
		fmt.Println("checking profile")
	case "sync":
		fmt.Println("syncing profile")
		if err := syncConfigs(profile); err != nil {
			log.Fatal(err)
		}
	case "install":
		fmt.Println("installing packages")
		if err := installPackages(profile); err != nil {
			log.Fatal(err)
		}
	case "apply":
		fmt.Println("applying profile")
		if err := syncConfigs(profile); err != nil {
			log.Fatal(err)
		}

		if err := installPackages(profile); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown command: %s", command)
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

func syncConfigs(profile Profile) error {
	for name, path := range profile.Dotfiles.Links {
		root := os.ExpandEnv(profile.Dotfiles.Root)
		source := filepath.Join(root, name)
		target := os.ExpandEnv(path)

		if err := syncConfig(source, target); err != nil {
			return err
		}
	}

	return nil
}

func syncConfig(source, target string) error {
	// source(~/.dotfiles) chceck
	sourceInfo, err := os.Stat(source)
	sourceExists := !os.IsNotExist(err)

	if err != nil && sourceExists {
		return fmt.Errorf("cannot access source %s: %w", source, err)
	}

	if sourceExists && !sourceInfo.IsDir() && !sourceInfo.Mode().IsRegular() {
		return fmt.Errorf("source %s is not directory or a file", source)
	}

	//target(~/.config/*) check
	targetInfo, err := os.Lstat(target)
	targetExists := !os.IsNotExist(err)

	if err != nil && targetExists {
		return fmt.Errorf("cannot access target %s: %w", target, err)
	}

	if targetExists && targetInfo.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(target)
		if err != nil {
			return err
		}

		if link != source {
			return fmt.Errorf("target %s points to %s instead of %s", target, link, source)
		}

		if !sourceExists {
			return fmt.Errorf("source %s does not exist", source)
		}
		return nil
	}

	// source and target do not exist
	if !sourceExists && !targetExists {
		return fmt.Errorf("source %s and target %s do not exist", source, target)
	}

	// source exist, but target is not
	if sourceExists && !targetExists {
		err := os.Symlink(source, target)
		if err != nil {
			return err
		}
		return nil
	}

	if !sourceExists && targetExists && (targetInfo.IsDir() || targetInfo.Mode().IsRegular()) {
		err := os.Rename(target, source)
		if err != nil {
			return err
		}

		err = os.Symlink(source, target)
		if err != nil {
			return err
		}
	}

	if sourceExists && targetExists && targetInfo.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("conflict: source %s and target %s both exist", source, target)
	}

	return nil
}

func installPackages(profile Profile) error {
	return nil
}
