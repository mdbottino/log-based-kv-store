package filesystem

import (
	"io"
	"os"
)

type FileLike interface {
	Seek(int64, int) (int64, error)
	io.Reader
	io.Writer
}

type File struct {
	handle FileLike
}

type FileCreator interface {
	Create(string) (FileLike, error)
}

type FileSystem struct{}

func NewFileSystem() FileSystem {
	return FileSystem{}
}

func (fs FileSystem) Create(name string) (FileLike, error) {
	return os.Create(name)
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	return f.handle.Seek(offset, whence)
}
