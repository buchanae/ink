package win

type Event int

const (
	QuitEvent Event = iota
	// TODO app events belong in the app package
	RefreshEvent
	SnapshotEvent
	ReturnEvent
)
