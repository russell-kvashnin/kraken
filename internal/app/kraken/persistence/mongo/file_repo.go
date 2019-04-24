package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
)

// Mongo collection for files
const FileCollectionName = "files"

// File document repository
type FileRepository struct {
	ConcurrentRepository
}

// Constructor
func NewFileRepository(mongo *Mongo) *FileRepository {
	repo := new(FileRepository)
	repo.mongo = mongo
	repo.dbName = mongo.GetCurrentDBName()
	repo.cName = FileCollectionName

	return repo
}

// Store file info in mongo
func (repo *FileRepository) Store(model file.File) error {
	var (
		e error
	)

	err := repo.CC(func(c *mgo.Collection) {
		err := c.Insert(model)

		if err != nil {
			e = kerr.NewErr(kerr.ErrLvlError, file.ErrorDomain, file.StoreErrorCode, err, nil)
		}
	})

	if err != nil {
		return err
	}

	return e
}

// Get file from database
func (repo *FileRepository) Get(shortUrl string) (file.File, error) {
	var (
		doc file.File
		e   error
	)

	err := repo.CC(func(c *mgo.Collection) {
		model := new(file.File)

		cErr := c.Find(bson.M{
			"short_url": shortUrl,
		}).One(model)

		if cErr != nil {
			details := make(map[string]string)
			details["short_url"] = shortUrl

			e = kerr.NewErr(kerr.ErrLvlError, file.ErrorDomain, file.NotFoundErrorCode, cErr, details)
		}

		doc = *model
	})

	if err != nil {
		return file.File{}, err
	}

	return doc, e
}
