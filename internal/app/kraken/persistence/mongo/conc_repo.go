package mongo

import (
	"github.com/globalsign/mgo"
)

// Concurrent repo operating with mgo.Session copy
// to make real concurrency rock'n'roll
type ConcurrentRepository struct {
	mongo  *Mongo
	dbName string
	cName  string
}

// Execute statement on collection in regular mode
func (repo *ConcurrentRepository) C(exec func(c *mgo.Collection)) error {
	s, err := repo.mongo.GetSession(false)
	if err != nil {
		return err
	}

	db := s.DB(repo.dbName)
	c := db.C(repo.cName)

	exec(c)

	return nil
}

// Execute statement on collection in concurrent mode
func (repo *ConcurrentRepository) CC(exec func(c *mgo.Collection)) error {
	s, err := repo.mongo.GetSession(false)
	if err != nil {
		return err
	}

	db := s.DB(repo.dbName)
	c := db.C(repo.cName)

	exec(c)

	return nil
}
