package postgres

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
	Url     string `json:"url"`
	MaxIdle int    `json:"maxIdle" default:"10"`
	MaxOpen int    `json:"maxOpen" default:"100"`
}

func NewPostgres(conf *PostgresConfig) *sqlx.DB {
	uri, err := url.Parse(conf.Url)
	if err != nil {
		panic(fmt.Sprintf("parse %s err:%s", conf.Url, err.Error()))
	}

	poll, err := sqlx.Open(uri.Scheme, conf.Url)
	if err != nil {
		panic(fmt.Sprintf("open %s err:%s", conf.Url, err.Error()))
	}

	err = poll.Ping()
	if err != nil {
		panic(fmt.Sprintf("ping %s err:%s", conf.Url, err.Error()))
	}
	poll.SetMaxIdleConns(conf.MaxIdle)
	poll.SetMaxOpenConns(conf.MaxOpen)
	poll.SetConnMaxLifetime(2 * time.Minute)
	poll.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return poll
}
