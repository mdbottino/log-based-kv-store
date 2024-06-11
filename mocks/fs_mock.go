package mocks

import (
	"errors"
	"io/fs"
	"time"

	"github.com/mdbottino/log-based-kv-store/filesystem"
)

type MockState struct {
	data         []byte
	length       int
	read         bool
	write_offset int
}

type MockFile struct {
	name  string
	state *MockState
}

type MockFileInfo struct {
	name string
	size int64
}

func (mfi MockFileInfo) Name() string {
	return mfi.name
}

func (mfi MockFileInfo) Size() int64 {
	return mfi.size
}

func (mfi MockFileInfo) Mode() fs.FileMode {
	return 600
}

func (mfi MockFileInfo) ModTime() time.Time {
	return time.Now()
}

func (mfi MockFileInfo) IsDir() bool {
	return false
}

func (mfi MockFileInfo) Sys() any {
	return nil
}

const MAX_SIZE int = 1024

func NewMockFile() MockFile {
	rows := make([]byte, MAX_SIZE)
	return MockFile{
		"some-file.log",
		&MockState{rows, 0, false, 0},
	}
}

func (mf MockFile) Seek(offset int64, whence int) (int64, error) {
	mf.state.read = false

	return 1, nil
}

func (mf MockFile) Read(p []byte) (n int, err error) {
	if mf.state.read {
		return 0, errors.New("EOF reached")
	}

	copied := copy(p, mf.state.data[:mf.state.length])
	mf.state.read = true

	return copied, nil
}

func (mf MockFile) Write(p []byte) (n int, err error) {
	if len(p)+mf.state.write_offset >= MAX_SIZE {
		return 0, errors.New("Write limit reached")
	}
	written := copy(mf.state.data[mf.state.write_offset:], p)
	mf.state.write_offset += written
	mf.state.length += written

	return written, nil
}

func (mf MockFile) GetSize() int {
	return mf.state.length
}

func (mf MockFile) Stat() (fs.FileInfo, error) {
	return MockFileInfo{mf.name, int64(mf.state.length)}, nil
}

type MockFileSystem struct {
	Handle *MockFile
}

func NewMockFileSystem() MockFileSystem {
	mf := NewMockFile()
	mfs := MockFileSystem{}

	mfs.Handle = &mf
	return mfs
}

func (mfs MockFileSystem) Create(name string) (filesystem.FileLike, error) {
	return mfs.Handle, nil
}
