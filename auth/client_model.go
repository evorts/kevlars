/**
 * @Author: steven
 * @Description:
 * @File: model
 * @Date: 24/12/23 21.51
 */

package auth

import (
	"database/sql"
	"time"
)

type Client struct {
	ID         int          `db:"id"`
	Name       string       `db:"name"`
	Secret     string       `db:"secret"`
	ExpiredAt  sql.NullTime `db:"expired_at"`
	Disabled   bool         `db:"disabled"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
	DisabledAt sql.NullTime `db:"disabled_at"`
}

type Clients []*Client

type ClientScope struct {
	ID         int          `db:"id"`
	ClientID   int          `db:"client_id"`
	Resource   string       `db:"resource"`
	Scopes     Scopes       `db:"scopes"`
	Disabled   bool         `db:"disabled"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
	DisabledAt sql.NullTime `db:"disabled_at"`
}

type ClientScopes []*ClientScope

type ClientWithScopes struct {
	*Client
	Scopes ClientScopes
}

type ClientsWithScopes []*ClientWithScopes

type clientDataForAuthorization struct {
	ClientName string
	Scopes     Scopes
	Disabled   bool
	ExpiredAt  *time.Time
}

// map[secret][resource]ClientForAuthorization
type mapClientAuthorization map[string]map[string]clientDataForAuthorization
