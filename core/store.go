package core

import (
	"encoding/json"
	"os"
)

const GrimoireFile = ".grim"

type Grimoire struct {
	Meta     MetaData `json:"meta"`
	Document []Doc    `json:"documents"`
}

type MetaData struct {
	Version           string   `json:"version"`
	Repository        string   `json:"repository"`
	Branch            string   `json:"branch"`
	Commits           int      `json:"commits"`
	LastCommit        string   `json:"last_commit"`
	LastCommitMessage string   `json:"last_commit_message"`
	LastCommitDate    string   `json:"last_commit_date"`
	Contributors      []string `json:"contributors"`
}

type Doc struct {
	ID          string `json:"id"`
	LinkedFile  string `json:"linked_file"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
	Status      string `json:"status"`
}

func Load() (*Grimoire, error) {
	data, err := os.ReadFile(GrimoireFile)
	if err != nil {
		return nil, err
	}
	var grim Grimoire
	err = json.Unmarshal(data, &grim)
	if err != nil {
		return nil, err
	}
	return &grim, nil
}

func Save(grim *Grimoire) error {
	data, err := json.MarshalIndent(grim, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(GrimoireFile, data, 0644)
}
