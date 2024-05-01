package mocks

import (
	"errors"

	"github.com/mdbottino/log-based-kv-store/filesystem"
)

type MockState struct {
	data         []byte
	length       int
	read         bool
	write_offset int
}

type MockFile struct {
	state *MockState
}

const MAX_SIZE int = 1024

func NewMockFile() MockFile {
	rows := make([]byte, MAX_SIZE)
	return MockFile{
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
