package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Type            string `yaml:"type" json:"type"`
	Address         string `yaml:"address" json:"address"`
	Port            string `yaml:"port" json:"port"`
	Username        string `yaml:"username" json:"username"`
	Password        string `yaml:"password" json:"password"`
	DBName          string `yaml:"dbname" json:"dbname"`
	MaxOpenConn     int    `yaml:"maxopenconn" json:"maxopenconn"`
	MaxIdleConn     int    `yaml:"maxidleconn" json:"maxidleconn"`
	ConnMaxIdleTime int    `yaml:"connmaxidletime" json:"connmaxidletime"`
	ConnMaxLifeTime int    `yaml:"connmaxlifetime" json:"connmaxlifetime"`
}

type Database struct {
	config Config
}

func NewConnection(dbconf Config) (DBInterface, error) {
	dbcon := &Database{
		config: dbconf,
	}

	if dbconf.Type == "mysql" {
		return dbcon.ConnectionMySQL()
	}

	if dbconf.Type == "postgres" {
		return dbcon.ConnectionPostgres()
	}

	return nil, fmt.Errorf("database type %s not supported", dbconf.Type)
}

func (d *Database) connectiondb(connStr string) (DBInterface, error) {
	cfg := d.config

	db, err := sqlx.Connect(cfg.Type, connStr)
	if err != nil {
		return nil, err
	}
	if cfg.MaxOpenConn != -1 {
		db.SetMaxOpenConns(cfg.MaxOpenConn)
	}
	if cfg.MaxIdleConn != -1 {
		db.SetMaxIdleConns(cfg.MaxIdleConn)
	}
	if cfg.ConnMaxLifeTime != -1 {
		db.SetConnMaxLifetime(time.Second * time.Duration(cfg.ConnMaxLifeTime))
	}
	if cfg.ConnMaxIdleTime != -1 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(cfg.ConnMaxIdleTime))
	}

	return db, nil
}
