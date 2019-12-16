// +build !sendonly

package app

import "github.com/buchanae/ink/win"

func DefaultConfig() Config {
	return Config{
		Snapshot: SnapshotConfig{
			Width:  2400,
			Height: 2400,
			Dir:    "snapshots",
		},
		Window: win.Config{
			Name:   "zinc",
			X:      10,
			Y:      10,
			Width:  800,
			Height: 800,
		},
	}
}

type Config struct {
	Snapshot SnapshotConfig
	Window   win.Config
}

type SnapshotConfig struct {
	Width, Height int
	Dir           string
}
