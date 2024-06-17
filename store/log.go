package store

import (
	"errors"

	"github.com/mdbottino/log-based-kv-store/filesystem"
)

type Log struct {
	segments []Segment
	folder   string
	fs       filesystem.FileCreator
}

// This is a pretty low value to deliverately trigger segment rotation more often
const MAX_SEGMENT_SIZE = 100

func NewLog(folder string, fs filesystem.FileCreator) Log {
	segment := NewSegment(folder, fs)

	return Log{
		[]Segment{segment},
		folder,
		fs,
	}
}

func (l *Log) Append(key, value string) error {
	latestSegment := l.GetLatestSegment()

	size, err := latestSegment.Size()
	if err != nil {
		return errors.New("an error occurred while checking for the log stats")
	}

	if size > MAX_SEGMENT_SIZE {
		newSegment := l.AddSegment()
		return newSegment.Append(key, value)
	}

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

func (l *Log) AddSegment() Segment {
	segment := NewSegment(l.folder, l.fs)
	l.segments = append(l.segments, segment)

	return segment
}
