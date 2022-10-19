package main

import (
	"fmt"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/version"
)

// main loop
// 1. take latest set of "known card ids"
// 2. crawl latest cards until we intersect with "latest known"
// 3. push jobs for downloading corresponding images
// 4. update "latest known card ids"
// 5. notify pipeline

func main() {
	//	var knownCards []CardID

	// TODO: load from disk
	//	knownCards = make([]CardID, 0) // empty for now
	fmt.Printf("Starting up...\nVersion: %s\nGit commit: %.6s\n", version.AppVersion, version.GitCommit)
}
