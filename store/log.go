package store

import (
	"errors"

	"github.com/mdbottino/log-based-kv-store/filesystem"
)

// 1KB
const DEFAULT_MAX_SEGMENT_SIZE = 1024

type Log struct {
	segments       []Segment
	folder         string
	fs             filesystem.FileCreator
	maxSegmentSize int
}

func NewLog(folder string, fs filesystem.FileCreator, cfg Config) Log {
	segment := NewSegment(folder, fs)

	size := DEFAULT_MAX_SEGMENT_SIZE
	if cfg.Segments.MaxSize != 0 {
		size = cfg.Segments.MaxSize
	}

	return Log{
		[]Segment{segment},
		folder,
		fs,
		size,
	}
}

func (l *Log) Append(key, value string) error {
	latestSegment := l.GetLatestSegment()

	size, err := latestSegment.Size()
	if err != nil {
		return errors.New("an error occurred while checking for the log stats")
	}

	if size > l.maxSegmentSize {
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
