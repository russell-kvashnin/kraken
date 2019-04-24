// This package provides file system abstraction via MMAP syscall
package mmap_fs

import (
	"github.com/edsrzf/mmap-go"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"os"
	"path"
)

const (
	errorDomain         = "MMAP_FS"
	createFileErrorCode = "CREATE_FILE_ERROR"
	resizeFileErrorCode = "RESIZE_FILE_ERROR"
	mapFileErrorCode    = "MAP_FILE_ERROR"
)

// Write file to fs via MMAP
func MakeWriteMap(blobSize int64, filePath string) (*mmap.MMap, kerr.Error) {
	var (
		f       *os.File
		m       *mmap.MMap
		err     error
		e       kerr.Error
		details map[string]string
	)

	details = make(map[string]string)
	details["filepath"] = filePath

	f, err = os.Create(path.Join(filePath))
	if err != nil {
		return nil, kerr.NewErr(kerr.ErrLvlError, errorDomain, createFileErrorCode, err, details)
	}
	defer f.Close()

	err = f.Truncate(blobSize)
	if err != nil {
		return nil, kerr.NewErr(kerr.ErrLvlError, errorDomain, resizeFileErrorCode, err, details)
	}

	m, e = makeMmap(f, mmap.RDWR)
	if e != nil {
		return nil, e
	}

	return m, nil
}

// Map file to memory
func makeMmap(f *os.File, prot int) (*mmap.MMap, kerr.Error) {
	var (
		err error
		m   mmap.MMap
	)

	switch prot {
	case mmap.RDONLY:
		break
	case mmap.RDWR:
		break
	default:
		prot = mmap.RDONLY
	}

	m, err = mmap.Map(f, prot, 0)
	if err != nil {
		return nil, kerr.NewErr(kerr.ErrLvlError, errorDomain, mapFileErrorCode, err, nil)
	}

	return &m, nil
}
