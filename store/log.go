package store

import (
	"github.com/mdbottino/log-based-kv-store/filesystem"
)

type Log struct {
	segments []Segment
	folder   string
	fs       filesystem.FileCreator
}

func NewLog(folder string, fs filesystem.FileCreator) Log {
	segment := NewSegment(folder, fs)

	return Log{
		[]Segment{segment},
		folder,
		fs,
	}
}

func (l Log) Append(key, value string) error {
	latestSegment := l.GetLatestSegment()

	return latestSegment.Append(key, value)
}

func (l Log) Find(key string) (string, error) {
	length := len(l.segments)
	for i := length - 1; i > 0; i-- {
		segment := l.segments[i]
		val, err := segment.Find(key)
		if err != nil {
			continue
		}

		// We found the key
		return val, err
	}

	// We try at last with the oldest segment
	return l.segments[0].Find(key)
}

func (l Log) GetLatestSegment() Segment {
	idx := len(l.segments) - 1
	return l.segments[idx]
}

func (l *Log) AddSegment() {
	segment := NewSegment(l.folder, l.fs)
	l.segments = append(l.segments, segment)
}
