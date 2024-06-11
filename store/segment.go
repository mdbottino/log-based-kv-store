package store

import (
	"bufio"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/mdbottino/log-based-kv-store/filesystem"
)

type Segment struct {
	filename string
	handle   filesystem.FileLike
}

func NewSegment(folder string, fs filesystem.FileCreator) Segment {
	timestamp := time.Now().Unix()
	filename := path.Join(folder, fmt.Sprintf("%d.log", timestamp))

	handle, err := fs.Create(filename)
	if err != nil {
		panic("failed to create the log file")
	}

	return Segment{
		filename,
		handle,
	}
}

func (s Segment) Append(key, value string) error {
	_, err := s.handle.Seek(0, 2)
	if err != nil {
		return errors.New("failed to move the offset to end of the file")
	}

	_, err = s.handle.Write([]byte(fmt.Sprintf("%s:%s\n", key, value)))
	if err != nil {
		return errors.New("failed to write to the log")
	}

	return nil
}

func getKeyAndValue(line string) (string, string, error) {
	parts := strings.Split(line, ":")

	if len(parts) == 2 {
		return parts[0], parts[1], nil
	}

	return "", "", errors.New("failed to retrieve key and value from line")
}

func (s Segment) Find(key string) (string, error) {
	_, err := s.handle.Seek(0, 0)
	if err != nil {
		return "", errors.New("failed to move the offset to beginning of the file")
	}

	scanner := bufio.NewScanner(s.handle)

	value := ""
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		k, v, err := getKeyAndValue(line)

		if err == nil && k == key {
			found = true
			value = v
		}
	}

	if found {
		return value, nil
	}

	return "", errors.New("couldn't find the key in the file")
}

func (s Segment) Size() (int, error) {
	info, err := s.handle.Stat()
	if err != nil {
		return 0, errors.New("couldn't obtain stat from file")

	}

	return int(info.Size()), nil
}
