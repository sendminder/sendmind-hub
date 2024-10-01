package database

import (
	"log"
	"sendmind-hub/pkg/config"
	"sendmind-hub/pkg/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type DB struct {
	Conn *pg.DB
}

func NewDB(cfg *config.Config) *DB {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.PostGresAddr,
		User:     cfg.PostGresUser,
		Password: cfg.PostGresPassword,
		Database: cfg.PostGresDB,
	})

	err := createSchema(db)
	if err != nil {
		log.Fatalf("Error creating schema: %v\n", err)
	}

	return &DB{Conn: db}
}

func createSchema(db *pg.DB) error {
	// 테이블 스키마 생성
	return db.Model((*model.User)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
}
