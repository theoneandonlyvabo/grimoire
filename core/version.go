package core

import (
	_ "embed"
	"encoding/json"
)

//go:embed version.json
var versionData []byte

type VersionInfo struct {
	Release   string `json:"release"`
	Timestamp string `json:"timestamp"`
	Changelog string `json:"changelog"`
}

func GetVersion() VersionInfo {
	var v VersionInfo
	json.Unmarshal(versionData, &v)
	return v
}
