// Package postgres implements postgres connection.
package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// Postgres -.
type Postgres struct {
	*gorm.DB
}

// New -.
func New(url string) (*Postgres, error) {

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

// Close -.
func (p *Postgres) Close() {
	d, err := p.DB.DB()
	if err != nil {
		log.Printf("error closing database: %s", err.Error())
	}

	d.Close()
}
