package main

import (
	"fmt"
	"os"
	"path/filepath"
)

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

		info("unchanged %s", target)
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

		success("linked %s -> %s", target, source)
		return nil
	}

	if !sourceExists && targetExists && (targetInfo.IsDir() || targetInfo.Mode().IsRegular()) {
		err := os.MkdirAll(filepath.Dir(source), 0755)
		if err != nil {
			return err
		}

		err = os.Rename(target, source)
		if err != nil {
			return fmt.Errorf("cannot import %s: %w", target, err)
		}

		if err := os.Symlink(source, target); err != nil {
			if rollbackErr := os.Rename(source, target); rollbackErr != nil {
				return fmt.Errorf("cannot link %s: %w; rollback failed: %v", target, err, rollbackErr)
			}
		}

		success("imported: %s", target)
		success("linked %s -> %s", target, source)

		return nil
	}

	if sourceExists && targetExists && targetInfo.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("conflict: source %s and target %s both exist", source, target)
	}

	return nil
}
