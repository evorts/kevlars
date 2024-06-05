/**
 * @Author: steven
 * @Description:
 * @File: helpers
 * @Date: 17/12/23 22.08
 */

package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type IHelper interface {
	BuildSqlAndArgs() (string, []interface{})
	BuildSqlAndArgsWithWherePrefix() (string, []interface{})
	BuildSqlAndArgsFilterOnly() (string, []interface{})
	BuildSqlAndArgsFilterOnlyWithWherePrefix() (string, []interface{})

	OrdersBy() OrdersBy
	Filters() *Filters
	Pagination() *Pagination
}

type helper struct {
	ordersBy   OrdersBy
	filters    *Filters
	pagination *Pagination

	separator Separator
}

func (h *helper) OrdersBy() OrdersBy {
	return h.ordersBy
}

func (h *helper) Filters() *Filters {
	return h.filters
}

func (h *helper) Pagination() *Pagination {
	return h.pagination
}

func (h *helper) BuildSqlAndArgs() (string, []interface{}) {
	q := make([]string, 0)
	args := make([]interface{}, 0)
	if h.filters != nil {
		fq, fv := h.filters.Build(h.separator.String())
		q = append(q, fq)
		args = append(args, fv...)
	}
	if h.ordersBy != nil && len(h.ordersBy) > 0 {
		q = append(q, h.ordersBy.Build())
	}
	// limit, offset should be at the last order
	if h.pagination != nil {
		q = append(q, h.pagination.Build())
	}
	return fmt.Sprintf(" %s", strings.Join(q, " ")), args
}

func (h *helper) BuildSqlAndArgsWithWherePrefix() (string, []interface{}) {
	where, args := h.BuildSqlAndArgs()
	if len(args) < 1 {
		return where, args
	}
	return fmt.Sprintf(" WHERE %s", where), args
}

func (h *helper) BuildSqlAndArgsFilterOnly() (string, []interface{}) {
	q := make([]string, 0)
	args := make([]interface{}, 0)
	if h.filters != nil {
		fq, fv := h.filters.Build(h.separator.String())
		q = append(q, fq)
		args = append(args, fv...)
	}
	return strings.Join(q, " "), args
}

func (h *helper) BuildSqlAndArgsFilterOnlyWithWherePrefix() (string, []interface{}) {
	filterWhere, args := h.BuildSqlAndArgsFilterOnly()
	if len(filterWhere) < 1 {
		return filterWhere, args
	}
	return fmt.Sprintf(" WHERE %s", filterWhere), args
}

func NewHelper(separator Separator, opts ...IHelperOption) IHelper {
	h := &helper{
		separator:  separator,
		pagination: nil,
		ordersBy:   make(OrdersBy, 0),
		filters:    nil,
	}
	for _, opt := range opts {
		opt.apply(h)
	}
	return h
}

func ToJsonObjectFromMap(v map[string]interface{}) *JsonObject {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	jo := new(JsonObject)
	jo.JSONText = b
	jo.Valid = true
	return jo
}

func ToJsonObjectFromInterface(v interface{}) *JsonObject {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	jo := new(JsonObject)
	jo.JSONText = b
	jo.Valid = true
	return jo
}

func ToJsonArray(v []interface{}) *JsonArray {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	jo := new(JsonArray)
	jo.JSONText = b
	jo.Valid = true
	return jo
}

func PlaceholderRepeat(placeholder string, repeat int) []string {
	rs := make([]string, 0)
	for i := 0; i < repeat; i++ {
		rs = append(rs, placeholder)
	}
	return rs
}

func BuildPlaceholder(count int) string {
	return strings.Join(strings.Split(strings.Repeat("?", count), ""), ",")
}

func ToDriverValueFromStringArray(collection []string) []driver.Value {
	rs := make([]driver.Value, 0)
	for _, v := range collection {
		rs = append(rs, v)
	}
	return rs
}
