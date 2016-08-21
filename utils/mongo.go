package utils

import (
	"errors"
	mgo "gopkg.in/mgo.v2"
)

/*MongoDb is a wrapper on top of mgo*/
type MongoDb struct {
	Host     string
	Database string
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
	}

	return session
}

func NewMongoDBWithConfig(config map[string]interface{}) *MongoDb {
	return &MongoDb{
		Host:     config["host"].(string),
		Database: config["database"].(string),
	}

}
