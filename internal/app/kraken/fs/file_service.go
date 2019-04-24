package fs

import (
	"github.com/edsrzf/mmap-go"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"github.com/russell-kvashnin/kraken/internal/pkg/mmap_fs"
	"io"
	"os"
	"path"
)

// File system service error codes
const (
	ErrorDomain = "FILE_SYSTEM"

	OpenErrorCode      = "FS_OPEN_FILE_ERROR"
	CloseErrorCode     = "FS_CLOSE_FILE_ERROR"
	CreateErrorCode    = "FS_CREATE_FILE_ERROR"
	CreateDirErrorCode = "FS_CREATE_DIR_ERROR"
	TruncateErrorCode  = "FS_TRUNCATE_FILE_ERROR"
	MmapErrorCode      = "FS_MMAP_ERROR_CODE"

	uploadDirPerms = 0755
)

// File service
type FileService struct {
	cfg config.FSConfig
}

// Constructor
func NewFileService(cfg config.FSConfig) *FileService {
	svc := new(FileService)
	svc.cfg = cfg

	return svc
}

// Save file to FS via mmap
func (svc *FileService) Save(file []byte, fSize int64, fileName string) (string, error) {
	var (
		err       error
		uploadDir string
		filePath  string
		details   map[string]string
	)
	details = make(map[string]string)

	uploadDir, err = svc.mkUploadedFileDir(fileName)
	if err != nil {
		return "", kerr.NewErr(kerr.ErrLvlError, ErrorDomain, CreateDirErrorCode, err, details)
	}

	fileName = fileName[len(fileName)-1:]
	filePath = path.Join(uploadDir, fileName)
	m, e := mmap_fs.MakeWriteMap(fSize, filePath)
	if e != nil {
		return "", e
	}

	copy(*m, file)

	return filePath, nil
}

// Read file via mmap
func (svc *FileService) Read(filepath string) (*mmap.MMap, error) {
	f, err := os.Open(filepath)
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, OpenErrorCode, err, nil)
		return nil, e
	}

	m, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, MmapErrorCode, err, nil)
		return nil, e
	}

	if err := f.Close(); err != nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, CloseErrorCode, err, nil)

		return nil, e
	}

	return &m, nil
}

// Create directory
func (svc *FileService) Mkdir(path string, perms os.FileMode) error {
	err := os.Mkdir(path, perms)
	if err != nil {
		details := make(map[string]string)
		details["dir"] = path
		details["perms"] = string(perms)

		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, CreateDirErrorCode, err, details)

		return e
	}

	return nil
}

// Create uploaded file directory
func (svc *FileService) mkUploadedFileDir(filename string) (string, error) {
	uploadDir := svc.cfg.UploadDir
	pLen := len(filename) - 1

	for i, bytes := range filename {
		if i == pLen {
			break
		}

		uploadDir += string(bytes) + "/"
	}

	err := os.MkdirAll(uploadDir, uploadDirPerms)
	if err != nil {
		return "", err
	}

	return uploadDir, nil
}

// Prepare file for mmap write
func (svc *FileService) prepareFile(uploadDir string, fileName string, fileSize int64) (*os.File, error) {
	fileRune := string(fileName[len(fileName)-1])

	f, err := os.Create(path.Join(uploadDir, fileRune))
	if err != nil {
		details := make(map[string]string)
		details["path"] = uploadDir + fileRune

		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, CreateErrorCode, err, details)

		return nil, e
	}

	err = f.Truncate(fileSize)
	if err != nil {
		details := make(map[string]string)
		details["size"] = string(fileSize) + " bytes."

		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, TruncateErrorCode, err, details)

		return nil, e
	}

	return f, nil
}

// Save file via mmap syscall
func (svc *FileService) saveFile(file *os.File, mFile io.Reader) error {
	m, err := mmap.Map(file, mmap.RDWR, 0)
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, MmapErrorCode, err, nil)
		return e
	}

	_, err = mFile.Read(m)
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, MmapErrorCode, err, nil)
		return e
	}

	err = m.Flush()
	if err != nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, MmapErrorCode, err, nil)
		return e
	}

	if err := m.Unmap(); err != nil {
		e := kerr.NewErr(kerr.ErrLvlError, ErrorDomain, MmapErrorCode, err, nil)

		return e
	}

	return nil
}
