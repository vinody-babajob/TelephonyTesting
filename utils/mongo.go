package mongo

import (
	"Babajob/Telephony.Utilities/configuration"
	"Babajob/Telephony.Utilities/constants"
	"Babajob/Telephony.Utilities/logger"
	"errors"
	mgo "gopkg.in/mgo.v2"
	"os"
)

/*MongoDb is a wrapper on top of mgo*/
type MongoDb struct {
	Host     string
	Database string
	logger   logger.Logger
}

/*GetCollection returns the collection given the name"*/
func (mongoDb *MongoDb) GetCollection(name string) (*mgo.Session, *mgo.Collection, error) {
	session := mongoDb.getSession()

	allDatabaseCollections, err := session.DB(mongoDb.Database).CollectionNames()

	if err != nil {
		defer session.Close()
		return nil, nil, err
	}

	for _, collectionName := range allDatabaseCollections {
		if collectionName == name {
			collection := session.DB(mongoDb.Database).C(name)
			return session, collection, nil
		}
	}

	defer session.Close()

	return nil, nil, errors.New("Invalid collection name : " + name)
}

func (mongoDb *MongoDb) getSession() *mgo.Session {
	session, err := mgo.Dial(mongoDb.Host)

	if err != nil {
		mongoDb.logger.LogFatal("Failed to connect to MongoDB")
	}

	return session
}

func NewMongoDBWithConfig(config map[string]interface{}) *MongoDb {
	return &MongoDb{
		Host:     config["host"].(string),
		Database: config["database"].(string),
		logger:   logger.NewLogger(),
	}

}

func NewMongoDb(basePath string) *MongoDb {
	var configReader = configuration.NewConfigurationReader(
		os.Getenv(constants.TelephonyEnvironment),
		basePath,
	)

	return &MongoDb{
		Host:     configReader.GetValue("host"),
		Database: configReader.GetValue("database"),
		logger:   logger.NewLogger(),
	}
}
