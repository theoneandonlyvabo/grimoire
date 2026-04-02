package ui

import (
	"fmt"

	"github.com/theoneandonlyvabo/grimoire/core"
)

func RunForge() error {
	if !core.IsGitRepo() {
		return fmt.Errorf("not a git repository, run 'git init' first")
	}

	metadata, err := core.GetMetadata()
	if err != nil {
		return fmt.Errorf("failed to read git metadata: %w", err)
	}

	files, err := core.ScanFiles(".")
	if err != nil {
		return fmt.Errorf("failed to scan project files: %w", err)
	}

	author := core.GetUserName()
	var documents []core.Doc
	for _, file := range files {
		doc := core.NewDoc(file, author)
		doc.Functions = core.ParseFunctions(file)
		documents = append(documents, doc)
	}

	grimoire := &core.Grimoire{
		Meta:     metadata,
		Document: documents,
	}

	if err := core.Save(grimoire); err != nil {
		return fmt.Errorf("failed to forge grimoire: %w", err)
	}

	fmt.Println("Grimoire forged.")
	return nil
}

func RunCarve() error {
	grimoire, err := core.Load()
	if err != nil {
		return fmt.Errorf("no grimoire found, run 'grimoire forge' first")
	}
	return Start(grimoire, false)
}

func RunCast() error {
	grimoire, err := core.Load()
	if err != nil {
		return fmt.Errorf("no grimoire found, run 'grimoire forge' first")
	}
	return Start(grimoire, true)
}
