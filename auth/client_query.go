/**
 * @Author: steven
 * @Description:
 * @File: client_query
 * @Date: 05/06/24 19.44
 */

package auth

import (
	"fmt"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/utils"
	"strings"
)

/** Add Query **/
var (
	// addClientQuery
	addClientQuery = map[db.SupportedDriver]struct {
		placeholder func(int) []string
		query       func(v string) string
	}{
		db.DriverPostgreSQL: {
			placeholder: func(repeat int) []string {
				return db.PlaceholderRepeat(
					fmt.Sprintf(
						`(%s,CASE WHEN ? THEN current_timestamp END)`,
						strings.Join(utils.RepeatInSlice("?", 4), ","),
					), repeat,
				)
			},
			query: func(v string) string {
				return `INSERT INTO ` + tableClients + `(name, secret, disabled, expired_at, disabled_at)
					VALUES ` + v + `
					RETURNING id, name, disabled, expired_at, created_at, disabled_at
		`
			},
		},
	}
	// addScopeQuery
	addScopeQuery = map[db.SupportedDriver]struct {
		placeholder func(int) []string
		query       func(v string) string
	}{
		db.DriverPostgreSQL: {
			placeholder: func(repeat int) []string {
				return db.PlaceholderRepeat(
					fmt.Sprintf(
						`(%s,CASE WHEN ? THEN current_timestamp END)`,
						strings.Join(utils.RepeatInSlice("?", 4), ","),
					), repeat,
				)
			},
			query: func(v string) string {
				return `INSERT INTO ` + tableClientScope + `(client_id, resource, scopes, disabled, disabled_at)
						VALUES ` + v + `
						RETURNING id, client_id, resource, scopes, disabled, created_at, disabled_at
				`
			},
		},
	}
)

/** Remove/Void Query **/
var (
	removeClientByIdsQuery = map[db.SupportedDriver]struct {
		query func(pl int) string
	}{
		db.DriverPostgreSQL: {
			query: func(pl int) string {
				return `DELETE FROM` + " " + tableClients + ` WHERE id IN(` + db.BuildPlaceholder(pl) + `)`
			},
		},
	}
	voidClientByIdsQuery = map[db.SupportedDriver]struct {
		query func(pl int) string
	}{
		db.DriverPostgreSQL: {
			query: func(pl int) string {
				return `UPDATE ` + tableClients + ` SET disabled=true, disabled_at=current_timestamp WHERE id IN(` + db.BuildPlaceholder(pl) + `)`
			},
		},
	}
	removeClientScopesByIdsQuery = map[db.SupportedDriver]struct {
		query func(pl int) string
	}{
		db.DriverPostgreSQL: {
			query: func(pl int) string {
				return `DELETE FROM` + " " + tableClientScope + ` WHERE id IN(` + db.BuildPlaceholder(pl) + `)`
			},
		},
	}
	voidClientScopesByIdsQuery = map[db.SupportedDriver]struct {
		query func(pl int) string
	}{
		db.DriverPostgreSQL: {
			query: func(pl int) string {
				return `UPDATE ` + tableClientScope + ` SET disabled=true, disabled_at=current_timestamp WHERE id IN(` + db.BuildPlaceholder(pl) + `)`
			},
		},
	}
)

/** Get Query **/
var (
	getClientWithScopesByQuery = map[db.SupportedDriver]struct {
		query func(qf string) string
	}{
		db.DriverPostgreSQL: {
			query: func(qf string) string {
				return `SELECT 
		    				c.id, c.name, c.secret, c.expired_at, c.disabled, 
		    				c.created_at, c.updated_at, c.disabled_at,
		    				cs.id as scope_id, cs.resource, cs.scopes, cs.disabled as scope_disabled,
		    				cs.created_at as scope_created_at, cs.updated_at as scope_updated_at,
		    				cs.disabled_at as scope_disabled_at
						FROM` + " " + tableClients + ` c JOIN ` + tableClientScope +
					` cs ON cs.client_id = c.id ` + qf
			},
		},
	}

	getClientsByQuery = map[db.SupportedDriver]struct {
		query func(qf string) string
	}{
		db.DriverPostgreSQL: {
			query: func(qf string) string {
				return `SELECT 
					id, name, secret, expired_at, disabled, created_at, updated_at, disabled_at
				FROM` + " " + tableClients + " " + qf
			},
		},
	}

	getClientScopesByQuery = map[db.SupportedDriver]struct {
		query func(qf string) string
	}{
		db.DriverPostgreSQL: {
			query: func(qf string) string {
				return `SELECT 
					id, client_id, resource, scopes, disabled, created_at, updated_at, disabled_at
				FROM` + " " + tableClientScope + " " + qf
			},
		},
	}
)
