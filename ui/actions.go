package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/theoneandonlyvabo/grimoire/core"
)

func RunForge() error {
	if !core.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	if _, err := core.Load(); err == nil {
		return fmt.Errorf("grimoire already exists, use 'grimoire carve' to edit")
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
		Version:  "1.0.0",
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

	syncFiles(grimoire)
	core.Save(grimoire)

	return runTUI(grimoire, false)
}

func RunCast() error {
	grimoire, err := core.Load()
	if err != nil {
		return fmt.Errorf("no grimoire found, run 'grimoire forge' first")
	}

	syncFiles(grimoire)
	core.Save(grimoire)

	return runTUI(grimoire, true)
}

func runTUI(grimoire *core.Grimoire, readOnly bool) error {
	m := NewModel(grimoire, readOnly)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func syncFiles(grimoire *core.Grimoire) {
	files, err := core.ScanFiles(".")
	if err != nil {
		return
	}

	existing := map[string]core.Doc{}
	for _, doc := range grimoire.Document {
		existing[doc.LinkedFile] = doc
	}

	author := core.GetUserName()
	var updated []core.Doc
	for _, file := range files {
		if doc, found := existing[file]; found {
			updated = append(updated, doc)
		} else {
			doc := core.NewDoc(file, author)
			doc.Functions = core.ParseFunctions(file)
			updated = append(updated, doc)
		}
	}

	grimoire.Document = updated
}
