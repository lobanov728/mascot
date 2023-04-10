package adapters

import (
	"fmt"
	"net/url"
)

type Postgres struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

// URL for postgres.
func (p *Postgres) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User,
		url.QueryEscape(p.Password),
		p.Host,
		p.Port,
		p.Database,
	)
}
