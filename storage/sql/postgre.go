package sql

import (
	"github.com/Lawliet-Chan/yu/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSql struct {
	*gorm.DB
}

func NewPostgreSql(dsn string) (*PostgreSql, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PostgreSql{db}, nil
}

func (p *PostgreSql) Db() *gorm.DB {
	return p.DB
}

func (p *PostgreSql) CreateIfNotExist(table interface{}) error {
	if p.Migrator().HasTable(table) {
		return nil
	}
	return p.Migrator().CreateTable(table)
}

func (*PostgreSql) Type() storage.StoreType {
	return storage.Server
}

func (*PostgreSql) Kind() storage.StoreKind {
	return storage.SQL
}
