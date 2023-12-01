package postgres

import (
	"github.com/madyar997/sso-jcode/internal/database/drivers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Postgres struct {
	url    string
	client *gorm.DB
}

func New(conf map[string]string) (drivers.DataStore, error) {
	if conf == nil {
		return nil, drivers.ErrInvalidConfigStruct
	}

	if _, ok := (conf)["url"]; !ok {
		return nil, drivers.ErrInvalidConfigStruct
	}
	return &Postgres{url: (conf)["url"]}, nil
}

func (p *Postgres) Name() string {
	return "postgres"
}

func (p *Postgres) Close() error {
	d, err := p.client.DB()
	if err != nil {
		log.Printf("error closing database: %s", err.Error())
		return err
	}

	return d.Close()
}

func (p *Postgres) Connect() error {
	db, err := gorm.Open(postgres.Open(p.url), &gorm.Config{})
	if err != nil {
		return err
	}

	p.client = db

	return nil
}
