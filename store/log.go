package store

import (
	"github.com/mdbottino/log-based-kv-store/filesystem"
)

type Log struct {
	segments []Segment
}

func NewLog(folder string, fs filesystem.FileCreator) Log {
	segment := NewSegment(folder, fs)

	return Log{
		[]Segment{segment},
	}
}

func (l Log) Append(key, value string) error {
	idx := len(l.segments) - 1
	latestSegment := l.segments[idx]

	return latestSegment.Append(key, value)
}

func (l Log) Find(key string) (string, error) {
	idx := len(l.segments) - 1
	latestSegment := l.segments[idx]

	return latestSegment.Find(key)
}
