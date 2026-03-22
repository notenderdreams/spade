package db

type Script struct {
	Name    string
	Command string
	Args    []string
	Runner  string
}

type Chain struct {
	ID    int
	Name  string
	Steps []ChainStep
}

type ChainStep struct {
	Seq    int
	Script Script
}

type ExportFile struct {
	Version string   `json:"version"`
	Scripts []Script `json:"scripts"`
}
