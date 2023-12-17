/**
 * @Author: steven
 * @Description:
 * @File: models_pg
 * @Date: 17/12/23 22.14
 */

package db

import (
	"database/sql/driver"
	"github.com/lib/pq"
)

type PqArrayOfString pq.StringArray

//goland:noinspection GoMixedReceiverTypes
func (a *PqArrayOfString) base() *pq.StringArray {
	base := pq.StringArray(*a)
	return &base
}

// Scan implements the sql.Scanner interface.
//
//goland:noinspection GoMixedReceiverTypes
func (a *PqArrayOfString) Scan(src interface{}) error {
	return a.base().Scan(src)
}

// Value implements the driver.Valuer interface.
//
//goland:noinspection GoMixedReceiverTypes
func (a PqArrayOfString) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return a.base().Value()
}
