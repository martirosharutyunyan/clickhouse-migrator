package types

type Migration struct {
	ShardNumber int
	Version     int64
	Source      string
	IsApplied   bool
}
