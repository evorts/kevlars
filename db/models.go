/**
 * @Author: steven
 * @Description:
 * @File: models
 * @Date: 17/12/23 22.13
 */

package db

import (
	"database/sql"
	"database/sql/driver"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
	"time"
)

type JsonObject struct {
	types.JSONText
	Valid bool
}

type JsonArray struct {
	types.JSONText
	Valid bool
}

var emptyJSON = types.JSONText("{}")
var emptyJSONArray = types.JSONText("[]")

//goland:noinspection GoUnusedGlobalVariable
var EmptyJsonObject = &JsonObject{
	JSONText: emptyJSON,
	Valid:    true,
}

//goland:noinspection GoUnusedGlobalVariable
var EmptyJsonArray = &JsonArray{
	JSONText: emptyJSONArray,
	Valid:    true,
}

// Scan implements the Scanner interface.
//
//goland:noinspection GoMixedReceiverTypes
func (n *JsonObject) Scan(value interface{}) error {
	if value == nil {
		n.JSONText, n.Valid = emptyJSON, false
		return nil
	}
	n.Valid = true
	return n.JSONText.Scan(value)
}

// Value implements the driver Valuer interface.
//
//goland:noinspection GoMixedReceiverTypes
func (n JsonObject) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.JSONText.Value()
}

func (n *JsonArray) Scan(value interface{}) error {
	if value == nil {
		n.JSONText, n.Valid = emptyJSONArray, false
		return nil
	}
	n.Valid = true
	return n.JSONText.Scan(value)
}

func (n JsonArray) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.JSONText.Value()
}

type NullableTime sql.NullTime

func (a NullableTime) From(dt *time.Time) NullableTime {
	if dt == nil {
		return NullableTime{Valid: false}
	}
	return NullableTime{
		Time:  *dt,
		Valid: true,
	}
}

func (a NullableTime) GetTime() *time.Time {
	if a.Valid {
		return &a.Time
	}
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (a NullableTime) Base() *sql.NullTime {
	v := sql.NullTime(a)
	return &v
}

// Scan implements the sql.Scanner interface.
//
//goland:noinspection GoMixedReceiverTypes
func (a *NullableTime) Scan(src interface{}) error {
	v := sql.NullTime(*a)
	err := v.Scan(src)
	*a = NullableTime(v)
	return err
}

// Value implements the driver.Valuer interface.
//
//goland:noinspection GoMixedReceiverTypes
func (a *NullableTime) Value() (driver.Value, error) {
	return a.Base().Value()
}
