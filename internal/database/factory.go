package database

import (
	"github.com/madyar997/sso-jcode/internal/database/drivers"
	"github.com/madyar997/sso-jcode/internal/database/drivers/mongo"
	"github.com/madyar997/sso-jcode/internal/database/drivers/postgres"
	"log"
)

type DataStoreFactory func(conf map[string]string) (drivers.DataStore, error)

var dataStoreFactories = make(map[string]DataStoreFactory)

func Register(name string, factory DataStoreFactory) {
	if factory == nil {
		log.Panicf("datastore factory %s does not exist.", name)
	}

	_, registered := dataStoreFactories[name]
	if registered {
		log.Printf("datastore factory %s already registered. Ignoring.", name)
	} else {
		dataStoreFactories[name] = factory
	}
}

func init() {
	Register("postgres", postgres.New)
	Register("mongo", mongo.New)
	//регистрация монги
}

func New(conf map[string]string) (drivers.DataStore, error) {
	if conf == nil {
		return nil, ErrEmptyConfigStruct
	}

	engineName := (conf)["datastore"]
	engineFactory, ok := dataStoreFactories[engineName]

	if !ok {

		availableDataStores := make([]string, 0, len(dataStoreFactories)-1)
		for k := range dataStoreFactories {
			availableDataStores = append(availableDataStores, k)
		}

		return nil, ErrInvalidDatastoreName.Error(availableDataStores)
	}

	return engineFactory(conf)
}

func Connect(conf map[string]string) (drivers.DataStore, error) {
	ds, err := New(conf)
	if err != nil {
		log.Printf("[ERROR] cannot create datastore: %v", err)
		return nil, err
	}

	if err = ds.Connect(); err != nil {
		log.Printf("[ERROR] cannot connect to database %s: %v", ds.Name(), err)
		return nil, err
	}

	return ds, nil
}
