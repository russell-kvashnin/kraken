package file

import (
	"github.com/dchest/uniuri"
	"path/filepath"
	"strings"
	"time"
)

// Config error codes
const (
	ErrorDomain = "FILE_DOMAIN"

	StoreErrorCode    = "FILE_STORE_ERROR"
	NotFoundErrorCode = "FILE_NOT_FOUND_ERROR"

	UploadedEventAlias = "file.uploaded"

	defaultAlphabet     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultStringLength = 6
)

// File document structure
type File struct {
	ShortUrl string `bson:"short_url" json:"shortUrl"`
	Path     string `bson:"path" json:"-"`
	Location []Node `bson:"location" json:"-"`

	UploadedAt time.Time `bson:"uploaded_at" json:"uploaded_at"`

	Meta Meta `bson:"meta" json:"meta"`
}

// File location node
type Node struct {
	NodeId string `bson:"node_id"`
}

// File meta info
type Meta struct {
	Tags             []string `bson:"tags" json:"tags,omitempty"`
	FileType         string   `bson:"file_type" json:"fileType"`
	FileSize         int64    `bson:"file_size" json:"fileSize"`
	OriginalFileName string   `bson:"original_file_name" json:"originalFileName"`
	OriginalFileExt  string   `bson:"original_file_ext" json:"originalFileExt"`
}

// Upload file
func (f *File) Upload(cmd UploadFileCommand) {
	originalExt := filepath.Ext(cmd.Header.Filename)
	originalName := strings.TrimSuffix(cmd.Header.Filename, originalExt)
	fileType := cmd.Header.Header.Get("Content-Type")

	f.ShortUrl = uniuri.NewLenChars(defaultStringLength, []byte(defaultAlphabet))
	f.Path = f.generatePath(cmd.UploadRoot, f.ShortUrl)
	f.UploadedAt = time.Now()
	f.Meta = Meta{
		Tags:             cmd.Tags,
		FileType:         fileType,
		FileSize:         cmd.Header.Size,
		OriginalFileName: originalName,
		OriginalFileExt:  originalExt,
	}
	f.Location = []Node{
		{
			NodeId: cmd.NodeId,
		},
	}
}

// Update file
func (f *File) Update(cmd UpdateFileCommand) {
	f.Path = cmd.Path
}

// Generate file path
func (f *File) generatePath(uploadRoot string, filename string) string {
	var path string

	pLen := len(filename) - 1

	for i, bytes := range filename {
		if i == pLen {
			break
		}

		path = uploadRoot + string(bytes) + "/"
	}

	return path
}

// File uploaded event
type UploadedEvent struct {
	NodeId   string
	ShortUrl string
}

// Uploaded event handler
type UploadedEventHandler interface {
	Handle(event UploadedEvent)
}

// File model repository interface
type Repository interface {
	Store(file File) error
	Get(shortUrl string) (File, error)
}
