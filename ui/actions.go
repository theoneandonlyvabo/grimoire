package ui

import (
	"fmt"

	"github.com/theoneandonlyvabo/grimoire/core"
)

func RunForge() error {
	if !core.IsGitRepo() {
		return fmt.Errorf("Not a git repository, run 'git init' first")
	}

	metadata, err := core.GetMetadata()
	if err != nil {
		return fmt.Errorf("Failed to read git metadata: %w", err)
	}

	files, err := core.ScanFiles(".")
	if err != nil {
		return fmt.Errorf("Failed to scan project files: %w", err)
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
		return fmt.Errorf("Failed to forge grimoire: %w", err)
	}

	fmt.Println("Grimoire forged.")
	return nil
}

func RunCarve() error {
	grimoire, err := core.Load()
	if err != nil {
		return fmt.Errorf("No grimoire found, run 'grimoire forge' first")
	}

	if err := refreshMetadata(grimoire); err == nil {
		core.Save(grimoire)
	}

	return Start(grimoire, false)
}

func RunCast() error {
	grimoire, err := core.Load()
	if err != nil {
		return fmt.Errorf("No grimoire found, run 'grimoire forge' first")
	}

	if err := refreshMetadata(grimoire); err == nil {
		core.Save(grimoire)
	}

	return Start(grimoire, true)
}

func refreshMetadata(grimoire *core.Grimoire) error {
	meta, err := core.GetMetadata()
	if err != nil {
		return err
	}
	grimoire.Meta = meta
	return nil
}

func runFromMenu(command string) error {
	switch command {
	case "Forge":
		return RunForge()
	case "Carve":
		return RunCarve()
	case "Cast":
		return RunCast()
	}
	return nil
}
