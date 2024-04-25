package store

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type Log struct {
	filename string
	handle   *os.File
}

func NewLog(folder string) Log {
	timestamp := time.Now().Unix()
	filename := path.Join(folder, fmt.Sprintf("%d.log", timestamp))

	handle, err := os.Create(filename)
	if err != nil {
		panic("failed to create the log file")
	}

	return Log{
		filename,
		handle,
	}

}

func (l Log) Append(key, value string) error {
	_, err := l.handle.Seek(0, 2)
	if err != nil {
		return errors.New("failed to move the offset to end of the file")
	}

	_, err = l.handle.Write([]byte(fmt.Sprintf("%s: %s\n", key, value)))
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

func (l Log) Find(key string) (string, error) {
	_, err := l.handle.Seek(0, 0)
	if err != nil {
		return "", errors.New("failed to move the offset to beginning of the file")
	}

	scanner := bufio.NewScanner(l.handle)

	for scanner.Scan() {
		line := scanner.Text()
		k, v, err := getKeyAndValue(line)

		if err == nil && k == key {
			return v, nil
		}
	}

	return "", errors.New("couldn't find the key in the file")
}
