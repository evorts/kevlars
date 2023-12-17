/**
 * @Author: steven
 * @Description:
 * @File: helpers_model
 * @Date: 18/12/23 00.18
 */

package db

import (
	"fmt"
	"strings"
)

type Pagination struct {
	Page   int
	Limit  int
	offset int
}

func (p Pagination) calcOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p Pagination) Build() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit, p.calcOffset())
}

func (p Pagination) BuildWithSpacePrefix() string {
	return " " + p.Build()
}

type FilterItem struct {
	Field string
	Op    Operator
	Value interface{}
}

func (f *FilterItem) build() (string, interface{}) {
	return fmt.Sprintf("%s %v %v", f.Field, f.Op, "?"), f.Value
}

type FilterItems []FilterItem

type FilterIn struct {
	Field  string
	Values []interface{}
}

func (f *FilterIn) build() (string, []interface{}) {
	return fmt.Sprintf("%s in (%s)", f.Field, strings.Join(strings.Split(strings.Repeat("?", len(f.Values)), ""), ",")), f.Values
}

type FilterIns []FilterIn

type OrderBy struct {
	Field string
	Sort  Sort
}

func (o *OrderBy) build() string {
	return fmt.Sprintf("%s %s", o.Field, o.Sort)
}

type OrdersBy []OrderBy

func (o OrdersBy) FindField(fl string, fb FindBy) OrdersBy {
	rs := make(OrdersBy, 0)
	for _, by := range o {
		if !fb.Match(fl, by.Field) {
			continue
		}
		rs = append(rs, by)
	}
	return rs
}

func (o OrdersBy) Build() string {
	orderBy := make([]string, 0)
	for _, order := range o {
		orderBy = append(orderBy, order.build())
	}
	return fmt.Sprintf("ORDER BY %s", strings.Join(orderBy, ", "))
}

func (o OrdersBy) BuildWithSpacePrefix() string {
	return " " + o.Build()
}

type FindBy string

const (
	FindByExact  FindBy = "exact"
	FindByPrefix FindBy = "prefix"
)

func (t FindBy) Match(input, target string) bool {
	switch t {
	case FindByPrefix:
		return strings.HasPrefix(target, input)
	case FindByExact:
		fallthrough
	default:
		return input == target
	}
}

type Filters struct {
	Ands FilterItems
	Ors  FilterItems
	Ins  FilterIns
	In   *FilterIn
}

func (f *Filters) FindField(fl string, fb FindBy) Filters {
	rs := Filters{
		Ands: make(FilterItems, 0),
		Ors:  make(FilterItems, 0),
		Ins:  make(FilterIns, 0),
		In:   nil,
	}
	if f.Ands != nil && len(f.Ands) > 0 {
		for _, and := range f.Ands {
			if !fb.Match(fl, and.Field) {
				continue
			}
			rs.Ands = append(rs.Ands, and)
		}
	}
	if f.Ors != nil && len(f.Ors) > 0 {
		for _, or := range f.Ors {
			if !fb.Match(fl, or.Field) {
				continue
			}
			rs.Ors = append(rs.Ors, or)
		}
	}
	if f.Ins != nil && len(f.Ins) > 0 {
		for _, in := range f.Ins {
			if !fb.Match(fl, in.Field) {
				continue
			}
			rs.Ins = append(rs.Ins, in)
		}
	}
	if f.In != nil && fb.Match(fl, f.In.Field) {
		rs.In = f.In
	}
	return rs
}

func (f *Filters) Build(separator string) (string, []interface{}) {
	fs := make([]string, 0)
	fargs := make([]interface{}, 0)
	if f.Ands != nil && len(f.Ands) > 0 {
		ands := make([]string, 0)
		for _, and := range f.Ands {
			s, v := and.build()
			ands = append(ands, s)
			fargs = append(fargs, v)
		}
		fs = append(fs, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
	}
	if f.Ors != nil && len(f.Ors) > 0 {
		ors := make([]string, 0)
		for _, or := range f.Ors {
			s, v := or.build()
			ors = append(ors, s)
			fargs = append(fargs, v)
		}
		fs = append(fs, fmt.Sprintf("(%s)", strings.Join(ors, " OR ")))
	}
	if f.In != nil {
		s, v := f.In.build()
		fs = append(fs, s)
		fargs = append(fargs, v...)
	}
	if f.Ins != nil && len(f.Ins) > 0 {
		ins := make([]string, 0)
		for _, in := range f.Ins {
			s, v := in.build()
			ins = append(ins, s)
			fargs = append(fargs, v...)
		}
		fs = append(fs, fmt.Sprintf("(%s)", strings.Join(ins, " AND ")))
	}
	return strings.Join(fs, separator), fargs
}

func (f *Filters) BuildWithWherePrefix(separator string) (string, []interface{}) {
	where, args := f.Build(separator)
	if len(args) < 1 {
		return where, args
	}
	return fmt.Sprintf(" WHERE %s", where), args
}
