package db

type Script struct {
	Name    string
	Command string
	Args    []string
}

type ExportFile struct {
	Version string      `json:"version"`
	Scripts []Script `json:"scripts"`
}