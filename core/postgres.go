package core

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Uri     string `yaml:"uri"`
	MaxIdle int    `yaml:"maxIdle" default:"10"`
	MaxOpen int    `yaml:"maxOpen" default:"100"`
}

func NewPostgres(conf *PostgresConfig) *sqlx.DB {
	uri, err := url.Parse(conf.Uri)
	if err != nil {
		panic(fmt.Sprintf("parse %s err:%s", conf.Uri, err.Error()))
	}

	poll, err := sqlx.Open(uri.Scheme, conf.Uri)
	if err != nil {
		panic(fmt.Sprintf("open %s err:%s", conf.Uri, err.Error()))
	}

	err = poll.Ping()
	if err != nil {
		panic(fmt.Sprintf("ping %s err:%s", conf.Uri, err.Error()))
	}
	poll.SetMaxIdleConns(conf.MaxIdle)
	poll.SetMaxOpenConns(conf.MaxOpen)
	poll.SetConnMaxLifetime(2 * time.Minute)
	poll.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return poll
}
