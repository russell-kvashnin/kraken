package mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"time"
)

// Mongo error codes
const (
	ErrorDomain = "MONGO_DB"

	DialErrCode  = "MONGO_DIAL_ERROR"
	NotConnected = "MONGO_NOT_CONNECTED"
)

// MongoDB service
type Mongo struct {
	cfg       config.MongoConfig
	session   *mgo.Session
	connLimit int
}

// Mongo connection channel
type ConnectInfo struct {
	Error error
	Done  bool
}

// Mongo service constructor
func NewMongo(cfg config.MongoConfig) *Mongo {
	mongo := new(Mongo)
	mongo.cfg = cfg

	return mongo
}

func (mongo *Mongo) GetCurrentDBName() string {
	return mongo.cfg.DBName
}

// Obtain mongodb connection
func (mongo *Mongo) Connect(c chan ConnectInfo) {
	defer close(c)

	uri := mongo.cfg.GetMongoUri()

	var (
		session  *mgo.Session
		err      error
		attempts int
	)

	for {
		session, err = mgo.Dial(uri)
		attempts++

		if err != nil {
			details := make(map[string]string)
			details["uri"] = uri

			e := kerr.NewErr(kerr.ErrLvlFatal, ErrorDomain, DialErrCode, err, details)
			c <- ConnectInfo{
				Error: e,
				Done:  false,
			}

			if mongo.cfg.Reconnect == false {
				break
			}

			timeout := mongo.cfg.ReconnectTimeout * time.Second
			time.Sleep(timeout)

			continue
		}

		break
	}

	mongo.session = session

	c <- ConnectInfo{
		Error: nil,
		Done:  true,
	}
}

// Close and destroy mgo session
func (mongo *Mongo) Shutdown() {
	mongo.session.Close()

	mongo.session = nil
}

// Returns mongodb session
// For real concurrency rock'n'roll copy session, count connection limit
func (mongo *Mongo) GetSession(copy bool) (*mgo.Session, error) {
	if mongo.session == nil {
		e := kerr.NewErr(
			kerr.ErrLvlError,
			ErrorDomain,
			NotConnected,
			fmt.Errorf("mongo not connected yet"),
			nil)

		return nil, e
	}

	if copy == true && mongo.connLimit < mongo.cfg.MaxSessions {
		return mongo.session.Copy(), nil
	}

	return mongo.session, nil
}
