package mongo

import (
	"context"
	"github.com/madyar997/sso-jcode/internal/database/drivers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	connectionTimeout = 3 * time.Second
	ensureIdxTimeout  = 10 * time.Second
	retries           = 1
)

type Mongo struct {
	MongoURL string
	client   *mongo.Client
	dbname   string

	DB      *mongo.Database
	Context context.Context

	retries           int
	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration
}

func (m *Mongo) Name() string { return "Mongo" }

func New(conf map[string]string) (drivers.DataStore, error) {
	if conf == nil {
		return nil, drivers.ErrInvalidConfigStruct
	}

	if _, ok := (conf)["url"]; !ok {
		return nil, drivers.ErrInvalidConfigStruct
	}

	if _, ok := (conf)["db"]; !ok {
		return nil, drivers.ErrInvalidConfigStruct
	}

	return &Mongo{
		MongoURL:          (conf)["url"],
		dbname:            (conf)["db"],
		retries:           retries,
		connectionTimeout: connectionTimeout,
		ensureIdxTimeout:  ensureIdxTimeout,
	}, nil
}

func (m *Mongo) Connect() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.MongoURL))
	if err != nil {
		return err
	}

	if err := m.Ping(); err != nil {
		return err
	}

	m.DB = m.client.Database(m.dbname)

	// убеждаемся что созданы все необходимые индексы
	return m.ensureIndexes()
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	return m.client.Ping(ctx, readpref.Primary())
}

func (m *Mongo) Close() error {
	return m.client.Disconnect(m.Context)
}

// убеждается что все индексы построены
func (m *Mongo) ensureIndexes() error {
	//ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	//defer cancel()

	return nil
}
