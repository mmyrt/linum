package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	path := os.ExpandEnv("$HOME/projects/linum/profile.toml")

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
		info("checking profile")
	case "sync":
		info("syncing profile")
		if err := syncConfigs(profile); err != nil {
			fatal(err)
		}
	case "install":
		info("installing packages")
		if err := installPackages(profile); err != nil {
			fatal(err)
		}
	case "apply":
		info("applying profile")
		if err := syncConfigs(profile); err != nil {
			fatal(err)
		}

		if err := installPackages(profile); err != nil {
			fatal(err)
		}
	default:
		fatal("Unknown command: ", command)
	}
}
