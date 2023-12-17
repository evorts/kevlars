/**
 * @Author: steven
 * @Description:
 * @File: helpers_test
 * @Date: 17/12/23 23.53
 */

package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type helperTestSuite struct {
	suite.Suite
}

func TestHelper_Build(t *testing.T) {
	orders := []OrderBy{
		{
			Field: "Field",
			Sort:  "ASC",
		},
	}
	filters := Filters{
		Ands: []FilterItem{
			{
				Field: "status",
				Op:    OpEq,
				Value: "SAVED",
			},
		},
		Ors: nil,
		In:  nil,
	}
	h := NewHelper(SeparatorAND, WithPagination(1, 10), WithOrdersBy(orders), WithFilters(filters))
	qf, args := h.BuildSqlAndArgs()
	assert.Equal(t, []interface{}{"SAVED"}, args, "Arguments test")
	assert.Equal(t, " (status = ?) ORDER BY Field ASC LIMIT 10 OFFSET 0", qf, "Query filter")
}
