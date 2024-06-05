/**
 * @Author: steven
 * @Description:
 * @File: common
 * @Date: 14/05/24 12.28
 */

package auth

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/evorts/kevlars/utils"
	"github.com/lib/pq"
	"net/http"
	"strings"
)

type Scope string

const (
	ScopeRead      Scope = "read"
	ScopeWrite     Scope = "write"
	ScopeDelete    Scope = "delete"
	ScopeUndefined Scope = "undefined"
)

//goland:noinspection GoMixedReceiverTypes
func (s Scope) FromHttpMethod(method string) Scope {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		return ScopeRead
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return ScopeWrite
	case http.MethodDelete:
		return ScopeDelete
	default:
		return ScopeUndefined
	}
}

//goland:noinspection GoMixedReceiverTypes
func (s Scope) String() string { return string(s) }

// Scan implements the sql.Scanner interface.
//
//goland:noinspection GoMixedReceiverTypes
func (s *Scope) Scan(src interface{}) error {
	if v, ok := src.(string); ok {
		*s = Scope(v)
		return nil
	}
	return errors.New("invalid scope value")
}

// Value implements the driver.Valuer interface.
//
//goland:noinspection GoMixedReceiverTypes
func (s Scope) Value() (driver.Value, error) {
	return driver.Value(s.String()), nil
}

type Scopes []Scope

// Scan implements the sql.Scanner interface.
//
//goland:noinspection GoMixedReceiverTypes
func (s *Scopes) Scan(src interface{}) error {
	if src == nil {
		*s = Scopes{}
		return nil
	}
	var (
		err         error
		arrOfString = &pq.StringArray{}
	)

	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &arrOfString)
	case []byte:
		err = arrOfString.Scan(src)
	default:
		err = errors.New("incompatible source type for scopes")
	}
	if err != nil {
		return err
	}
	*s = Scopes{}.FromStringArray(*arrOfString)
	return nil
}

// Value implements the driver.Valuer interface.
//
//goland:noinspection GoMixedReceiverTypes
func (s Scopes) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return pq.StringArray(s.ToStringArray()).Value()
}

//goland:noinspection GoMixedReceiverTypes
func (s Scopes) FromStringArray(values []string) Scopes {
	rs := make(Scopes, len(values))
	for i, v := range values {
		rs[i] = Scope(v)
	}
	return rs
}

//goland:noinspection GoMixedReceiverTypes
func (s Scopes) ToStringArray() []string {
	rs := make([]string, len(s))
	for i, v := range s {
		rs[i] = v.String()
	}
	return rs
}

//goland:noinspection GoMixedReceiverTypes
func (s Scopes) AllowedToRead() bool {
	if len(s) == 0 {
		return false
	}
	return utils.InArray(s, ScopeRead)
}

//goland:noinspection GoMixedReceiverTypes
func (s Scopes) AllowedToWrite() bool {
	if len(s) == 0 {
		return false
	}
	return utils.InArray(s, ScopeWrite)
}

//goland:noinspection GoMixedReceiverTypes
func (s Scopes) AllowedTo(scope Scope) bool {
	if len(s) == 0 {
		return false
	}
	return utils.InArray(s, scope)
}

//goland:noinspection GoMixedReceiverTypes
func (s Scopes) IsAllowedByHttpMethod(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		return utils.InArray(s, ScopeRead)
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return utils.InArray(s, ScopeWrite)
	case http.MethodDelete:
		return utils.InArray(s, ScopeDelete)
	default:
		return false
	}
}
