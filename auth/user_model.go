/**
 * @Author: steven
 * @Description:
 * @File: auth_model
 * @Date: 14/05/24 11.51
 */

package auth

import (
	"database/sql"
	"time"
)

type UserAuthRecord struct {
	ID         int        `db:"id"`
	UserID     int        `db:"user_id"`
	Email      string     `db:"email"`
	Phone      string     `db:"phone"`
	Creds      string     `db:"creds"`
	Disabled   bool       `db:"disabled"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	DisabledAt *time.Time `db:"disabled_at"`
	ExpiredAt  *time.Time `db:"expired_at"`
}

type UserAuthRecords []*UserAuthRecord

type UserAccessRecord struct {
	ID         int64        `db:"id"`
	UserID     int64        `db:"user_id"`
	Resource   string       `db:"resource"`
	Scopes     Scopes       `db:"scopes"`
	Disabled   bool         `db:"disabled"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
	DisabledAt sql.NullTime `db:"disabled_at"`
}

type UserAccessRecords []*UserAccessRecord

type UserWithAccessRecord struct {
	*UserAuthRecord
	Permissions UserAccessRecords `db:"permissions"`
}

type UserWithAccessRecords []*UserAccessRecord
