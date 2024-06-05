/**
 * @Author: steven
 * @Description:
 * @File: db_test
 * @Date: 05/06/24 09.47
 */

package db

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type dbTestSuite struct {
	suite.Suite
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(dbTestSuite))
}
