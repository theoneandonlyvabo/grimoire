package core

import (
	"encoding/json"
	"os"
	"time"
)

const GrimoireFile = ".grimoire"

type Grimoire struct {
	Version  string `json:"version"`
	Document []Doc  `json:"documents"`
}

type Doc struct {
	ID          string     `json:"id"`
	LinkedFile  string     `json:"linked_file"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Author      string     `json:"author"`
	UpdatedAt   string     `json:"updated_at"`
	Functions   []Function `json:"functions"`
}

type Function struct {
	Name      string `json:"name"`
	Signature string `json:"signature"`
	Notes     string `json:"notes"`
	Author    string `json:"author"`
	UpdatedAt string `json:"updated_at"`
}

func NewDoc(file string, author string) Doc {
	return Doc{
		ID:         file,
		LinkedFile: file,
		Status:     "wip",
		Author:     author,
		UpdatedAt:  time.Now().Format("2006-01-02 15:04"),
		Functions:  []Function{},
	}
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
