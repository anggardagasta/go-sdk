package database

import (
	"fmt"

	_ "github.com/lib/pq"
)

func (d *Database) ConnectionPostgres() (DBInterface, error) {
	cfg := d.config
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", cfg.Address, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	return d.connectiondb(connStr)
}
